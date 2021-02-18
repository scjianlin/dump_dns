package main

import (
	"dns/pkg"
	"flag"
)

var env = flag.String("env", "", "`env` is environment value.")

func main() {
//	go func() {
//		fmt.Println(http.ListenAndServe("0.0.0.0:8080", nil))
//	}()
	flag.Parse()
	pkg.Start(*env)
}
