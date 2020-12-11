package govkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func SendRequestVkApi(accessToken, methodName string, paramsMap map[string]string) ([]byte, error) {
	var qh VkApiQueryHandler
	qh.makeURL(accessToken, methodName, paramsMap)
	qh.sendQuery()
	qh.parseRespBody()

	if qh.Error != nil {
		return nil, qh.Error
	}

	return qh.Response, nil
}

type VkApiQueryHandler struct {
	QueryURL *url.URL
	RespBody []byte
	Error    *Error
	Response json.RawMessage
}

func (qh *VkApiQueryHandler) makeURL(accessToken, methodName string, paramsMap map[string]string) {
	var q QueryURL
	q.createURL(methodName)
	q.setParams(accessToken, paramsMap)
	qh.QueryURL = q.URL
}

func (qh *VkApiQueryHandler) sendQuery() {
	resp, err := http.Get(qh.QueryURL.String())
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	qh.RespBody = body
}

func (qh *VkApiQueryHandler) parseRespBody() {
	var handler struct {
		Error    *Error
		Response json.RawMessage
	}
	var err error
	err = json.Unmarshal(qh.RespBody, &handler)
	if err != nil {
		panic(err.Error())
	}
	qh.Error = handler.Error
	qh.Response = handler.Response
}

type QueryURL struct {
	URL *url.URL
}

func (q *QueryURL) createURL(methodName string) {
	apiURL := "https://api.vk.com/method/"
	var err error
	q.URL, err = url.Parse(apiURL + methodName)
	if err != nil {
		panic(err.Error())
	}
}

func (q *QueryURL) setParams(accessToken string, paramsMap map[string]string) {
	params := url.Values{}
	for key, value := range paramsMap {
		params.Set(key, value)
	}
	params.Set("access_token", accessToken)
	params.Set("lang", "0")

	q.URL.RawQuery = params.Encode()
}

type Error struct {
	Code          int    `json:"error_code"`
	Message       string `json:"error_msg"`
	RequestParams []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"request_params"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code %d: %s", e.Code, e.Message)
}
