package utils

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
)

func Get(urls string) ([]byte, error) {

	resp, err := getClient().Get(urls)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	tempData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err

	}
	return tempData, nil
}

func getClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	return client
}
