package storage

import "github.com/dgraph-io/badger/v4"

type Storage struct {
	db *badger.DB
}

func (s *Storage) Put(key []byte, value []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
}
func (s *Storage) Get(key []byte) (value []byte, _ error) {
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			value = val
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	return value, err
}

func (s *Storage) Del(key []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		return err
	})
}

func (s *Storage) Close() error {
	err := s.db.Close()
	return err
}

func BuildDB(path string) (storage *Storage, err error) {
	opt := badger.DefaultOptions(path)
	opt.Dir = path
	opt.ValueDir = path
	db, err := badger.Open(opt)
	return &Storage{db: db}, err
}
