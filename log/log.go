package log

import (
	stdLog "log"

	"github.com/davecgh/go-spew/spew"
)

var (
	Fatal   = stdLog.Fatal
	Fatalf  = stdLog.Fatalf
	Fatalln = stdLog.Fatalln
	Panic   = stdLog.Panic
	Panicf  = stdLog.Panicf
	Panicln = stdLog.Panicln
	Prefix  = stdLog.Prefix
	Print   = stdLog.Print
	Printf  = stdLog.Printf
	Println = stdLog.Println
)

func init() {
	stdLog.SetFlags(0)
}

func Dump(v ...interface{}) {
	spew.Dump(v...)
}

func FatalIf(err error) {
	if err != nil {
		stdLog.Fatal(err)
	}
}

func PanicIf(err error) {
	if err != nil {
		stdLog.Panic(err)
	}
}

func PanicIff(err error, template string, v ...interface{}) {
	if err != nil {
		stdLog.Panicf(template, v...)
	}
}

func SetPrefix(s string) {
	stdLog.SetPrefix("[" + s + "] ")
}
