package dataaccess

import (
	"encoding/json"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type dataAccessLevel struct {
	DataAccess
	db *leveldb.DB
}

// SetLevel Setting LevelDB
func SetLevel(path string) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Fatal(err)
	}
	DA = dataAccessLevel{
		db: db,
	}
}

func (da dataAccessLevel) Get(key string) (DataSet, error) {
	var ds DataSet
	data, err := da.db.Get([]byte(key), nil)
	if err != nil {
		return ds, err
	}
	err = json.Unmarshal(data, &ds)
	if err != nil {
		return ds, err
	}
	return ds, nil
}

func (da dataAccessLevel) Has(key string) (bool, error) {
	return da.db.Has([]byte(key), nil)
}

func (da dataAccessLevel) Set(key string, value DataSet) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return da.db.Put([]byte(key), j, nil)
}

func (da dataAccessLevel) Delete(key string) error {
	return da.db.Delete([]byte(key), nil)
}

func (da dataAccessLevel) ListAll() ([]param, error) {
	p := []param{}
	iter := da.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		var ds DataSet
		json.Unmarshal(value, &ds)
		p = append(p, param{
			Key:   string(key),
			Value: ds,
		})
	}
	iter.Release()
	err := iter.Error()
	return p, err
}

func (da dataAccessLevel) Close() error {
	return da.db.Close()
}
