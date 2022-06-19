package dataaccess

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBolt(t *testing.T) {
	dbName := "test.db"
	SetBolt(dbName)
	d := DataSet{
		SubscribeID:       "12",
		ExponentPushToken: "qwerty",
		PushPrivateKey:    "zxcvbn",
		PushAuth:          "0987654321",
		ServerKey:         "qazwsc",
		CreatedAt:         123,
		LastUpdatedAt:     456,
	}
	err := DA.Set("test", d)
	if err != nil {
		t.Error(err)
	}
	has, err := DA.Has("test")
	if err != nil {
		t.Error(err)
	}
	has, _ = DA.Has("testNotfound")
	if has {
		t.Error("Mismatch Notfound")
	}
	dag, err := DA.Get("test")
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(d, dag); diff != "" {
		t.Error(diff)
	}
	sp, err := DA.ListAll()
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff([]param{{Key: "test", Value: d}}, sp); diff != "" {
		t.Error(diff)
	}
	err = DA.Delete("test")
	if err != nil {
		t.Error(err)
	}
	_, err = DA.Get("test")
	if err == nil {
		t.Error("Mismatch Notfound")
	}
	err = DA.Close()
	if err != nil {
		t.Error(err)
	}
	os.Remove(dbName)
}
