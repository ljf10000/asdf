package asdf

import (
	"unsafe"
)

const (
	MPARRAY_STR_BORDER = "\x12\x34\x56\x78\x9a\xbc\xde\xff"
	MPARRAY_PATH       = "/tmp"
	MPARRAY_FILE_META  = MPARRAY_PATH + "/mparray.meta.shm"
	MPARRAY_FILE_BMAP  = MPARRAY_PATH + "/mparray.bmap.shm"
	MPARRAY_FILE_DATA  = MPARRAY_PATH + "/mparray.data.shm"
)

const (
	MPARRAY_BORDER_SIZE = SizeofPointer
	MPARRAY_VERSION     = 0
	MPARRAY_TIMEOUT     = 3

	MPARRAY_ZONE_END      = 32 * SizeofK
	MPARRAY_ENTRY_MAXSIZE = SizeofM

	MPARRAY_READER_MINSIZE = 4
	MPARRAY_READER_MAXSIZE = 8192
)

const (
	MPARRAY_OBJ_META = iota
	MPARRAY_OBJ_BMAP
	MPARRAY_OBJ_DATA

	MPARRAY_OBJ_END
)

var MPARRAY_BORDER [MPARRAY_BORDER_SIZE]byte

func init() {
	copy(MPARRAY_BORDER[:], ([]byte)(MPARRAY_STR_BORDER[:MPARRAY_BORDER_SIZE]))
}

/******************************************************************************/
type mparray_border_t struct {
	body [MPARRAY_BORDER_SIZE]byte
}

func (me *mparray_border_t) init() {
	copy(me.body[:], MPARRAY_BORDER[:])
}

func (me *mparray_border_t) is_good() bool {
	return me.body == MPARRAY_BORDER
}

/******************************************************************************/
var __mparray_entry_hdr mparray_entry_hdr

const sizeof_mparray_entry_hdr = unsafe.Sizeof(__mparray_entry_hdr)

type mparray_entry_hdr struct {
	head_border mparray_border_t
	lock        MpLockTm
	body_border mparray_border_t
}

type mparray_entry_t struct {
	mparray_entry_hdr

	body [MPARRAY_ENTRY_MAXSIZE]byte
}

func MPARRAY_ENTRY_BODYSIZE(size uint32) uint32 {
	return Align32(size, SizeofPointer)
}

func MPARRAY_ENTRY_REALSIZE(size uint32) uint32 {
	return uint32(sizeof_mparray_entry_hdr) + MPARRAY_ENTRY_BODYSIZE(size) + MPARRAY_BORDER_SIZE
}

func MPARRAY_ENTRY_SIZE(size uint32) uint32 {
	return Align32(MPARRAY_ENTRY_REALSIZE(size), SizeofCacheLine)
}

func MPARRAY_ENTRY_OFFSET(idx uint32, size uint32) uint64 {
	return uint64(idx) * uint64(MPARRAY_ENTRY_SIZE(size))
}

func (me *mparray_entry_t) slice(size uint32) []byte {
	total := int(MPARRAY_ENTRY_REALSIZE(size))

	return MakeSlice(unsafe.Pointer(me), total, total)
}

func (me *mparray_entry_t) entry(size uint32) []byte {
	return me.body[:size]
}

func (me *mparray_entry_t) foot_border(size uint32) *mparray_border_t {
	address := uintptr(unsafe.Pointer(me))
	foot_address := address + uintptr(sizeof_mparray_entry_hdr) + uintptr(MPARRAY_ENTRY_BODYSIZE(size))

	return (*mparray_border_t)(unsafe.Pointer(foot_address))
}

func (me *mparray_entry_t) is_good(size uint32) bool {
	if false == me.head_border.is_good() {
		Panic("invalid head border")
	} else if false == me.body_border.is_good() {
		Panic("invalid body border")
	} else if false == me.foot_border(size).is_good() {
		Panic("invalid foot border")
	}

	return true
}

func (me *mparray_entry_t) init(size uint32) {
	me.head_border.init()
	me.body_border.init()
	me.foot_border(size).init()
}

/******************************************************************************/
type mparray_bset_t = BitSet

type mparray_data_t struct {
	_ [0]byte
}

func (me *mparray_data_t) zone(desc *mparray_desc_t) *mparray_data_t {
	offset := desc.offset[MPARRAY_OBJ_DATA]
	address := uintptr(unsafe.Pointer(me)) + uintptr(offset)

	return (*mparray_data_t)(unsafe.Pointer(address))
}

