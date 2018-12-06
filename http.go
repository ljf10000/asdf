package asdf

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	HttpMethodGet    HttpMethod = 0
	HttpMethodPost   HttpMethod = 1
	HttpMethodDelete HttpMethod = 2
	HttpMethodEnd    HttpMethod = 3
)

type HttpMethod byte

var httpMethods = &EnumMapper{
	Enum: "HttpMethod",
	Names: []string{
		HttpMethodGet:    "GET",
		HttpMethodPost:   "POST",
		HttpMethodDelete: "DELETE",
	},
}

func (me HttpMethod) IsGood() bool {
	return httpMethods.IsGoodIndex(int(me))
}

func (me HttpMethod) String() string {
	return httpMethods.Name(int(me))
}

func (me *HttpMethod) FromString(s string) error {
	idx, err := httpMethods.Index(s)
	if nil == err {
		*me = HttpMethod(idx)
	}

	return err
}

/******************************************************************************/

const (
	httpBodyType = "text/plain"
)

// only for string
// buf==>crypt==>hex
func httpEncode(buf []byte, codec ICodec) []byte {
	if nil != buf {
		if nil != codec {
			b := codec.Encode(buf)
			s := base64.StdEncoding.EncodeToString(b)

			return []byte(s)
		} else {
			return buf
		}
	} else {
		return nil
	}
}

// only for string
// hex==>crypt==>buf
func httpDecode(buf []byte, codec ICodec) ([]byte, error) {
	if nil != buf {
		if nil != codec {
			b, _ := base64.StdEncoding.DecodeString(string(buf))

			return codec.Decode(b)
		} else {
			return buf, nil
		}
	} else {
		return nil, ErrBadBuf
	}
}

type HttpCodec struct {
	Crypt *Crypt
}

func NewHttpCodec(crypt *Crypt) *HttpCodec {
	return &HttpCodec{
		Crypt: crypt,
	}
}

func (me *HttpCodec) Decode(b []byte) ([]byte, error) {
	return httpDecode(b, me.Crypt)
}

func (me *HttpCodec) Encode(b []byte) []byte {
	return httpEncode(b, me.Crypt)
}

type IHttpUser interface {
	HttpUser() string
}

type IHttpUrl interface {
	HttpUrl() string
}

type IHttpBodyType interface {
	HttpBodyType() string
}

type IHttpPost interface {
	IHttpUser
	IHttpUrl
	IHttpBodyType
}

type IHttpGet interface {
	IHttpUser
	IHttpUrl
}

func HttpError(w http.ResponseWriter, error int, codec ICodec) {
	HttpReply(w, NewStdError(error), codec)
}

type IHttpLogEnable interface {
	LogEnable() bool
}

// [hex==>crypt==>]json==>obj
func HttpBody(rBody io.ReadCloser, iBody interface{}, codec ICodec) error {
	body, err := ioutil.ReadAll(rBody)
	rBody.Close()

	if nil != err {
		return err
	} else if body, err = httpDecode(body, codec); nil != err {
		return err
	} else if err = json.Unmarshal(body, iBody); nil != err {
		Log.Debug("body[%s] to json error:%s", string(body), err.Error())

		return err
	} else if obj, ok := iBody.(IHttpLogEnable); !ok || obj.LogEnable() {
		Log.Debug("http body[%s] to obj[%+v]", string(body), iBody)
	}

	return nil
}

func HttpGet(iGet IHttpGet, output interface{}, codec ICodec) error {
	Log.Debug("%s get %s ...", iGet.HttpUser(), iGet.HttpUrl())

	r, err := http.Get(iGet.HttpUrl())
	if nil != err {
		Log.Debug("%s get %s error:%s",
			iGet.HttpUser(), iGet.HttpUrl(), err.Error())

		return err
	}

	// Log.Debug("%s get %s ok.", iGet.HttpUser(), iGet.HttpUrl())

	if nil != output {
		return HttpBody(r.Body, output, codec)
	} else {
		return nil
	}
}

func HttpPost(iPost IHttpPost, input, output interface{}, codec ICodec) error {
	buf, _ := json.Marshal(input)
	Log.Debug("%s post %s ...", iPost.HttpUser(), iPost.HttpUrl())

	buf = httpEncode(buf, codec)

	bodyType := iPost.HttpBodyType()
	if Empty == bodyType {
		bodyType = httpBodyType
	}

	r, err := http.Post(iPost.HttpUrl(), bodyType, bytes.NewBuffer(buf))
	if nil != err {
		Log.Debug("%s post %s type:%s data:%s error:%s",
			iPost.HttpUser(), iPost.HttpUrl(), bodyType, string(buf), err.Error())

		return err
	}

	// Log.Debug("%s post %s ok.", iPost.HttpUser(), iPost.HttpUrl())

	if nil != output {
		return HttpBody(r.Body, output, codec)
	} else {
		return nil
	}
}

// output==>json==>crypt==>hex
func HttpReply(w http.ResponseWriter, output interface{}, codec ICodec) {
	if buf, err := json.Marshal(output); nil == err {
		buf = httpEncode(buf, codec)

		if obj, ok := output.(IHttpLogEnable); !ok || obj.LogEnable() {
			Log.Debug("http reply: %s", string(buf))
		}

		w.Write(buf)
	}
}

func StdHttpBody(rBody io.ReadCloser, iBody interface{}) error {
	return HttpBody(rBody, iBody, nil)
}

func StdHttpGet(iGet IHttpGet, output interface{}) error {
	return HttpGet(iGet, output, nil)
}

func StdHttpPost(iPost IHttpPost, input, output interface{}) error {
	return HttpPost(iPost, input, output, nil)
}

func StdHttpReply(w http.ResponseWriter, output interface{}) {
	HttpReply(w, output, nil)
}
