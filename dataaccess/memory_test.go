package dataaccess

import (
	"reflect"
	"testing"
)

func TestSetMemory(t *testing.T) {
	SetMemory()
	if reflect.TypeOf(DA) != reflect.TypeOf(dataAccessMemory{db: map[string]DataSet{}}) {
		t.Errorf("Invaild Type: %v", reflect.TypeOf(DA))
	}
	DA = nil
}

func TestSetHasGetDeleteListAllClose(t *testing.T) {
	SetMemory()
	if reflect.TypeOf(DA) != reflect.TypeOf(dataAccessMemory{db: map[string]DataSet{}}) {
		t.Errorf("Invaild Type: %v", reflect.TypeOf(DA))
	}
	act := DataSet{
		SubscribeID:       "ABCDEFG123456",
		ExponentPushToken: "Expo[xxxxxx]",
		PushPrivateKey:    "PushPrivateKey",
		PushAuth:          "PushAuth",
		ServerKey:         "ServerKey",
		CreatedAt:         1600000000,
		LastUpdatedAt:     1600000000,
	}
	// Set
	DA.Set("key", act)

	// Has
	b, _ := DA.Has("key")
	if !b {
		t.Errorf("Has: Key Not found")
	}

	// Get exist
	exp, _ := DA.Get("key")
	if !reflect.DeepEqual(act, exp) {
		t.Errorf("Get: Invaild Param: %v %v", act, exp)
	}

	// Get notExist
	_, err := DA.Get("notExist")
	if err == nil {
		t.Errorf("Get: notExist NoErr")
	}

	// ListAll
	p, _ := DA.ListAll()
	listExp := []param{
		{
			Key:   "key",
			Value: act,
		},
	}
	if !reflect.DeepEqual(p, listExp) {
		t.Errorf("Get: Invaild Param: %v %v", p, listExp)
	}

	// Delete
	DA.Delete("key")

	// Close
	if DA.Close() != nil {
		t.Errorf("Close: Err")
	}

	// Destroy
	DA = nil
}
