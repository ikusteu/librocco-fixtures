package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ikusteu/librocco-fixtures/pkg/response"
)

const basepath = "https://www.googleapis.com/books/v1/volumes"

type QueryFields int

const (
	Title QueryFields = iota
	Authors
	Publisher
	PublishedDate
	IndustryIdentifiers
	Categories
	Language
)

func (q QueryFields) String() string {
	return [...]string{
		"volumeInfo/title",
		"volumeInfo/authors",
		"volumeInfo/publisher",
		"volumeInfo/publishedDate",
		"volumeInfo/industryIdentifiers",
		"volumeInfo/categories",
		"volumeInfo/language",
	}[q]
}

func AllFields() []QueryFields {
	return []QueryFields{
		Title,
		Authors,
		Publisher,
		PublishedDate,
		IndustryIdentifiers,
		Categories,
		Language,
	}
}

type SearchConfig struct {
	Term       string
	Additional map[string]string
}

type VolumeReq struct {
	url     string
	queries []string
}

func New() *VolumeReq {
	return &VolumeReq{url: basepath}
}

func (vr *VolumeReq) addQuery(q string) {
	vr.queries = append(vr.queries, q)

	qs := strings.Join(vr.queries, "&")

	vr.url = fmt.Sprintf("%s?%s", basepath, qs)
}

// Search method adds a search query to the request ("q=")
func (vr *VolumeReq) Search(s SearchConfig) {
	params := []string{}

	// Add main serch term
	if sstr := s.Term; sstr != "" {
		params = append(params, sstr)
	}

	// Add field specific search
	if add := s.Additional; len(add) > 0 {
		for fld, srch := range add {
			sstr := fmt.Sprintf("%s:%s", fld, srch)
			params = append(params, sstr)
		}
	}

	q := fmt.Sprintf("q=%s", strings.Join(params, "+"))

	vr.addQuery(q)
}

// Fields method adds "fields" query restricting the response received from the API
// fieldsConf is provided in form of [field]subfields map
func (vr *VolumeReq) Fields(fields []QueryFields) {

	strs := []string{}
	for _, fIota := range fields {
		f := fIota.String()
		strs = append(strs, f)
	}

	q := fmt.Sprintf("fields=items(%s)", strings.Join(strs, ","))

	vr.addQuery(q)
}

// NumItems method adds "maxItems" query param, essentially specifying number of items to receive
func (vr *VolumeReq) NumItems(n int) {
	q := fmt.Sprintf("maxResults=%d", n)
	vr.addQuery(q)
}

// StartIndex method adds "startIndex" query param for pagination
func (vr *VolumeReq) StartIndex(i int) {
	q := fmt.Sprintf("startIndex=%d", i)
	vr.addQuery(q)
}

// Send method executes the request using http.DefaultClient.Do and returns expected response or error
func (vr *VolumeReq) Send(i int) (r *response.VolumeRes, err error) {
	fmt.Printf("Sending request to (req %d):\n%s\n\n", i, vr.url)

	req, err := http.NewRequest(http.MethodGet, vr.url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if c := res.StatusCode; c >= 400 {
		// Pretty print the json response
		// if request failed
		jsonb := &bytes.Buffer{}
		json.Indent(jsonb, b, "", "  ")

		errStr := fmt.Sprintf("Request to '%s' failed with status code %d, response:\n%s\n", vr.url, c, jsonb.String())
		err := errors.New(errStr)

		return nil, err
	}

	return response.FromJSON(b)
}
