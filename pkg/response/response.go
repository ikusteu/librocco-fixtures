package response

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ikusteu/librocco-fixtures/pkg/books"
	"github.com/ikusteu/librocco-fixtures/pkg/utils"
)

type VolumeRes struct {
	Items   []books.BookItem `json:"items"`
	ISBNSet *utils.StringSet
}

// FromJSON creates a VolumeRes struct from raw json, received from the request
func FromJSON(jsonr []byte) (res *VolumeRes, err error) {
	res = &VolumeRes{}

	if err := json.Unmarshal(jsonr, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Categories returns a set (non repeating slice) of available categories
// in the response
func (vr *VolumeRes) Categories() *utils.StringSet {
	categories := utils.NewStringSet()
	for _, b := range vr.Items {
		if cats := b.VolumeInfo.Categories; len(cats) > 0 {
			categories.AddSlice(cats)
		}
	}
	return categories
}

func (vr *VolumeRes) FilterByCategory(cat string) *VolumeRes {
	nvr := &VolumeRes{}
	for _, el := range vr.Items {
		if el.HasCategory(cat) {
			nvr.Items = append(nvr.Items, el)
		}
	}
	return nvr
}

type Filter func(r *VolumeRes, b *books.BookItem) bool

func (vr *VolumeRes) AddItems(items []books.BookItem, f Filter) {
	for _, it := range items {
		if f == nil || f(vr, &it) {
			vr.Items = append(vr.Items, it)
			if isbn, err := it.ISBN10(); err == nil {
				vr.ISBNSet.Add(isbn)
			}
		}
	}
}

func (vr *VolumeRes) Count() int {
	return len(vr.Items)
}

type CountMap map[string]int

func (vr *VolumeRes) CountCategries() CountMap {
	count := make(map[string]int)
	for _, el := range vr.Items {
		if cats := el.VolumeInfo.Categories; len(cats) > 0 {
			for _, c := range cats {
				count[c]++
			}
		}
	}
	return count
}

func (m CountMap) Print() {
	for v, c := range m {
		fmt.Printf("[%s]: %d\n", v, c)
	}
}

func (vr *VolumeRes) Store(dirpath string) error {
	fmt.Print("Initializing save to fs...\n\n")
	fmt.Printf("Looking for dir: %s\n", dirpath)
	if _, err := os.ReadDir(dirpath); err == nil {
		fmt.Printf("Dir '%s' found!\n\n", dirpath)
	} else {
		fmt.Printf("No dir '%s' found, creating new...\n", dirpath)
		if err := os.Mkdir(dirpath, 0755); err != nil {
			return err
		} else {
			fmt.Printf("Dir '%s' successfully created!\n\n", dirpath)
		}
	}

	fmt.Printf("Storing files to '%s'...\n\n", dirpath)
	for _, book := range vr.Items {
		isbn, _ := book.ISBN10()
		fn := fmt.Sprintf("%v/%v.json", dirpath, isbn)
		fmt.Printf("Writing to file: %v\n", fn)

		jsons, err := json.MarshalIndent(book, "", "  ")
		if err != nil {
			return err
		}

		os.WriteFile(fn, jsons, 0755)
		fmt.Printf("Success!\n")
	}

	return nil
}