func (me *mparray_data_t) entry(desc *mparray_desc_t, idx MpArrayIndex) *mparray_entry_t {
	offset := MPARRAY_ENTRY_OFFSET(idx.Entry, desc.config.Size)
	address := uintptr(unsafe.Pointer(me)) + uintptr(offset)

	return (*mparray_entry_t)(unsafe.Pointer(address))
}

/******************************************************************************/
type mparray_bmap_t struct {
	_ [0]byte
}

func (me *mparray_bmap_t) zone(desc *mparray_desc_t) *mparray_bset_t {
	offset := desc.offset[MPARRAY_OBJ_BMAP]
	address := uintptr(unsafe.Pointer(me)) + uintptr(offset)

	return (*mparray_bset_t)(unsafe.Pointer(address))
}

func MPARRAY_BMAP_SIZE(count uint32) uint32 {
	return Align32(__BITSET_SIZE(count), SizeofCacheLine)
}

/******************************************************************************/
var __mparray_desc_t mparray_desc_t

const sizeof_mparray_desc_t = unsafe.Sizeof(__mparray_desc_t)

type mparray_desc_t struct {
	config MpArrayConfig

	offset [MPARRAY_OBJ_END]uint64
}

func (me *mparray_desc_t) is_good_index(idx uint32) bool {
	return idx < me.config.Count
}

/******************************************************************************/
func mparray_calc(config []MpArrayConfig, offset [][MPARRAY_OBJ_END]uint64, size []uint64) {
	size[MPARRAY_OBJ_META] = uint64(sizeof_mparray_meta_hdr)
	size[MPARRAY_OBJ_BMAP] = 0
	size[MPARRAY_OBJ_DATA] = 0

	count := len(config)
	for i := 0; i < count; i++ {
		for j := 0; j < MPARRAY_OBJ_END; j++ {
			offset[i][j] = size[j]

			Log.Debug("offset[%d][%d]=%d", i, j, offset[i][j])
		}

		size[MPARRAY_OBJ_META] += uint64(sizeof_mparray_desc_t)
		size[MPARRAY_OBJ_BMAP] += uint64(MPARRAY_BMAP_SIZE(config[i].Count))
		size[MPARRAY_OBJ_DATA] += MPARRAY_ENTRY_OFFSET(config[i].Count, config[i].Size)
	}
}

/******************************************************************************/
var __mparray_meta_hdr mparray_meta_hdr

const sizeof_mparray_meta_hdr = unsafe.Sizeof(__mparray_meta_hdr)

type mparray_meta_hdr struct {
	count uint32
	_     uint32

	size [MPARRAY_OBJ_END]uint64
}

type mparray_meta_t struct {
	mparray_meta_hdr

	desc [MPARRAY_ZONE_END]mparray_desc_t
}

func (me *mparray_meta_t) is_good_index(idx uint32) bool {
	return idx < me.count
}

/******************************************************************************/

type MpArrayConfig struct {
	Count uint32
	Size  uint32
}

type MpArrayIndex struct {
	Array uint32
	Entry uint32
}

type MpArrayCursor struct {
	desc *mparray_desc_t
	bset *mparray_bset_t
	data *mparray_data_t
}

type MpArrayReader struct {
	Ok bool

	Body []byte
}

type MpArray struct {
	File [MPARRAY_OBJ_END]string
	mem  [MPARRAY_OBJ_END]MMap

	meta *mparray_meta_t
	bmap *mparray_bmap_t
	data *mparray_data_t

	Trace  bool
	locked bool
	idx    MpArrayIndex
}

func (me *MpArray) trace_start(idx MpArrayIndex) {
	if me.Trace {
		me.locked = true
		me.idx = idx
	}
}

func (me *MpArray) trace_stop(idx MpArrayIndex) {
	if me.Trace {
		me.locked = false
	}
}

func (me *MpArray) trace(idx MpArrayIndex, cb func()) {
	me.trace_start(idx)
	cb()
	me.trace_stop(idx)
}

func (me *MpArray) cursor(idx MpArrayIndex, cursor *MpArrayCursor) error {
	if false == me.meta.is_good_index(idx.Array) {
		return ErrLog("invalid mparray index:%d", idx.Array)
	}

	desc := &me.meta.desc[idx.Array]
	if false == desc.is_good_index(idx.Entry) {
		return ErrLog("invalid mparray entry index:%d", idx.Array)
	}

	cursor.desc = desc
	cursor.bset = me.bmap.zone(desc)
	cursor.data = me.data.zone(desc)

	return nil
}

