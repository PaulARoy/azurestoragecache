// Testing azurestoragecache
package main

import (
	"flag"
	"sourcegraph.com/PaulARoy/azurestoragecache"
)

func main() {
	flag.Parse()
	azurestoragecache.New(nil, nil, "Cache")
}