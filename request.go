package viber

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// PostData to viber API
func (v *Viber) PostData(url string, i interface{}) ([]byte, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	Log.Println("Post data:", string(b))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Close = true
	req.Header.Add("X-Viber-Auth-Token", v.AppKey)

	c := &http.Client{
		Timeout: time.Duration(v.RequestTimeout) * time.Second,
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
