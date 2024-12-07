package main

import (
	"DestributedKV/kv/storage"
)

func main() {
	s, _ := storage.BuildDB("D:\\TestPath")
	s.Put([]byte("111"), []byte("222"))

}
