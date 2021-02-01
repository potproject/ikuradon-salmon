package dataaccess

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	redismock "github.com/go-redis/redismock/v8"
)

func TestSetHasGetDeleteListAllCloseRedis(t *testing.T) {
	client, mock := redismock.NewClientMock()
	DA = dataAccessRedis{
		client: client,
	}
	if reflect.TypeOf(DA) != reflect.TypeOf(dataAccessRedis{client: client}) {
		t.Errorf("Invaild Type: %v", reflect.TypeOf(DA))
	}
	act := DataSet{
		SubscribeID:       "ABCDEFG123456",
		UserID:            "100200300",
		Username:          "UserName",
		Domain:            "server.mastodon.net",
		AccessToken:       "AccessToken",
		ExponentPushToken: "Expo[xxxxxx]",
		PushPrivateKey:    "PushPrivateKey",
		PushPublicKey:     "PushPublicKey",
		PushAuth:          "PushAuth",
		ServerKey:         "ServerKey",
		CreatedAt:         1600000000,
		LastUpdatedAt:     1600000000,
	}
	// Set
	DA.Set("key", act)

	// Has
	mock.ExpectExists("key").SetVal(1)
	b, _ := DA.Has("key")
	if !b {
		t.Errorf("Has: Key Not found")
	}

	// Get exist
	jact, _ := json.Marshal(act)
	mock.ExpectGet("key").SetVal(string(jact))
	exp, _ := DA.Get("key")
	if !reflect.DeepEqual(act, exp) {
		t.Errorf("Get: Invaild Param: %v %v", act, exp)
	}

	// Get notExist
	mock.ExpectGet("notExist").SetErr(errors.New("Not Exist"))
	_, err := DA.Get("notExist")
	if err == nil {
		t.Errorf("Get: notExist NoErr")
	}

	// ListAll
	listAllkey := []string{"keyl"}
	mock.ExpectKeys("*").SetVal(listAllkey)
	mock.ExpectGet("keyl").SetVal(string(jact))
	p, err := DA.ListAll()
	if err != nil {
		t.Error("ListAll:", err)
	}
	listExp := []param{
		{
			Key:   "keyl",
			Value: act,
		},
	}
	if !reflect.DeepEqual(p, listExp) {
		t.Errorf("Get: Invaild Param: %v %v", p, listExp)
	}

	// Delete
	mock.ExpectDel("key").SetVal(1)
	DA.Delete("key")

	// Close
	if DA.Close() != nil {
		t.Errorf("Close: Err")
	}

	// Destroy
	DA = nil
}
