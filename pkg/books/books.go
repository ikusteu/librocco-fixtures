package books

import (
	"encoding/json"
	"errors"
)

type VolumeInfo struct {
	Title         string   `json:"title"`
	Authors       []string `json:"authors"`
	Publisher     string   `json:"publisher"`
	PublishedDate string   `json:"publishedDate"`

	IndustryIdentifiers []struct {
		Type       string `json:"type"`
		Identifier string `json:"identifier"`
	} `json:"industryIdentifiers"`

	Categories []string `json:"categories"`
	Language   string   `json:"language"`
}

var ErrNoISBN = errors.New("NO_ISBN")

// #region BookItem
type BookItem struct {
	VolumeInfo VolumeInfo `json:"volumeInfo"`
}

// ISBN10 extracts an ISBN_10 from volume info and returns error if none found
func (b *BookItem) ISBN10() (isbn string, err error) {
	if ids := b.VolumeInfo.IndustryIdentifiers; len(ids) > 0 {
		for _, id := range ids {
			if id.Type == "ISBN_10" {
				return id.Identifier, nil
			}
		}
	}
	return "", ErrNoISBN
}

// ISBN13 extracts an ISBN_13 from volume info and returns error if none found
func (b *BookItem) ISBN13() (isbn string, err error) {
	if ids := b.VolumeInfo.IndustryIdentifiers; len(ids) > 0 {
		for _, id := range ids {
			if id.Type == "ISBN_13" {
				return id.Identifier, nil
			}
		}
	}
	return "", ErrNoISBN
}

// HasCategory checks if a book entry has a provided category
func (b *BookItem) HasCategory(cat string) bool {
	if cats := b.VolumeInfo.Categories; len(cats) > 0 {
		for _, c := range cats {
			if c == cat {
				return true
			}
		}
	}
	return false
}

// #endregion BookItem

// #region BookStock
type BookStock struct {
	*BookItem
	Warehouse string `json:"warehouse"`
	Quantity  int    `json:"quantity"`
}

func FromJSON(jsonStr []byte, warehouse string) (b *BookStock, err error) {
	b = &BookStock{Warehouse: warehouse}
	if err := json.Unmarshal(jsonStr, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *BookStock) Add(n int) {
	b.Quantity += n
}

func (b *BookStock) ToJSON() (jsonStr []byte, err error) {
	jsonStr, err = json.MarshalIndent(b, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonStr, nil
}

// #endregion BookStock
