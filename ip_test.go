package asdf

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestIp(t *testing.T) {
	ip := IpAddress(0)
	bin := [4]byte{}

	ipstring := "192.168.0.1"
	ip.FromString(ipstring)

	binary.BigEndian.PutUint32(bin[:], uint32(ip))

	fmt.Printf("ipstring=%s, ip=%x, bin=%v"+Crlf, ipstring, uint32(ip), bin[:])
}

func TestIpTuple(t *testing.T) {
	const COUNT = 3000

	tuples := [COUNT * COUNT]Ip2Tuple{}
	maps := map[Ip2Tuple]uint64{}

	for i := 0; i < COUNT; i++ {
		for j := 0; j < COUNT; j++ {
			tuple := Ip2Tuple{
				Sip: IpAddress(i),
				Dip: IpAddress(j),
			}

			tuples[i*COUNT+j] = tuple

			maps[tuple] = 1
		}
	}

	hits := uint64(0)

	old := NowTime64()

	for i := 0; i < COUNT*COUNT; i++ {
		tuple := tuples[i]

		if _, ok := maps[tuple]; ok {
			hits++
		}
	}

	now := NowTime64()
	ns := uint64(now - old)

	t.Logf("hits: %d, time: %dms, pps: %d\n", hits, ns/1000000, hits*1000000000/ns)
}
