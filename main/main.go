// Testing azurestoragecache
package main

import (
	"flag"
	"github.com/PaulARoy/azurestoragecache"
)

func main() {
	flag.Parse()
	azurestoragecache.New("", "", "Cache")
}