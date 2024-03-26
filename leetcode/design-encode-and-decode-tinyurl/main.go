package main

import (
	"crypto/md5"
	"encoding/hex"
)

type Codec struct {
	LongURLHashMap  map[string]string // long -> short
	ShortURLHashMap map[string]string // short -> long
}

func Constructor() Codec {
	return Codec{
		LongURLHashMap:  make(map[string]string),
		ShortURLHashMap: make(map[string]string),
	}
}

// Encodes a URL to a shortened URL.
func (this *Codec) encode(longUrl string) string {
	v, ok := this.LongURLHashMap[longUrl]
	if ok {
		return v
	}
	hash := GetMD5Hash(longUrl)
	url := "http://tinyurl.com" + string(hash)
	this.ShortURLHashMap[url] = longUrl
	return url

}

// Decodes a shortened URL to its original URL.
func (this *Codec) decode(shortUrl string) string {
	return this.ShortURLHashMap[shortUrl]
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
