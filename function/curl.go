package function

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

func PostForm(url string, data []byte) (body []byte, err error) {
	client := &http.Client{Timeout: time.Second * 15}
	rsp, err := client.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(data))
	if err != nil {
		return
	}
	if rsp != nil && rsp.Body != nil {
		defer rsp.Body.Close()
	}

	if rsp.StatusCode != http.StatusOK {
		err = errors.New(rsp.Status)
		return
	}
	body, err = ioutil.ReadAll(rsp.Body)
	return
}
