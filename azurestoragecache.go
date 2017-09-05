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

	"github.com/Azure/azure-sdk-for-go/storage"
)

// Cache stores and retrieves data using Azure Storage.
type Cache struct {
	// The Azure Blob Storage Client
	client storage.BlobStorageClient

	// container name to use to store blobs
	container string
}

var noLogErrors, _ = strconv.ParseBool(os.Getenv("NO_LOG_AZUREBSCACHE_ERRORS"))

// Get the cached value with the specified key.
func (c *Cache) Get(key string) (resp []byte, ok bool) {
	rdr, err := c.client.GetBlob(c.container, key)
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

// Set the cached value with the specified key.
func (c *Cache) Set(key string, value []byte) {
	err := c.client.CreateBlockBlobFromReader(c.container,
		key,
		uint64(len(value)),
		bytes.NewReader(value),
		nil)
	if err != nil {
		if !noLogErrors {
			log.Printf("azurestoragecache.Set failed: %s", err)
		}
		return
	}
}

// Delete the cached value with the specified key.
func (c *Cache) Delete(key string) {
	res, err := c.client.DeleteBlobIfExists(c.container, key, nil)
	if !noLogErrors {
		log.Printf("azurestoragecache.Delete result: %s", res)
	}
	if err != nil {
		if !noLogErrors {
			log.Printf("azurestoragecache.Delete failed: %s", err)
		}
	}
}

// New returns a new Cache with underlying client for Azure Storage.
//
// accountName and accountKey are the Azure Storage credentials.  If either are
// empty, the contents of the environment variables AZURESTORAGE_ACCOUNT_NAME
// and AZURESTORAGE_ACCESS_KEY will be used.
//
// containerName is the container name in which cached values will be stored.
// If not specified, "cache" will be used.
func New(accountName string, accountKey string, containerName string) (*Cache, error) {
	if accountName == "" {
		accountName = os.Getenv("AZURESTORAGE_ACCOUNT_NAME")
	}

	if accountKey == "" {
		accountKey = os.Getenv("AZURESTORAGE_ACCESS_KEY")
	}

	if containerName == "" {
		containerName = "cache"
	}

	client, err := storage.NewBasicClient(accountName, accountKey)
	if err != nil {
		return nil, err
	}

	cache := Cache{
		client:    client.GetBlobService(),
		container: containerName,
	}

	_, err := cache.client.CreateContainerIfNotExists(cache.container, storage.ContainerAccessTypeBlob)
	if err != nil {
		return nil, err
	}

	return &cache, nil
}
