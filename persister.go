package main

import (
	"github.com/kennygrant/sanitize"
	"github.com/patrickmn/go-cache"
	"github.com/peacemakr-io/peacemakr-go-sdk/pkg/utils"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type filePersister struct {
	directoryToSave string
	cache           *cache.Cache
}

func (p *filePersister) getFilePath(key string) string {
	x := p.directoryToSave + "/" + sanitize.Path(strings.Replace(key, "/", "-", -1))
	return x
}

func (p *filePersister) getInCache(key string) *string {
	var foo string
	if x, found := p.cache.Get(p.getFilePath(key)); found {
		foo = x.(string)
		return &foo
	}
	return nil
}

func (p *filePersister) setInCache(key, value string) {
	p.cache.Set(p.getFilePath(key), value, cache.DefaultExpiration)
}

func (p *filePersister) Exists(key string) bool {
	exists := p.getInCache(key)
	if exists != nil {
		return true
	}

	if _, err := os.Stat(p.getFilePath(key)); err == nil {
		return true
	}

	return false
}

func (p *filePersister) Save(key, value string) error {
	p.setInCache(key, value)
	err := ioutil.WriteFile(p.getFilePath(key), []byte(value), 0700)
	return err
}

func (p *filePersister) Load(key string) (string, error) {
	exists := p.getInCache(key)
	if exists != nil {
		return *exists, nil
	}

	b, err := ioutil.ReadFile(p.getFilePath(key))
	if err != nil {
		return "", err
	}
	return string(b), err
}


func GetDiskPersister(path string) utils.Persister {

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)

	fileP := &filePersister{
		path,
		c,
	}
	return utils.Persister(fileP)
}
