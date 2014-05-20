/*
+The MIT License (MIT)
 +
 +Copyright (c) 2014 Levi Cook
 +
 +Permission is hereby granted, free of charge, to any person obtaining a copy
 +of this software and associated documentation files (the "Software"), to deal
 +in the Software without restriction, including without limitation the rights
 +to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 +copies of the Software, and to permit persons to whom the Software is
 +furnished to do so, subject to the following conditions:
 +
 +The above copyright notice and this permission notice shall be included in all
 +copies or substantial portions of the Software.
 +
 +THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 +IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 +FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 +AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 +LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 +OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 +SOFTWARE.

*/

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
