package asdf

const (
	IPPROTO_IP      = 0
	IPPROTO_ICMP    = 1
	IPPROTO_IGMP    = 2
	IPPROTO_IPIP    = 4
	IPPROTO_TCP     = 6
	IPPROTO_EGP     = 8
	IPPROTO_PUP     = 12
	IPPROTO_UDP     = 17
	IPPROTO_IDP     = 22
	IPPROTO_TP      = 29
	IPPROTO_DCCP    = 33
	IPPROTO_IPV6    = 41
	IPPROTO_RSVP    = 46
	IPPROTO_GRE     = 47
	IPPROTO_ESP     = 50
	IPPROTO_AH      = 51
	IPPROTO_MTP     = 92
	IPPROTO_BEETPH  = 94
	IPPROTO_ENCAP   = 98
	IPPROTO_PIM     = 103
	IPPROTO_COMP    = 108
	IPPROTO_SCTP    = 132
	IPPROTO_UDPLITE = 136
	IPPROTO_RAW     = 255
	IPPROTO_END     = 256
)

type IpProto byte

var ipProtos = EnumMapper{
	Type: "asdf.IpProto",
	Names: []string{
		IPPROTO_IP:      "ip",
		IPPROTO_ICMP:    "icmp",
		IPPROTO_IGMP:    "igmp",
		IPPROTO_IPIP:    "ipip",
		IPPROTO_TCP:     "tcp",
		IPPROTO_EGP:     "egp",
		IPPROTO_PUP:     "pup",
		IPPROTO_UDP:     "udp",
		IPPROTO_IDP:     "idp",
		IPPROTO_TP:      "tp",
		IPPROTO_DCCP:    "dccp",
		IPPROTO_IPV6:    "ipv6",
		IPPROTO_RSVP:    "rsvp",
		IPPROTO_GRE:     "gre",
		IPPROTO_ESP:     "esp",
		IPPROTO_AH:      "ah",
		IPPROTO_MTP:     "mtp",
		IPPROTO_BEETPH:  "beetph",
		IPPROTO_ENCAP:   "encap",
		IPPROTO_PIM:     "pim",
		IPPROTO_COMP:    "comp",
		IPPROTO_SCTP:    "sctp",
		IPPROTO_UDPLITE: "udplite",
		IPPROTO_RAW:     "raw",
	},
}

func (me IpProto) HavePort() bool {
	switch me {
	case IPPROTO_TCP, IPPROTO_UDP:
		return true
	default:
		return false
	}
}

func (me IpProto) IsGood() bool {
	return ipProtos.IsGoodIndex(int(me))
}

func (me IpProto) String() string {
	return ipProtos.Name(int(me))
}

func (me IpProto) Compare(obj IpProto) int {
	return CompareByte(byte(me), byte(obj))
}

func (me *IpProto) FromString(s string) error {
	idx, err := ipProtos.Index(s)
	if nil == err {
		*me = IpProto(idx)
	}

	return err
}
