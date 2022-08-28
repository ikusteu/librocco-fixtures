package main

import (
	"fmt"
	"sync"

	"github.com/ikusteu/librocco-fixtures/pkg/bookstore"
)

var Warehouses = [...]string{
	"all",
	"jazz",
	"science",
	"horses",
	"astronomy",
	"education",
}

var WarehouseStoreLookup = make(map[string]*bookstore.Store)

// Initialise all stores in warehouse store lookup
// not to break the app if called later on
func init() {
	for _, wh := range Warehouses {
		{
			WarehouseStoreLookup[wh] = bookstore.New()
		}
	}
}

// Ingest all raw data from the fs into appropriate stores
// before the app starts
func init() {
	wg := &sync.WaitGroup{}
	wg.Add(len(Warehouses) - 1)

	// We're skipping "all" warehouse for iterations as
	// we're filling it on all other iterations
	for _, wh := range Warehouses[1:] {
		func(wh string) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Print(r)
				}
				wg.Done()
			}()

			dir := "./raw_books/" + wh

			bss := []*bookstore.Store{
				WarehouseStoreLookup["all"],
				WarehouseStoreLookup[wh],
			}

			agg := bookstore.MultiStoreAggregator{Stores: bss}
			loader := &bookstore.JSONBookLoader{
				Dirname:   dir,
				Warehouse: wh,
			}

			loader.IngestDir(agg)
		}(wh)
	}

	fmt.Printf("\nInitialised app with bookstores:\n\n")
	for _, wh := range Warehouses {
		bs := WarehouseStoreLookup[wh]
		fmt.Printf("Bookstore: '%s', books loaded: %d\n", wh, len(bs.Isbns.ToSlice()))
	}
	fmt.Printf("\n\n\n")
}

func main() {
	db := bookstore.New()

	// bookstore := bookstore.New()

	// if err := bookstore.IngestDir(os.Args[1], os.Args[2]); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// rb, err := bookstore.GetRandom().ToJSON()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// fmt.Printf("Here's a random book for you\n%s", rb)

	// note := notes.GenerateNew(bookstore, "in-note").GenerateOutput()

	// file, err := os.Create("Note-1.json")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// defer file.Close()

	// note.Print(file)

	// fmt.Printf("Random note generated and stored in to a file\n")

	// fmt.Printf("\n\nGREAT SUCCESS !!!\n\n")
}

// func storeBooks(category string, books []BookItem) {
// 	for i, b := range books {
// 		bId, err := b.ISBN()
// 		if errors.Is(err, ErrNoISBN) {
// 			fmt.Printf("No ISBN found for book '%s' by '%s'", b.VolumeInfo.Title, b.VolumeInfo.Authors[0])
// 			bId = fmt.Sprintf("%s-%d", category, i)
// 		}

// 		json, err := json.MarshalIndent(b, "", "  ")
// 		if err != nil {
// 			fmt.Printf("Error marshaling JSON for book '%s', Error:\n%v\n", b.VolumeInfo.Title, err)
// 			break
// 		}

// 		fp := fmt.Sprintf("fixtures/%s/%s.json", category, bId)
// 		if err = os.WriteFile(fp, json, 0755); err != nil {
// 			fmt.Printf("Error writing to '%s', Error:\n%v\n", fp, err)
// 		}
// 	}
// }
