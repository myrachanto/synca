package load

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	httperors "github.com/myrachanto/erroring"
)

func Fetcher(url string) ([]*Product, httperors.HttpErr) {
	if url == "" {
		return nil, httperors.NewBadRequestError("url is empty")
	}
	response, err := http.Get(url)
	if err != nil {
		return nil, httperors.NewBadRequestError("something wentwrong fetching datat form that url")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, httperors.NewBadRequestError("something wentwrong reading data")
	}
	products := []*Product{}
	err = json.Unmarshal([]byte(body), &products)
	if err != nil {
		return nil, httperors.NewBadRequestError("something wentwrong reading data")
	}
	return products, nil
}
