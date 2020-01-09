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
	scFlatString,
}

func initEnum() {
	evTypes.Init()
	logTypes.Init()
	ipProtos.Init()
	frequencies.Init()
	httpMethods.Init()
	configFileTypes.Init()
}

func init() {
	initEnum()
	initCrypt()
}
