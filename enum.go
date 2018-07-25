package asdf

const InvalidEnum = -1

type IEnum interface {
	INumber
	IGood
	IString
	// todo: IFromString
}

func IsGoodEnum(idx interface{}) bool {
	n, ok := idx.(INumber)
	if !ok {
		return false
	}
	v := n.Int()

	return v >= n.Begin() && v < n.End()
}

type EnumBinding []string

// todo: reutrn string and error
func (me EnumBinding) EntryShow(idx interface{}) string {
	if nil == me {
		return Empty
	}

	e, ok := idx.(IEnum)
	if !ok {
		return Empty
	}

	if !e.IsGood() {
		return Empty
	}

	return me[e.Int()]
}

type EnumManager struct {
	iMap map[int]string
	sMap map[string]int
}

func EnumManagerCreate() *EnumManager {
	return &EnumManager{
		iMap: make(map[int]string),
		sMap: make(map[string]int),
	}
}

func (me *EnumManager) Register(idx int, name string) {
	me.iMap[idx] = name
	me.sMap[name] = idx
}

func (me *EnumManager) Index(name string) (int, bool) {
	idx, ok := me.sMap[name]

	return idx, ok
}

func (me *EnumManager) Index2(name string) int {
	idx, ok := me.sMap[name]
	if ok {
		return idx
	} else {
		return InvalidEnum
	}
}

func (me *EnumManager) Name(idx int) (string, bool) {
	name, ok := me.iMap[idx]

	return name, ok
}

func (me *EnumManager) Name2(idx int) string {
	name, ok := me.iMap[idx]
	if ok {
		return name
	} else {
		return Unknow
	}
}

func (me *EnumManager) IsGoodName(name string) bool {
	_, ok := me.sMap[name]

	return ok
}

func (me *EnumManager) IsGoodIndex(idx int) bool {
	_, ok := me.iMap[idx]

	return ok
}

type EnumMapper struct {
	Enum   string
	Names  []string
	values map[string]int
}

func (me *EnumMapper) Init() {
	me.values = map[string]int{}

	for k, v := range me.Names {
		me.values[v] = k
	}
}

func (me *EnumMapper) Index(name string) (int, error) {
	idx, ok := me.values[name]
	if ok {
		return idx, nil
	} else {
		return 0, ErrSprintf("invalid %s: %s", me.Enum, name)
	}
}

func (me *EnumMapper) Name(idx int) string {
	if me.IsGoodIndex(idx) {
		return me.Names[idx]
	} else {
		return Unknow
	}
}

func (me *EnumMapper) NameEx(idx int) (string, bool) {
	if me.IsGoodIndex(idx) {
		return me.Names[idx], true
	} else {
		return Unknow, false
	}
}

func (me *EnumMapper) IsGoodIndex(idx int) bool {
	return idx >= 0 && idx < len(me.Names)
}
