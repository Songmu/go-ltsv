package ltsv_test

import (
	"fmt"
	"net"
	"time"

	"github.com/Songmu/go-ltsv"
	"github.com/kr/pretty"
)

type log struct {
	Time    *logTime
	Host    net.IP
	Req     string
	Status  int
	Size    int
	UA      string
	ReqTime float64
	AppTime *float64
	VHost   string
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

type logTime struct {
	time.Time
}

func (lt *logTime) UnmarshalText(t []byte) error {
	ti, err := time.ParseInLocation(timeFormat, string(t), time.UTC)
	if err != nil {
		return err
	}
	lt.Time = ti
	return nil
}

func ExampleUnmarshal() {
	ltsvLog := "time:2016-07-13T00:00:04+09:00\t" +
		"host:192.0.2.1\t" +
		"req:POST /api/v0/tsdb HTTP/1.1\t" +
		"status:200\t" +
		"size:36\t" +
		"ua:ua:mackerel-agent/0.31.2 (Revision 775fad2)\t" +
		"reqtime:0.087\t" +
		"vhost:mackerel.io"
	l := log{}
	ltsv.Unmarshal([]byte(ltsvLog), &l)
	t := l.Time
	l.Time = nil
	fmt.Println(t)
	pretty.Println(l)
	// Output:
	// 2016-07-13 00:00:04 +0900 +0900
	// ltsv_test.log{
	//     Time:    (*ltsv_test.logTime)(nil),
	//     Host:    {0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xc0, 0x0, 0x2, 0x1},
	//     Req:     "POST /api/v0/tsdb HTTP/1.1",
	//     Status:  200,
	//     Size:    36,
	//     UA:      "ua:mackerel-agent/0.31.2 (Revision 775fad2)",
	//     ReqTime: 0.087,
	//     AppTime: (*float64)(nil),
	//     VHost:   "mackerel.io",
	// }
}
