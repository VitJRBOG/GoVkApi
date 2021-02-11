package govkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Method(methodName, accessToken string, values map[string]string) ([]byte, error) {
	u, err := createURL(methodName)
	if err != nil {
		return nil, err
	}

	u = addQueryParamsToURL(u, accessToken, values)

	r, err := sendRequest(u.String())
	if err != nil {
		return nil, err
	}

	v, err := parseVkApiResponseBody(r)
	if err != nil {
		return nil, err
	}

	if v.Error != nil {
		return nil, v.Error
	}

	return v.Response, nil
}

func createURL(methodName string) (*url.URL, error) {
	vkApiURL := "https://api.vk.com/method/"
	u, err := url.Parse(vkApiURL + methodName)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func addQueryParamsToURL(u *url.URL, accessToken string, params map[string]string) *url.URL {
	p := url.Values{}
	for key, param := range params {
		p.Set(key, param)
	}
	p.Set("access_token", accessToken)
	p.Set("lang", "0")

	u.RawQuery = p.Encode()

	return u
}

func sendRequest(u string) ([]byte, error) {
	response, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

type vkApiResponseBody struct {
	Error    *vkApiError
	Response json.RawMessage
}

type vkApiError struct {
	Code          int    `json:"error_code"`
	Message       string `json:"error_msg"`
	RequestParams []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"request_params"`
}

func (e *vkApiError) Error() string {
	return fmt.Sprintf("code %d: %s", e.Code, e.Message)
}

func parseVkApiResponseBody(b []byte) (*vkApiResponseBody, error) {
	var v vkApiResponseBody
	err := json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}
