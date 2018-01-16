package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	HttpClient *http.Client
)

func init() {
	HttpClient = http.DefaultClient
}

func httpDo(req *http.Request, v interface{}) error {
	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}
