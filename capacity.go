package asdf

import (
	"fmt"
)

type Capacity struct {
	Count uint32 `json:"count"`
	Cap   uint32 `json:"cap"`
}

func (me *Capacity) String() string {
	return fmt.Sprintf("cap(%d) count(%d)", me.Cap, me.Count)
}

func (me *Capacity) Zero() {
	me.Count = 0
	me.Cap = 0
}

func (me *Capacity) IsZero() bool {
	return 0 == me.Count || 0 == me.Cap
}

func (me *Capacity) IsFull() bool {
	return me.Count == me.Cap
}

/******************************************************************************/

type BlockCap32 struct {
	Size uint32 `json:"size"`
	Cap  uint32 `json:"cap"`
}

func (me *BlockCap32) String() string {
	return fmt.Sprintf("cap(%d) size(%d)", me.Cap, me.Size)
}

func (me *BlockCap32) Zero() {
	me.Size = 0
	me.Cap = 0
}

func (me *BlockCap32) IsFull() bool {
	return me.Size == me.Cap
}

func (me *BlockCap32) IsZero() bool {
	return 0 == me.Size && 0 == me.Cap
}

func (me *BlockCap32) AddAlign(v uint32, align uint32) {
	me.Size += v
	me.Cap += Align32(v, align)
}

/******************************************************************************/

type BlockCap64 struct {
	Size uint64 `json:"size"`
	Cap  uint64 `json:"cap"`
}

func (me *BlockCap64) String() string {
	return fmt.Sprintf("cap(%d) size(%d)", me.Cap, me.Size)
}

func (me *BlockCap64) Zero() {
	me.Size = 0
	me.Cap = 0
}

func (me *BlockCap64) IsFull() bool {
	return me.Size == me.Cap
}

func (me *BlockCap64) IsZero() bool {
	return 0 == me.Size && 0 == me.Cap
}

func (me *BlockCap64) Add32(v BlockCap32) {
	me.Size += uint64(v.Size)
	me.Cap += uint64(v.Cap)
}

func (me *BlockCap64) Sub32(v BlockCap32) {
	me.Size -= uint64(v.Size)
	me.Cap -= uint64(v.Cap)
}

/******************************************************************************/
//
type SizeCountStat struct {
	Size  uint64 `json:"size"`
	Count uint64 `json:"count"`
}

func (me *SizeCountStat) String() string {
	return "size: " + Utoa64(me.Size) +
		", count: " + Utoa64(me.Count)
}

func (me *SizeCountStat) Add(size, count uint64) {
	me.Size += size
	me.Count += count
}

//
type SizeCountStat32 struct {
	Size  uint32 `json:"size"`
	Count uint32 `json:"count"`
}

func (me *SizeCountStat32) String() string {
	return "size: " + Utoa32(me.Size) +
		", count: " + Utoa32(me.Count)
}

func (me *SizeCountStat32) Add(size, count uint32) {
	me.Size += size
	me.Count += count
}
