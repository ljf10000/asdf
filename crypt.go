package asdf

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

const (
	CryptKeySize = 32
)

const (
	CryptAes CryptType = iota

	CryptEnd
)

var cryptEnumManager = EnumManagerCreate()

func cryptEnumManagerInit() {
	cryptEnumManager.Register(int(CryptAes), "aes")
}

const (
	Sha224 HmacType = iota
	Sha256
	Sha384
	Sha512

	ShaEnd
)

var hmacEnumManager = EnumManagerCreate()

func hmacEnumManagerInit() {
	hmacEnumManager.Register(int(Sha224), "sha224")
	hmacEnumManager.Register(int(Sha256), "sha256")
	hmacEnumManager.Register(int(Sha384), "sha384")
	hmacEnumManager.Register(int(Sha512), "sha512")
}

type CryptKey [CryptKeySize]byte

func (me *CryptKey) String() string {
	return hex.EncodeToString((*me)[:])
}

func (me *CryptKey) FromString(s string) error {
	buf, err := hex.DecodeString(s)
	if nil != err {
		return err
	}

	copy((*me)[:], buf)
	return nil
}

type HmacType int

func (me HmacType) String() string {
	return hmacEnumManager.Name2(int(me))
}

func (me *HmacType) FromString(s string) error {
	idx, ok := hmacEnumManager.Index(s)
	if ok {
		*me = HmacType(idx)

		return nil
	} else {
		return Error
	}
}

type CryptType int

func (me CryptType) String() string {
	return cryptEnumManager.Name2(int(me))
}

func (me *CryptType) FromString(s string) error {
	idx, ok := cryptEnumManager.Index(s)
	if ok {
		*me = CryptType(idx)

		return nil
	} else {
		return Error
	}
}

func cryptToBinary(b *Crypt, s *CryptString) {
	(&b.Key).FromString(s.Key)
	(&b.CryptType).FromString(s.CryptType)
	(&b.HmacType).FromString(s.HmacType)

	b.init()
}

func cryptToString(s *CryptString, b *Crypt) {
	s.Key = (&b.Key).String()
	s.CryptType = b.CryptType.String()
	s.HmacType = b.HmacType.String()
}

type CryptString struct {
	Key       string `json:"key"`
	CryptType string `json:"crypt"`
	HmacType  string `json:"hmactype"`
}

func (me *CryptString) ToBinary() *Crypt {
	crypt := &Crypt{}

	cryptToBinary(crypt, me)

	return crypt
}

func (me *CryptString) FromBinary(crypt *Crypt) {
	cryptToString(me, crypt)
}

type Crypt struct {
	Key       CryptKey
	CryptType CryptType
	HmacType  HmacType

	Cipher cipher.Block
	Hash   hash.Hash
}

func (me *Crypt) ToString() *CryptString {
	crypt := &CryptString{}

	cryptToString(crypt, me)

	return crypt
}

func (me *Crypt) FromString(crypt *CryptString) {
	cryptToBinary(me, crypt)
}

func (me *Crypt) init() {
	me.Cipher, _ = aes.NewCipher(me.Key[:])
	me.Hash = hmac.New(sha256.New, me.Key[:])
}

func (me *Crypt) size(b []byte) int {
	return AlignI(len(b), me.Cipher.BlockSize())
}

func (me *Crypt) Align(b []byte) []byte {
	Size := me.size(b)

	if len(b) == Size {
		return b
	} else {
		buf := make([]byte, Size)

		copy(buf, b)

		return buf
	}
}

func (me *Crypt) IsAlign(b []byte) bool {
	return len(b) == me.size(b)
}

func (me *Crypt) Count(b []byte) int {
	return me.size(b) / me.Cipher.BlockSize()
}

func (me *Crypt) crypt(dst, src []byte, handler func(dst, src []byte)) error {
	if !me.IsAlign(src) {
		return ErrNoAlign
	} else if len(dst) < len(src) {
		return ErrNoSpace
	}

	blocksize := me.Cipher.BlockSize()
	count := me.Count(src)

	for i := 0; i < count; i++ {
		bsrc := src[i*blocksize : (i+1)*blocksize]
		bdsr := dst[i*blocksize : (i+1)*blocksize]

		handler(bdsr, bsrc)
	}

	return nil
}

func (me *Crypt) Encrypt(dst, src []byte) error {
	return me.crypt(dst, src, me.Cipher.Encrypt)
}

// only for string
func (me *Crypt) Decrypt(dst, src []byte) error {
	err := me.crypt(dst, src, me.Cipher.Decrypt)

	return err
}

func (me *Crypt) Encode(b []byte) []byte {
	buf := me.Align(b)
	count := me.Count(b)

	for i := 0; i < count; i++ {
		blocksize := me.Cipher.BlockSize()
		bin := buf[i*blocksize : (i+1)*blocksize]
		me.Cipher.Encrypt(bin, bin)
	}

	return buf
}

// only for string
func (me *Crypt) Decode(b []byte) ([]byte, error) {
	if !me.IsAlign(b) {
		return nil, ErrNoAlign
	}

	buf := b
	count := me.Count(b)

	for i := 0; i < count; i++ {
		blocksize := me.Cipher.BlockSize()
		bin := buf[i*blocksize : (i+1)*blocksize]
		me.Cipher.Decrypt(bin, bin)
	}

	count = len(buf)
	for i := 0; i < count; i++ {
		if 0 == buf[i] {
			return buf[:i], nil
		}
	}

	return buf, nil
}

// now, only support aes/hmacsha256
func NewCrypt(ctype CryptType, htype HmacType, key []byte) *Crypt {
	c := &Crypt{
		CryptType: CryptAes,
		HmacType:  Sha256,
	}
	if nil != key {
		copy(c.Key[:], key)
	} else {
		RandBuffer(c.Key[:]).Rand()
	}

	c.init()

	return c
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])

	return plantText[:(length - unpadding)]
}

func initCrypt() {
	cryptEnumManagerInit()
	hmacEnumManagerInit()
}
