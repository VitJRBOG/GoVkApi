package govkapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestMethod(t *testing.T) {
	methodName := "users.get"
	params := url.Values{
		"access_token": {"123"}, // вставить рабочий access token
		"user_ids":     {"1"},
		"v":            {"5.126"},
	}

	r, err := Method(methodName, params)

	if err != nil {
		t.Error("returned error:", err.Error())
	}

	var values map[string]interface{}
	err = json.Unmarshal(r, &values)

	if err, exist := values["error"]; exist {
		t.Error("returned error:", err)
	}
}

func TestCreateURL(t *testing.T) {
	methodName := "users.get"

	u := createURL(methodName)

	e := "https://api.vk.com/method/" + methodName

	if len(u) == 0 {
		t.Error("url is empty")
	} else {
		if u != e {
			t.Error("bad url\ngot:", u, "\nexpected:", e)
		}
	}
}

func TestSendRequest(t *testing.T) {
	u := "https://api.vk.com/method/users.get"
	params := url.Values{
		"access_token": {"123"}, // вставить рабочий access token
		"user_ids":     {"1"},
		"v":            {"5.126"},
	}

	r, err := sendRequest(u, params)
	if err != nil {
		t.Error("returned error:", err.Error())
	}

	var values map[string]interface{}
	err = json.Unmarshal(r, &values)

	if err, exist := values["error"]; exist {
		t.Error("returned error:", err)
	}
}

func TestParseVkApiResponseBody(t *testing.T) {
	u := "https://api.vk.com/method/users.get"
	params := url.Values{
		"access_token": {"123"}, // вставить рабочий access token
		"user_ids":     {"1"},
		"v":            {"5.126"},
	}

	response, err := http.PostForm(u, params)
	if err != nil {
		t.Error("returned error:", err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Error("returned error:", err.Error())
	}

	v, err := parseVkApiResponseBody(body)
	if err != nil {
		t.Error("returned error:", err.Error())
	}

	if v == nil {
		t.Error("the result of parsing is nil")
	}
}
