package google

import "net/url"

type Google struct {
}

func NewGoogle() (*Google) {
	google := new(Google)
	return google
}

func (google *Google) GetSearchURL(word string) (string) {
	u := new(url.URL)
	u.Scheme = "https"
	u.Host = "www.google.com"
	u.Path = "search"
	q := u.Query()
	q.Add("q", word)
	u.RawQuery = q.Encode()
	return u.String()
}
