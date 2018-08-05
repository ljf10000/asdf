package asdf

import (
	"fmt"
)

type Capacity struct {
	Count uint32
	Cap   uint32
}

func (me *Capacity) String() string {
	return fmt.Sprintf("cap(%d) count(%d)", me.Cap, me.Count)
}

func (me *Capacity) Zero() {
	me.Count = 0
	me.Cap = 0
}

func (me *Capacity) IsFull() bool {
	return me.Count == me.Cap
}

type Capacity32 struct {
	Size uint32
	Cap  uint32
}

func (me *Capacity32) String() string {
	return fmt.Sprintf("cap(%d) size(%d)", me.Cap, me.Size)
}

func (me *Capacity32) Zero() {
	me.Size = 0
	me.Cap = 0
}

func (me *Capacity32) IsFull() bool {
	return me.Size == me.Cap
}

func (me *Capacity32) AddAlign(v uint32, align uint32) {
	me.Size += v
	me.Cap += Align32(v, align)
}

type Capacity64 struct {
	Size uint64
	Cap  uint64
}

func (me *Capacity64) String() string {
	return fmt.Sprintf("cap(%d) size(%d)", me.Cap, me.Size)
}

func (me *Capacity64) Zero() {
	me.Size = 0
	me.Cap = 0
}

func (me *Capacity64) IsFull() bool {
	return me.Size == me.Cap
}

func (me *Capacity64) Add32(v Capacity32) {
	me.Size += uint64(v.Size)
	me.Cap += uint64(v.Cap)
}

func (me *Capacity64) Sub32(v Capacity32) {
	me.Size -= uint64(v.Size)
	me.Cap -= uint64(v.Cap)
}
