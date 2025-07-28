package xsoar

import (
	"encoding/json"
	"net/http"
)

func Decode[T any](resp *http.Response) (T, error) {
	defer resp.Body.Close()
	v := new(T)
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	return *v, decoder.Decode(v)
}

func GetMessage(resp *http.Response) string {
	message, err := Decode[json.RawMessage](resp)
	if err != nil {
		return ""
	}
	return string(message)
}
