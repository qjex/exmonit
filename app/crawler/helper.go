package crawler

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func getJson(httpClient *http.Client, req *http.Request) (map[string]interface{}, error) {
	rs, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error performing http request")
	}
	if rs.StatusCode != 200 {
		return nil, errors.Wrapf(err, "status code=%d", rs.StatusCode)
	}
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body bytes")
	}
	_ = rs.Body.Close()
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrapf(err, "error deserializing response %s", body)
	}
	return data, nil
}
