package utils

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
)

func Get(urls string) ([]byte, error) {

	req, err := http.NewRequest("GET", urls, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Googlebot/2.1 (+http://www.google.com/bot.html)")

	resp, err := getClient().Do(req)

	tempData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err

	}
	resp.Body.Close()
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
