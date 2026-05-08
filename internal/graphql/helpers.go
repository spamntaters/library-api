package graphql

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/library-api/internal/db"
	"github.com/user/library-api/internal/openlibrary"
)

func dbToGraphQLSeries(s db.Series) *Series {
	result := &Series{
		ID:   int(s.ID),
		Name: s.Name,
	}
	if s.Description.Valid {
		result.Description = &s.Description.String
	}
	return result
}

func dbToGraphQLTag(t db.Tag) *Tag {
	return &Tag{
		ID:   int(t.ID),
		Name: t.Name,
	}
}

func dbToGraphQLAuthor(a db.Author) *Author {
	author := &Author{
		ID:   int(a.ID),
		Name: a.Name,
	}
	if a.Bio.Valid {
		author.Bio = &a.Bio.String
	}
	return author
}

func dbToGraphQLBook(b db.Book) *Book {
	book := &Book{
		ID:       int(b.ID),
		Title:    b.Title,
		AuthorID: int(b.AuthorID),
		Owned:    b.Owned,
		Read:     b.Read,
	}
	if b.Isbn.Valid {
		s := b.Isbn.String
		book.Isbn = &s
	}
	if b.Isbn13.Valid {
		s := b.Isbn13.String
		book.Isbn13 = &s
	}
	if b.PublishedDate.Valid {
		t := b.PublishedDate.Time
		book.PublishedDate = &t
	}
	if b.PageCount.Valid {
		n := int(b.PageCount.Int32)
		book.PageCount = &n
	}
	if b.Description.Valid {
		s := b.Description.String
		book.Description = &s
	}
	if b.CoverUrl.Valid {
		s := b.CoverUrl.String
		book.CoverURL = &s
	}
	return book
}

func olToExternalBook(data *openlibrary.BookData) *ExternalBook {
	result := &ExternalBook{}
	if data.Title != "" {
		result.Title = &data.Title
	}
	for _, a := range data.Authors {
		name := a.Name
		result.AuthorName = append(result.AuthorName, &name)
	}
	for _, isbn := range data.Identifiers.ISBN13 {
		s := isbn
		result.Isbn = append(result.Isbn, &s)
	}
	if len(result.Isbn) == 0 {
		for _, isbn := range data.Identifiers.ISBN10 {
			s := isbn
			result.Isbn = append(result.Isbn, &s)
		}
	}
	if data.URL != "" {
		result.Key = &data.URL
	}
	if data.Cover.Medium != "" {
		result.CoverID = &data.Cover.Medium
	}
	return result
}

func docToExternalBook(doc openlibrary.Doc) *ExternalBook {
	result := &ExternalBook{
		Title: &doc.Title,
	}
	if len(doc.AuthorName) > 0 {
		names := make([]*string, len(doc.AuthorName))
		for i, n := range doc.AuthorName {
			names[i] = &n
		}
		result.AuthorName = names
	}
	if len(doc.ISBN) > 0 {
		isbns := make([]*string, len(doc.ISBN))
		for i, isbn := range doc.ISBN {
			isbns[i] = &isbn
		}
		result.Isbn = isbns
	}
	if doc.FirstPublishYear > 0 {
		result.FirstPublishYear = &doc.FirstPublishYear
	}
	if doc.CoverID != "" {
		result.CoverID = &doc.CoverID
	}
	if doc.Key != "" {
		result.Key = &doc.Key
	}
	return result
}

func ptrToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrToBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func ptrToInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func timeToPgDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func intToPgInt4(i *int) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*i), Valid: true}
}
