package store

import (
	"AdachiAndShimamura/DistributedKV/kv/engine_util"
)

type RegionReader struct {
	storage *engine_util.Engines
}

func (r *RegionReader) GetCF(cf string, key []byte) ([]byte, error) {
	val, err := engine_util.GetCF(r.storage.Kv, cf, key)
	return val, err
}

func (r *RegionReader) IterCF(cf string) engine_util.DBIterator {
	//TODO implement me
	panic("implement me")
}

func (r *RegionReader) Close() {
	//TODO implement me
	panic("implement me")
}
