package api

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetErrorFromHtml - retrieve error text from CTFd HTML
func GetErrorFromHtml(res http.Response) (*string, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	s := doc.Find("div.alert span").First().Text()
	return &s, nil
}
