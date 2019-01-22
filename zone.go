package asdf

type IZoneNode interface {
	String() string
	IsGood() bool
	Zero()

	Up() IZoneNode
	Down() IZoneNode

	Compare(v IZoneNode) int
	InRange(a, b IZoneNode) bool
	InZone(z IZone) bool
	Add(v IZoneNode) IZoneNode
}

type IZoneEx interface {
	NodeBegin() IZoneNode
	NodeEnd() IZoneNode

	String() string
	IsGood() bool
	Zero()
	Include(v IZone) bool
	Match(v IZone) bool
	Intersect(v IZone) (IZone, bool)
}

type IZone interface{}
