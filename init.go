package asdf

var sizeCheckers = []*SizeChecker{
	scListNode,
	scList,
	scHashNode,
	scUnsafeStimerNode,
	scIp2Tuple,
}

func initEnum() {
	logTypes.Init()
}

func init() {
	initEnum()
	initCrypt()
}
