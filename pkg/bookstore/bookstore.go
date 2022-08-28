package bookstore

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"sync"

	"github.com/ikusteu/librocco-fixtures/pkg/books"
	"github.com/ikusteu/librocco-fixtures/pkg/utils"
)

// #region Store
type Store struct {
	Books map[string]*books.BookStock
	Isbns *utils.StringSet
	mu    *sync.Mutex
}

func New() *Store {
	m := make(map[string]*books.BookStock)
	k := utils.NewStringSet()
	return &Store{m, k, &sync.Mutex{}}
}

func (bs Store) Add(b *books.BookStock) error {
	isbn, err := b.ISBN10()
	if err != nil {
		return err
	}

	bs.mu.Lock()

	bs.Books[isbn] = b
	bs.Isbns.Add(isbn)

	bs.mu.Unlock()

	return nil
}

func (bs *Store) GetRandom() *books.BookStock {
	isbn := bs.Isbns.GetRandom()
	return bs.Books[isbn]
}

// #region Store

// #region MultiStoreAggregator
type StoreAggregator interface {
	Add(b *books.BookStock) error
}

type MultiStoreAggregator struct {
	Stores []*Store
}

func (ag MultiStoreAggregator) Add(b *books.BookStock) error {
	if len(ag.Stores) == 0 {
		return errors.New("no stores to aggregate")
	}
	for _, s := range ag.Stores {
		s.Add(b)
	}
	return nil
}

// #endregion MultiStoreAggregator

// #region JSONBookLoader
type JSONBookLoader struct {
	Dirname   string
	Warehouse string
}

func (l *JSONBookLoader) IngestDir(store StoreAggregator) error {
	files, err := os.ReadDir(l.Dirname)
	if err != nil {
		return err
	}

	jsonFiles := filterJSONFiles(files)
	nf := len(jsonFiles)
	if nf == 0 {
		fmt.Printf("No JSON files found in '%s'\n", l.Dirname)
		return nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(nf)

	for _, fn := range jsonFiles {
		fPath := fmt.Sprintf("%s/%s", l.Dirname, fn)

		go func(fPath string) {
			defer func() {
				if rec := recover(); rec != nil {
					fmt.Println(rec)
				}
				wg.Done()
			}()

			jsonr, err := os.ReadFile(fPath)
			if err != nil {
				panic(err.Error())
			}

			b, err := books.FromJSON(jsonr, l.Warehouse)
			if err != nil {
				panic(err.Error())
			}

			if err := store.Add(b); err != nil {
				panic(err.Error())
			}
		}(fPath)
	}

	wg.Wait()

	return nil
}

func filterJSONFiles(contents []fs.DirEntry) []string {
	res := []string{}
	for _, el := range contents {
		fn := el.Name()
		regex := regexp.MustCompile(".json$")

		if !el.IsDir() && regex.Match([]byte(fn)) {
			res = append(res, fn)
		}
	}

	return res
}

// #endregion JSONBookLoader
