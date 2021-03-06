package govkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Method(methodName string, values url.Values) ([]byte, error) {
	u := createURL(methodName)

	r, err := sendRequest(u, values)
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

func createURL(methodName string) string {
	u := fmt.Sprintf("https://api.vk.com/method/%s", methodName)
	return u
}

func sendRequest(u string, values url.Values) ([]byte, error) {
	response, err := http.PostForm(u, values)
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