func (me *MpArray) bind(idx int, mem MMap) {
	switch idx {
	case MPARRAY_OBJ_META:
		me.meta = (*mparray_meta_t)(SlicePointer(mem))
	case MPARRAY_OBJ_BMAP:
		me.bmap = (*mparray_bmap_t)(SlicePointer(mem))
	case MPARRAY_OBJ_DATA:
		me.data = (*mparray_data_t)(SlicePointer(mem))
	default:
		Panic("invalid zone[%d]", idx)
	}

	me.mem[idx] = mem
}

type MpArrayEachZone func(array *MpArray, idx int, cursor *MpArrayCursor) (bool, error)

func mparray_init_zone(array *MpArray, idx int, cursor *MpArrayCursor) (bool, error) {
	cursor.bset.Init(cursor.desc.config.Count)

	Log.Debug("mparray init zone[%d] count[%d]", idx, cursor.desc.config.Count)

	return true, nil
}

type MpArrayEachEntry func(array *MpArray, idx MpArrayIndex, cursor *MpArrayCursor) (bool, error)

func mparray_init_entry(array *MpArray, idx MpArrayIndex, cursor *MpArrayCursor) (bool, error) {
	entry := cursor.data.entry(cursor.desc, idx)
	entry.lock.Init()
	entry.init(cursor.desc.config.Size)

	// Log.Debug("mparray init zone[%d] entry[%d]\n%s", idx.Array, idx.Entry, BinSprintf(entry.slice(cursor.desc.config.Size)))

	return true, nil
}

func (me *MpArray) init() error {
	if err := me.EachZone(mparray_init_zone); nil != err {
		return err
	}

	if err := me.EachEntry(mparray_init_entry); nil != err {
		return err
	}

	return nil
}

/******************************************************************************/

func (me *MpArray) Init() {
	if Empty == me.File[MPARRAY_OBJ_META] {
		me.File[MPARRAY_OBJ_META] = MPARRAY_FILE_META
	}

	if Empty == me.File[MPARRAY_OBJ_BMAP] {
		me.File[MPARRAY_OBJ_BMAP] = MPARRAY_FILE_BMAP
	}

	if Empty == me.File[MPARRAY_OBJ_DATA] {
		me.File[MPARRAY_OBJ_DATA] = MPARRAY_FILE_DATA
	}
}

func (me *MpArray) EachZone(each MpArrayEachZone) error {
	cursor := MpArrayCursor{}
	count := int(me.meta.count)

	for i := 0; i < count; i++ {
		desc := &me.meta.desc[i]

		cursor.desc = desc
		cursor.bset = me.bmap.zone(desc)

		if ok, err := each(me, i, &cursor); false == ok {
			return err
		}
	}

	return nil
}

func (me *MpArray) EachEntry(each MpArrayEachEntry) error {
	cursor := MpArrayCursor{}
	count := int(me.meta.count)

	for i := 0; i < count; i++ {
		desc := &me.meta.desc[i]
		count_entry := int(desc.config.Count)

		for j := 0; j < count_entry; j++ {
			idx := MpArrayIndex{uint32(i), uint32(j)}

			if err := me.cursor(idx, &cursor); nil != err {
				return err
			}

			if ok, err := each(me, idx, &cursor); false == ok {
				return err
			}
		}
	}

	return nil
}

func (me *MpArray) Close() error {
	for i := 0; i < MPARRAY_OBJ_END; i++ {
		me.mem[i].Unmap()
	}

	return nil
}

func (me *MpArray) Open(config []MpArrayConfig) error {
	count := len(config)
	offset := make([][MPARRAY_OBJ_END]uint64, count)
	size := [MPARRAY_OBJ_END]uint64{}
	exist := [MPARRAY_OBJ_END]bool{}

	for i := 0; i < MPARRAY_OBJ_END; i++ {
		exist[i] = FileName(me.File[i]).Exist()
	}

	for i := 1; i < MPARRAY_OBJ_END; i++ {
		if exist[i] != exist[i-1] {
			return ErrLog("not all file exist")
		}
	}

	mparray_calc(config, offset[:], size[:])
	if exist[0] {
		// open exist
		for i := 0; i < MPARRAY_OBJ_END; i++ {
			fsize, _ := FileName(me.File[i]).FileSize()
			csize := PageAlign(int(size[i]))

			if int(fsize) != csize {
				return ErrLog("file %s fsize[%d] != csize[%d]", me.File[i], fsize, csize)
			}

			mem, err := MapOpen(me.File[i], 0)
			if nil != err {
				return err
			}
			me.bind(i, mem)

			if me.meta.size[i] != size[i] {
				return ErrLog("file %s meta size[%d] != size[%d]", me.File[i], me.meta.size[i], size[i])
			}
		}

		if me.meta.count != uint32(count) {
			return ErrLog("file meta count not match")
		}
	} else {
		// create new
		for i := 0; i < MPARRAY_OBJ_END; i++ {
			mem, err := MapCreate(me.File[i], int(size[i]), 0)
			if nil != err {
				return err
			}
			me.bind(i, mem)

			mem.Zero()
		}

		me.meta.count = uint32(count)
		me.meta.size = size

		for i := 0; i < count; i++ {
			desc := &me.meta.desc[i]

			desc.offset = offset[i]
			desc.config = config[i]
		}

		if err := me.init(); nil != err {
			return err
		}
	}

	return nil
}

