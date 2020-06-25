package dataAccess

var DA DataAccess

type DataAccess interface {
	Get(key string) (DataSet, error)
	Has(key string) (bool, error)
	Set(key string, value DataSet) error
	Delete(key string) error
	ListAll() ([]param, error)
	Close() error
}

type param struct {
	Key   string  // SubScribeId
	Value DataSet // DataSet
}

type DataSet struct {
	SubscribeId        string `json:"subscribe_id"`
	UserID             string `json:"user_id"`
	Username           string `json:"username"`
	Domain             string `json:"domain"`
	AccessToken        string `json:"access_token"`
	ExponentPushToken  string `json:"exponent_push_token"`
	PushPrivateKey     string `json:"push_private_key"`
	PushPublicKey      string `json:"push_public_key"`
	PushAuth           string `json:"push_auth"`
	ServerKey          string `json:"server_key"`
	CreatedAt          int64  `json:"created_at"`
	ExpiredAt          int64  `json:"expired_at"`
	LastUpdatedAt      int64  `json:"last_updated_at"`
	ServerLastId       int64  `json:"server_last_id"`
	NotificationsCount int64  `json:"notifications_count"`
}