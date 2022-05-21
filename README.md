go-ltsv
=======

[![Test Status](https://github.com/Songmu/go-ltsv/workflows/test/badge.svg?branch=main)][actions]
[![Coverage Status](https://codecov.io/gh/Songmu/go-ltsv/branch/main/graph/badge.svg)][codecov]
[![MIT License](https://img.shields.io/github/license/Songmu/go-ltsv)][license]
[![PkgGoDev](https://pkg.go.dev/badge/github.com/Songmu/go-ltsv)][PkgGoDev]

[actions]: https://github.com/Songmu/go-ltsv/actions?workflow=test
[codecov]: https://codecov.io/gh/Songmu/go-ltsv
[license]: https://github.com/Songmu/go-ltsv/blob/main/LICENSE
[PkgGoDev]: https://pkg.go.dev/github.com/Songmu/go-ltsv

LTSV library to map ltsv to struct.

## Synopsis

```go
import (
	"net"

	"github.com/Songmu/go-ltsv"
)

type log struct {
	Host    net.IP
	Req     string
	Status  int
	Size    int
	UA      string
	ReqTime float64
	AppTime *float64
	VHost   string
}

func main() {
	ltsvLog := "time:2016-07-13T00:00:04+09:00\t" +
		"host:192.0.2.1\t" +
		"req:POST /api/v0/tsdb HTTP/1.1\t" +
		"status:200\t" +
		"size:36\t" +
		"ua:ua:mackerel-agent/0.31.2 (Revision 775fad2)\t" +
		"reqtime:0.087\t" +
		"vhost:mackerel.io"
	l := &log{}
	ltsv.Unmarshal([]byte(ltsvLog), l)
	...
}
```

## Description

LTSV parser and encoder for Go with reflection

## Installation

```console
% go get github.com/Songmu/go-ltsv
```

## Author

[Songmu](https://github.com/Songmu)
