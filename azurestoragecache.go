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

// Package azurestoragecache provides an implementation of httpcache.Cache that
// stores and retrieves data using Azure Storage.
package azurestoragecache // import "github.com/PaulARoy/azurestoragecache"

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	vendorstorage "github.com/Azure/azure-sdk-for-go/storage"
)

// Cache objects store and retrieve data using Azure Storage
type Cache struct {
	// Our configuration for Azure Storage
	Config Config

	// The Azure Blob Storage Client
	Client vendorstorage.BlobStorageClient
}

type Config struct {
	// Account configuration for Azure Storage
	AccountName string
	AccountKey  string

	// Container name to use to store blob
	ContainerName string
}

var noLogErrors, _ = strconv.ParseBool(os.Getenv("NO_LOG_AZUREBSCACHE_ERRORS"))

func (c *Cache) Get(key string) (resp []byte, ok bool) {
	rdr, err := c.Client.GetBlob(c.Config.ContainerName, key)
	if err != nil {
		return []byte{}, false
	}

	resp, err = ioutil.ReadAll(rdr)
	if err != nil {
		if !noLogErrors {
			log.Printf("azurestoragecache.Get failed: %s", err)
		}
	}

	rdr.Close()
	return resp, err == nil
}

func (c *Cache) Set(key string, block []byte) {
	err := c.Client.CreateBlockBlobFromReader(c.Config.ContainerName,
		key,
		uint64(len(block)),
		bytes.NewReader(block),
		nil)
	if err != nil {
		if !noLogErrors {
			log.Printf("azurestoragecache.Set failed: %s", err)
		}
		return
	}
}

func (c *Cache) Delete(key string) {
	res, err := c.Client.DeleteBlobIfExists(c.Config.ContainerName, key, nil)
	if !noLogErrors {
		log.Printf("azurestoragecache.Delete result: %s", res)
	}
	if err != nil {
		if !noLogErrors {
			log.Printf("azurestoragecache.Delete failed: %s", err)
		}
	}
}

// New returns a new Cache with underlying client for Azure Storage
//
// accountName is the Azure Storage Account Name (part of credentials)
// accountKey is the Azure Storage Account Key (part of credentials)
// containerName is the container name in which images will be stored (/!\ LOWER CASE)
//
// The environment variables AZURESTORAGE_ACCOUNT_NAME and AZURESTORAGE_ACCESS_KEY
// are used as credentials if nothing is provided.
func New(accountName string, accountKey string, containerName string) (*Cache, bool, error) {
	accName := accountName
	accKey := accountKey
	contName := containerName

	if len(accName) <= 0 {
		accName = os.Getenv("AZURESTORAGE_ACCOUNT_NAME")
	}

	if len(accKey) <= 0 {
		accKey = os.Getenv("AZURESTORAGE_ACCESS_KEY")
	}

	if len(contName) <= 0 {
		contName = "cache"
	}

	cache := Cache{
		Config: Config{
			AccountName:   accName,
			AccountKey:    accKey,
			ContainerName: contName,
		},
	}

	api, err := vendorstorage.NewBasicClient(cache.Config.AccountName, cache.Config.AccountKey)
	if err != nil {
		return nil, false, err
	}

	cache.Client = api.GetBlobService()

	res, err := cache.Client.CreateContainerIfNotExists(cache.Config.ContainerName, vendorstorage.ContainerAccessTypeBlob)
	if err != nil {
		return nil, false, err
	}

	return &cache, res, nil
}
