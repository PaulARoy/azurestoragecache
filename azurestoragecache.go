// Package azurestoragecache provides an implementation of httpcache.Cache that
// stores and retrieves data using Azure Storage.
package azurestoragecache // import "github.com/PaulARoy/azurestoragecache"

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"bytes"

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
	AccountKey string

	// Container name to use to store blob
	ContainerName string
}

var noLogErrors, _ = strconv.ParseBool(os.Getenv("NO_LOG_AZUREBSCACHE_ERRORS"))

func (c *Cache) Get(key string) (resp []byte, ok bool) {
	rdr, err := c.Client.GetBlob(c.Config.ContainerName, key)
	if err != nil {
		return []byte{}, false
	}
	rdr.Close()
	
	resp, err = ioutil.ReadAll(rdr)
	if err != nil {
		if !noLogErrors {
			log.Printf("azurestoragecache.Get failed: %s", err)
		}
	}
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

func (c *Cache) Delete(key string) bool {
	res, err := c.Client.DeleteBlobIfExists(c.Config.ContainerName, key, nil)
	if err != nil {
		if !noLogErrors {
			log.Printf("azurestoragecache.Delete failed: %s", err)
		}
		return false
	}
	return res
}

// New returns a new Cache with underlying client for Azure Storage
//
// containerName is the container name for azure blob service
//
// The environment variables AZURESTORAGE_ACCOUNT_NAME and AZURESTORAGE_ACCESS_KEY 
// are used as credentials. To use different credentials, construct a Cache object 
// manually.
func New(accountName string, accountKey string, containerName string) *Cache {
	cache := Cache{
		Config: Config{
			AccountName: accountName, // || os.Getenv("AZURESTORAGE_ACCOUNT_NAME"),
			AccountKey: accountKey, // || os.Getenv("AZURESTORAGE_ACCESS_KEY"),
			ContainerName: containerName,
		},
	}

	api, err := vendorstorage.NewBasicClient(cache.Config.AccountName, cache.Config.AccountKey)
	if err != nil {
		return nil
	}
	
	cache.Client = api.GetBlobService()
	cache.Client.CreateContainerIfNotExists(cache.Config.ContainerName, 
											vendorstorage.ContainerAccessTypeBlob)
	return &cache
}