func (me *MpArray) Load() error {
	for i := 0; i < MPARRAY_OBJ_END; i++ {
		mem, err := MapOpen(me.File[i], 0)
		if nil != err {
			return err
		}

		me.bind(i, mem)
	}

	return nil
}

type IMpArrayObject interface {
	ISlice
}

func (me *MpArray) Write(idx MpArrayIndex, obj IMpArrayObject, timeout int) error {
	cursor := MpArrayCursor{}
	mem := obj.Slice()
	size := len(mem)

	Log.Debug("mparray write: zone[%d] entry[%d] size[%d] ...", idx.Array, idx.Entry, size)

	if me.locked {
		return ErrLog("mparray write: have locked")
	}

	if err := me.cursor(idx, &cursor); nil != err {
		return err
	}

	if size > int(cursor.desc.config.Size) {
		return ErrLog("mparray write: too long buffer")
	}

	entry := cursor.data.entry(cursor.desc, idx)
	if false == entry.is_good(cursor.desc.config.Size) {
		Panic("mparray write: invalid entry border")
	}
	// Log.Debug("mparray write:\n%s", BinSprintf(entry.slice(cursor.desc.config.Size)))
	Log.Debug("mparray check entry border: zone[%d] entry[%d] size[%d]", idx.Array, idx.Entry, size)

	state := entry.lock.HandleW(timeout, func() {
		copy(entry.body[:size], mem)

		cursor.bset.Set(idx.Entry)
	})

	if MPLOCK_TIMEOUT == state {
		return ErrLog("mparray write: timeout")
	}

	if false == entry.is_good(cursor.desc.config.Size) {
		Panic("mparray write: invalid entry border")
	}
	Log.Debug("mparray check entry border: zone[%d] entry[%d] size[%d]", idx.Array, idx.Entry, size)

	Log.Debug("mparray write: zone[%d] entry[%d] size[%d] ok.\n", idx.Array, idx.Entry, size)

	return nil
}

func (me *MpArray) Read(idx MpArrayIndex, obj IMpArrayObject, timeout int) (bool, error) {
	cursor := MpArrayCursor{}
	mem := obj.Slice()
	ok := false

	Log.Debug("mparray read: zone[%d] entry[%d] ...", idx.Array, idx.Entry)

	if me.locked {
		return false, ErrLog("mparray read: have locked")
	}

	if err := me.cursor(idx, &cursor); nil != err {
		return false, err
	}

	size := cursor.desc.config.Size
	if uint32(len(mem)) < size {
		return false, ErrLog("mparray read: no enough buffer")
	}

	entry := cursor.data.entry(cursor.desc, idx)
	if false == entry.is_good(size) {
		Panic("mparray read: invalid entry border")
	}
	Log.Debug("mparray check entry border: zone[%d] entry[%d] size[%d]", idx.Array, idx.Entry, size)

	state := entry.lock.HandleR(timeout, func() {
		ok = cursor.bset.IsSet(idx.Entry)
		if ok {
			copy(mem, entry.body[:size])
		}
	})

	if MPLOCK_TIMEOUT == state {
		return false, ErrLog("mparray read: timeout")
	}

	if false == entry.is_good(cursor.desc.config.Size) {
		Panic("mparray read: invalid entry border")
	}
	Log.Debug("mparray check entry border: zone[%d] entry[%d] size[%d]", idx.Array, idx.Entry, size)

	Log.Debug("mparray read: zone[%d] entry[%d] size[%d] ok.\n", idx.Array, idx.Entry, size)

	return ok, nil
}

/******************************************************************************/
