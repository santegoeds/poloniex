// +build debug

package trace

import "github.com/kr/pretty"

func Printf(fmt string, args ...interface{}) {
	pretty.Printf("TRACE: "+fmt, args...)
}

func Println(args ...interface{}) {
	pretty.Println(args...)
}
