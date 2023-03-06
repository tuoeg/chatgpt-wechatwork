package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPOptions struct {
	ContentType string
	Token       string
	Username    string
	Password    string
}

// 发送http请求接口
func SendHTTPRequest(client *http.Client, url, method string, reqBody interface{}, option *HTTPOptions) ([]byte, int, error) {
	bodyBuffer, ok := reqBody.(*bytes.Buffer)
	if !ok {
		reqData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		bodyBuffer = bytes.NewBuffer(reqData)
	}

	req, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if option != nil {
		if option.ContentType != "" {
			req.Header.Set("Content-Type", option.ContentType)
		}
		if option.Token != "" {
			fmt.Println("auth", option.Token)
			req.Header.Add("Authorization", option.Token)
		}
		if option.Username != "" && option.Password != "" {
			req.SetBasicAuth(option.Username, option.Password)
		}
	}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(rsp.Body)
		fmt.Println(string(body))
		return nil, http.StatusInternalServerError, fmt.Errorf("unnormal http status code: %d\n", rsp.StatusCode)
	}
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return body, rsp.StatusCode, nil
}
