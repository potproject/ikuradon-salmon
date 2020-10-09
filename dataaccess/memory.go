package dataaccess

import (
	"fmt"
)

type dataAccessMemory struct {
	DataAccess
	db map[string]DataSet
}

// SetMemory Setting Using Memory
func SetMemory() {
	DA = dataAccessMemory{
		db: map[string]DataSet{},
	}
}

func (da dataAccessMemory) Get(key string) (DataSet, error) {
	var ds DataSet
	ds, ok := da.db[key]
	if !ok {
		return ds, fmt.Errorf("%s Not Found", key)
	}
	return ds, nil
}

func (da dataAccessMemory) Has(key string) (bool, error) {
	_, ok := da.db[key]
	return ok, nil
}

func (da dataAccessMemory) Set(key string, value DataSet) error {
	da.db[key] = value
	return nil
}

func (da dataAccessMemory) Delete(key string) error {
	delete(da.db, key)
	return nil
}

func (da dataAccessMemory) ListAll() ([]param, error) {
	p := []param{}
	for k, v := range da.db {
		p = append(p, param{
			Key:   k,
			Value: v,
		})
	}
	return p, nil
}

func (da dataAccessMemory) Close() error {
	return nil
}
