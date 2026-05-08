package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error)
	GetAuthor(ctx context.Context, id int32) (Author, error)
	ListAuthors(ctx context.Context) ([]Author, error)
	UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) (Author, error)
	DeleteAuthor(ctx context.Context, id int32) error

	CreateBook(ctx context.Context, arg CreateBookParams) (Book, error)
	GetBook(ctx context.Context, id int32) (Book, error)
	ListBooks(ctx context.Context, arg ListBooksParams) ([]Book, error)
	UpdateBook(ctx context.Context, arg UpdateBookParams) (Book, error)
	DeleteBook(ctx context.Context, id int32) error
	ToggleOwned(ctx context.Context, id int32) (Book, error)
	GetBookByISBN(ctx context.Context, isbn pgtype.Text) (Book, error)

	CreateSeries(ctx context.Context, arg CreateSeriesParams) (Series, error)
	GetSeries(ctx context.Context, id int32) (Series, error)
	ListSeries(ctx context.Context) ([]Series, error)
	UpdateSeries(ctx context.Context, arg UpdateSeriesParams) (Series, error)
	DeleteSeries(ctx context.Context, id int32) error
	AddBookToSeries(ctx context.Context, arg AddBookToSeriesParams) (SeriesBook, error)
	RemoveBookFromSeries(ctx context.Context, arg RemoveBookFromSeriesParams) error
	GetSeriesBooks(ctx context.Context, seriesID int32) ([]GetSeriesBooksRow, error)
	GetMissingBooks(ctx context.Context, seriesID int32) ([]Book, error)

	CreateTag(ctx context.Context, name string) (Tag, error)
	GetTag(ctx context.Context, id int32) (Tag, error)
	ListTags(ctx context.Context) ([]Tag, error)
	AddTagToBook(ctx context.Context, arg AddTagToBookParams) error
	RemoveTagFromBook(ctx context.Context, arg RemoveTagFromBookParams) error
	GetBookTags(ctx context.Context, bookID int32) ([]Tag, error)
}

var _ Querier = (*Queries)(nil)
