package asdf

const (
	MAX_INT   = int(^uint(0) >> 1)
	MAX_INT32 = int32(^uint32(0) >> 1)
	MAX_INT64 = int64(^uint64(0) >> 1)

	MIN_INT   = ^MAX_INT
	MIN_INT32 = ^MAX_INT32
	MIN_INT64 = ^MAX_INT64

	MAX_UINT   = ^uint(0)
	MAX_UINT32 = ^uint32(0)
	MAX_UINT64 = ^uint64(0)
)
