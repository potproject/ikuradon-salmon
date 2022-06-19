package dataaccess

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
)

type dataAccessBolt struct {
	DataAccess
	db *bolt.DB
}

const BoltBucketName = "dataset"

// SetBolt Setting BoltDB
func SetBolt(path string) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BoltBucketName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	DA = dataAccessBolt{
		db: db,
	}
}

func (da dataAccessBolt) Get(key string) (DataSet, error) {
	var ds DataSet
	var data []byte
	err := da.db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(BoltBucketName)).Get([]byte(key))
		return nil
	})
	err = json.Unmarshal(data, &ds)
	if err != nil {
		return ds, err
	}
	return ds, nil

}

func (da dataAccessBolt) Has(key string) (bool, error) {
	_, err := da.Get(key)
	if err != nil {
		return false, err
	}
	return true, err
}

func (da dataAccessBolt) Set(key string, value DataSet) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return da.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(BoltBucketName)).Put([]byte(key), j)
	})
}

func (da dataAccessBolt) Delete(key string) error {
	return da.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(BoltBucketName)).Delete([]byte(key))
	})
}

func (da dataAccessBolt) ListAll() ([]param, error) {
	p := []param{}
	err := da.db.View(func(tx *bolt.Tx) error {
		tx.Bucket([]byte(BoltBucketName)).ForEach(func(k, v []byte) error {
			var ds DataSet
			json.Unmarshal(v, &ds)
			p = append(p, param{Key: string(k), Value: ds})
			return nil
		})
		return nil
	})
	return p, err
}

func (da dataAccessBolt) Close() error {
	return da.db.Close()
}
