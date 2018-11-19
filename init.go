package asdf

var sizeCheckers = []*SizeChecker{
	scListNode,
	scList,
	scHashNode,
	scUnsafeStimerNode,
	scPortTuple,
	scIp2Tuple,
	scIp4Tuple,
	scIp5Tuple,
	scIp6Tuple,
}

func initEnum() {
	logTypes.Init()
	ipProtos.Init()
	httpMethods.Init()
}

func init() {
	initEnum()
	initCrypt()
}
