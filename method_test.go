package govkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestMethod(t *testing.T) {
	methodName := "users.get"
	accessToken := "123" // вставить рабочий access token
	params := map[string]string{
		"user_ids": "1",
		"v":        "5.126",
	}

	r, err := Method(methodName, accessToken, params)

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

	u, err := createURL(methodName)

	e := "https://api.vk.com/method/" + methodName

	if err != nil {
		t.Error("returned error:", err.Error())
	}

	if u == nil {
		t.Error("url is nil")
	} else {
		if u.String() != e {
			t.Error("bad url\ngot:", u.String(), "\nexpected:", e)
		}
	}
}

func TestAddQueryParamsToURL(t *testing.T) {
	vkApiURL := "https://api.vk.com/method/"
	methodName := "users.get"
	accessToken := "123" // вставить рабочий access token
	params := map[string]string{
		"user_ids": "1",
		"v":        "5.126",
	}

	u, err := url.Parse(vkApiURL + methodName)
	if err != nil {
		t.Error("returned error:", err.Error())
	}

	e := fmt.Sprintf("%s%s?access_token=%s&lang=0&user_ids=1&v=5.126", vkApiURL, methodName, accessToken)

	u = addQueryParamsToURL(u, accessToken, params)
	if u == nil {
		t.Error("url is nil")
	} else {
		if u.String() != e {
			t.Error("bad url\ngot:", u.String(), "\nexpected:", e)
		}
	}
}

func TestSendRequest(t *testing.T) {
	vkApiURL := "https://api.vk.com/method/"
	methodName := "users.get"
	accessToken := "123" // вставить рабочий access token
	params := map[string]string{
		"user_ids": "1",
		"v":        "5.126",
	}

	u, err := url.Parse(vkApiURL + methodName)
	if err != nil {
		t.Error("returned error:", err.Error())
	}

	p := url.Values{}
	for key, param := range params {
		p.Set(key, param)
	}
	p.Set("access_token", accessToken)
	p.Set("lang", "0")

	u.RawQuery = p.Encode()

	r, err := sendRequest(u.String())
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
	vkApiURL := "https://api.vk.com/method/"
	methodName := "users.get"
	accessToken := "123" // вставить рабочий access token
	params := map[string]string{
		"user_ids": "1",
		"v":        "5.126",
	}

	u, err := url.Parse(vkApiURL + methodName)
	if err != nil {
		t.Error("returned error:", err.Error())
	}

	p := url.Values{}
	for key, param := range params {
		p.Set(key, param)
	}
	p.Set("access_token", accessToken)
	p.Set("lang", "0")

	u.RawQuery = p.Encode()

	response, err := http.Get(u.String())
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
