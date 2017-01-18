package asdf

const InvalidEnum = -1

type IEnum interface {
	INumber
	IGood
	IToString
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
