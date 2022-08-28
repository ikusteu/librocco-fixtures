package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/ikusteu/librocco-fixtures/pkg/books"
	"github.com/ikusteu/librocco-fixtures/pkg/request"
	"github.com/ikusteu/librocco-fixtures/pkg/response"
	"github.com/ikusteu/librocco-fixtures/pkg/utils"
)

func newBookFilter(category string, cap int) response.Filter {
	return func(r *response.VolumeRes, b *books.BookItem) bool {
		if c := r.Count(); c >= cap {
			fmt.Printf("Capacity (100) reached, total books in res: %d\n", c)
			return false
		}
		isbn, err := b.ISBN10()
		if err != nil {
			return false
		}
		if r.ISBNSet.Exists(isbn) {
			fmt.Printf("ISBN '%s' already exists, skipping\n", isbn)
			return false
		}
		return true
	}
}

func getBooks(n int, category string, r *response.VolumeRes, startIndex int) {
	var res *response.VolumeRes
	if r == nil {
		res = &response.VolumeRes{}
	} else {
		res = r
	}

	nreq := (n / 40) + 1

	wg := sync.WaitGroup{}
	wg.Add(nreq)

	for i := 0; i < n; i += 40 {
		go (func(i int) {
			vReq := request.New()

			vReq.Fields(request.AllFields())
			vReq.Search(request.SearchConfig{
				Additional: map[string]string{
					"subject": category,
				},
			})

			vReq.StartIndex(startIndex + i)

			if diff := n - i; diff >= 40 {
				vReq.NumItems(40)
			} else {
				vReq.NumItems(diff)
			}

			vRes, err := vReq.Send(startIndex + i)
			if err != nil {
				log.Fatal(err.Error())
			}

			res.AddItems(vRes.Items, newBookFilter(category, 100))

			wg.Done()
		})(i)
	}

	wg.Wait()
}

func get100Books(category string) *response.VolumeRes {
	res := &response.VolumeRes{ISBNSet: utils.NewStringSet()}

	for n, si := 100, 0; n > 0; n, si = 100-res.Count(), si+40 {
		getBooks(100, category, res, si)
	}

	return res
}

func getBooksForCategory(category string, wg *sync.WaitGroup) {
	defer wg.Done()

	res := get100Books(category)

	fmt.Printf("Retrieved %d books for '%s' category\n", res.Count(), category)

	err := res.Store("./raw_books/" + category)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func importBooks() {
	categories := []string{
		"Jazz",
		"Science",
		"Horses",
		"Astronomy",
		"Education",
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(categories))

	for _, cat := range categories {
		go getBooksForCategory(cat, wg)
	}

	wg.Wait()
}
