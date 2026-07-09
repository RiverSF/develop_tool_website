package net

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

func HttpPostRequest(apiUrl string, bytesData []byte, headers map[string]string, httpClient *http.Client) (*http.Response, []byte, error) {
	var err error

	isGzip := false
	if len(headers) > 0 {
		for key, item := range headers {
			if key == "Accept-Encoding" && item == "gzip" {
				isGzip = true
				break
			}
		}
	}

	var reader *bytes.Buffer

	if isGzip {
		var zBuf bytes.Buffer
		zw := gzip.NewWriter(&zBuf)
		if _, err = zw.Write(bytesData); err != nil {
			return nil, []byte{}, fmt.Errorf("gzip error='%w'", err)
		}
		zw.Close()
		reader = &zBuf
	} else {
		reader = bytes.NewBuffer(bytesData)
	}

	request, err := http.NewRequest("POST", apiUrl, reader)
	if err != nil {
		return nil, []byte{}, fmt.Errorf("newRequest error='%w'", err)
	}

	for key, item := range headers {
		request.Header.Set(key, item)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, []byte{}, fmt.Errorf("clientDo error='%w'", err)
	}

	defer response.Body.Close()

	body := response.Body
	if response.Header.Get("Content-Encoding") == "gzip" {
		body2, e := gzip.NewReader(response.Body)
		if e == nil {
			body = body2
		} else if e != io.EOF {
			return nil, []byte{}, fmt.Errorf("unzip error='%w'", e)
		}
	}

	data, err := io.ReadAll(body)
	_, _ = io.Copy(io.Discard, response.Body)
	if err != nil {
		return nil, []byte{}, fmt.Errorf("read body error='%w'", err)
	}

	return response, data, nil
}

func HttpGetRequest(apiUrl string, httpClient *http.Client) (*http.Response, []byte, error) {
	response, err := httpClient.Get(apiUrl)
	if err != nil {
		return nil, []byte{}, fmt.Errorf("clientDo error='%w'", err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	_, _ = io.Copy(io.Discard, response.Body)
	if err != nil {
		return nil, []byte{}, fmt.Errorf("read body error='%w'", err)
	}

	return response, body, nil
}
