// Testing azurestoragecache
package main

import (
	"flag"
	"azurestoragecache"
)

func main() {
	flag.Parse()
	azurestoragecache.New(nil, nil, "Cache")
}