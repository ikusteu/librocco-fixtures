package notes

import (
	"encoding/json"
	"io"

	"github.com/ikusteu/librocco-fixtures/pkg/books"
	"github.com/ikusteu/librocco-fixtures/pkg/bookstore"
	"github.com/ikusteu/librocco-fixtures/pkg/utils"
)

type NoteInternal struct {
	_type string
	books map[string]*books.BookStock
}
type NoteOutput struct {
	Type  string             `json:"type"`
	Books []*books.BookStock `json:"books"`
}

func New(_type string) *NoteInternal {
	books := make(map[string]*books.BookStock)
	return &NoteInternal{_type, books}
}

func GenerateNew(bStore *bookstore.Store, _type string) *NoteInternal {
	note := New(_type)

	nBooks := utils.RandInt(5, 20)

	for i := 0; i < nBooks; i++ {
		b := bStore.GetRandom()
		q := utils.RandInt(1, 10)

		note.AddBook(b, q)
	}

	return note
}

func (n *NoteInternal) AddBook(b *books.BookStock, q int) error {
	isbn, err := b.ISBN10()
	if err != nil {
		return err
	}

	if bs, ok := n.books[isbn]; ok {
		bs.Add(q)
		return nil
	}

	n.books[isbn] = b
	b.Add(q)

	return nil
}

func (ni *NoteInternal) GenerateOutput() *NoteOutput {
	no := &NoteOutput{Type: ni._type}

	for _, b := range ni.books {
		no.Books = append(no.Books, b)
	}

	return no
}

func (no *NoteOutput) Print(w io.Writer) error {
	buf, err := json.MarshalIndent(no, "", "  ")
	if err != nil {
		return err
	}

	w.Write(buf)
	return nil
}
