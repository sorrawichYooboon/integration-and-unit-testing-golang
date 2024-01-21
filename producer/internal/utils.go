package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func MarshalJSONData(req interface{}) []byte {
	jsonData, err := json.Marshal(req)
	if err != nil {
		panic(errors.New("cannot marshal json data"))
	}

	return jsonData
}

func UnmarshalJSONData(req []byte, target interface{}) {
	err := json.Unmarshal(req, target)
	if err != nil {
		panic(errors.New("cannot unmarshal json data"))
	}
}

func CreateHTTPRequest(method, url string, body []byte) *http.Request {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Content-Type", "application/json")
	return request
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
