// Copyright 2017 Paul Roy All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Testing azurestoragecache (set, get, delete)
//
// Parameters
// - accountName	: Azure Storage Account Name
// - accountKey		: Azure Storage Account Key
// - containerName	: Azure Storage Container Name /!\ must be lower-case
package main

import (
	"image"
	"image/png"
	"os"
	"fmt"
	"flag"
	"bytes"
	"io/ioutil"
	
	"github.com/PaulARoy/azurestoragecache"
)

func main() {
	// get arguments
	flag.Parse()
	var accountName = flag.String("accountName", "", "Azure Storage Account Name")
	var accountKey = flag.String("accountKey", "", "Azure Storage Account Key")
	var containerName = flag.String("containerName", "cache", "Azure Storage Container Name")

	// create cache
	cache, res, err := azurestoragecache.New(*accountName, *accountKey, *containerName)

	fmt.Println("***** CREATION *****")
	fmt.Println("Container has been created: ", res)
	fmt.Println("Error: ", err)
	fmt.Println("--OK--\n")
	
    // open file
    infile, err := os.Open("in.png")
	handle(err)
    defer infile.Close()

	// read file
	src, _, err := image.Decode(infile)
    handle(err)
	
	// to buffer
	buf := new(bytes.Buffer)
	err = png.Encode(buf, src)
	handle(err)
	
	// upload
	fmt.Println("***** SET *****")
	cache.Set("mykey", buf.Bytes())
	fmt.Println("--OK--\n")
	
	// download
	fmt.Println("***** GET *****")
	bytes, res := cache.Get("mykey")
	if !res {
		fmt.Println("/!\\ ERROR: GOT EMPTY FILE AFTER DOWNLOAD")
	}
	fmt.Println("--OK--\n")
	
    // write file
    err = ioutil.WriteFile("out.png", bytes, os.FileMode(0644))
	handle(err)
	
    // delete key
	fmt.Println("***** DELETE *****")
    cache.Delete("mykey")
	fmt.Println("--OK--\n")
	
    // check we have nothing left
	fmt.Println("***** CHECK-DELETE *****")
    bytes, res = cache.Get("mykey")
	if !res {
		fmt.Println("FILE CORRECTLY DELETED")
	}
	fmt.Println("--OK--\n")
}

func handle(err error) {
    if err != nil {
        fmt.Println("/!\\ ERROR: ", err)
		panic(err)
    }
}