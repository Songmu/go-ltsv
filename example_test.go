package ltsv_test

import (
	"net"
	"time"

	"github.com/Songmu/go-ltsv"
	"github.com/kr/pretty"
)

type log struct {
	Time    LogTime
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

type LogTime struct {
	time.Time
}

func (lt *LogTime) UnmarshalText(t []byte) error {
	ti, err := time.Parse(timeFormat, string(t))
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
	pretty.Println(l)
	// Output:
	// ltsv_test.log{
	//     Time: ltsv_test.LogTime{
	//         Time: time.Time{
	//             sec:  63603932404,
	//             nsec: 0,
	//             loc:  &time.Location{
	//                 name: "Local",
	//                 zone: {
	//                     {name:"JCST", offset:32400, isDST:false},
	//                     {name:"JDT", offset:36000, isDST:true},
	//                     {name:"JST", offset:32400, isDST:false},
	//                 },
	//                 tx: {
	//                     {when:-1017824400, index:0x2, isstd:false, isutc:false},
	//                     {when:-683794800, index:0x1, isstd:false, isutc:false},
	//                     {when:-672393600, index:0x2, isstd:false, isutc:false},
	//                     {when:-654764400, index:0x1, isstd:false, isutc:false},
	//                     {when:-640944000, index:0x2, isstd:false, isutc:false},
	//                     {when:-620290800, index:0x1, isstd:false, isutc:false},
	//                     {when:-609494400, index:0x2, isstd:false, isutc:false},
	//                     {when:-588841200, index:0x1, isstd:false, isutc:false},
	//                     {when:-578044800, index:0x2, isstd:false, isutc:false},
	//                 },
	//                 cacheStart: -578044800,
	//                 cacheEnd:   9223372036854775807,
	//                 cacheZone:  &time.zone{(CYCLIC REFERENCE)},
	//             },
	//         },
	//     },
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
