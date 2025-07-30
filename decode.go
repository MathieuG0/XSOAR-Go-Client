package xsoar

import (
	"encoding/json"
	"net/http"
)

func HTTPResponseDecode[T any](resp *http.Response) (T, error) {
	defer resp.Body.Close()
	v := new(T)
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	return *v, decoder.Decode(v)
}

func GetMessage(resp *http.Response) string {
	message, err := HTTPResponseDecode[json.RawMessage](resp)
	if err != nil {
		return ""
	}
	return string(message)
}
