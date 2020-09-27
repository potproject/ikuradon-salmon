package setting

import (
	"os"
	"strconv"
	"testing"
)

func TestSetSetting(t *testing.T) {
	appName := "ikuradon-salmon-test"
	appVersion := "99.99.99"
	salt := "SALT"
	baseURL := "https://ikuradon-salmon.example.com/"
	apiHost := "0.0.0.0"
	apiPort := 8080
	useRedis := false
	os.Setenv("APP_NAME", appName)
	os.Setenv("APP_VERSION", appVersion)
	os.Setenv("SALT", salt)
	os.Setenv("BASE_URL", baseURL)
	os.Setenv("API_HOST", apiHost)
	os.Setenv("API_PORT", strconv.Itoa(apiPort))
	os.Setenv("USE_REDIS", "false")
	defer os.Unsetenv("APP_NAME")
	defer os.Unsetenv("APP_VERSION")
	defer os.Unsetenv("SALT")
	defer os.Unsetenv("BASE_URL")
	defer os.Unsetenv("API_HOST")
	defer os.Unsetenv("API_PORT")
	defer os.Unsetenv("USE_REDIS")
	SetSetting()
	if S.AppName != appName {
		t.Error(" S.AppName / Actual:" + S.AppName + " Expect:" + appName)
	}
	if S.AppVersion != appVersion {
		t.Error(" S.AppVersion / Actual:" + S.AppVersion + " Expect:" + appVersion)
	}
	if S.Salt != salt {
		t.Error(" S.Salt / Actual:" + S.Salt + " Expect:" + salt)
	}
	if S.BaseURL != baseURL {
		t.Error(" S.BaseURL / Actual:" + S.BaseURL + " Expect:" + baseURL)
	}
	if S.APIHost != apiHost {
		t.Error(" S.APIHost / Actual:" + S.APIHost + " Expect:" + apiHost)
	}
	if S.APIPort != uint16(apiPort) {
		t.Errorf(" S.APIPort / Actual: %d Expect: %d", S.APIPort, apiPort)
	}
	if S.UseRedis != useRedis {
		t.Errorf(" S.APIPort / Actual: %t Expect: %t", S.UseRedis, useRedis)
	}
}
