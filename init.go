package asdf

var sizeCheckers = []*SizeChecker{
	scListNode,
	scList,
	scHashNode,
	scUnsafeStimerNode,
}

func initEnum() {
	logTypes.Init()
}

func init() {
	initEnum()
	initCrypt()
}
