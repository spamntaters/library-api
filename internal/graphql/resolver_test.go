package graphql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
	"github.com/user/library-api/internal/db"
	"github.com/user/library-api/internal/openlibrary"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateAuthor(ctx context.Context, arg db.CreateAuthorParams) (db.Author, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Author), args.Error(1)
}

func (m *MockStore) GetAuthor(ctx context.Context, id int32) (db.Author, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Author), args.Error(1)
}

func (m *MockStore) ListAuthors(ctx context.Context, arg db.ListAuthorsParams) ([]db.Author, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]db.Author), args.Error(1)
}

func (m *MockStore) UpdateAuthor(ctx context.Context, arg db.UpdateAuthorParams) (db.Author, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Author), args.Error(1)
}

func (m *MockStore) DeleteAuthor(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStore) GetAuthorByName(ctx context.Context, name string) (db.Author, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return db.Author{}, args.Error(1)
	}
	return args.Get(0).(db.Author), args.Error(1)
}

func (m *MockStore) GetBooksByAuthor(ctx context.Context, authorID int32) ([]db.Book, error) {
	args := m.Called(ctx, authorID)
	return args.Get(0).([]db.Book), args.Error(1)
}

func (m *MockStore) GetSeriesByAuthor(ctx context.Context, authorID int32) ([]db.Series, error) {
	args := m.Called(ctx, authorID)
	return args.Get(0).([]db.Series), args.Error(1)
}

func (m *MockStore) CreateBook(ctx context.Context, arg db.CreateBookParams) (db.Book, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Book), args.Error(1)
}

func (m *MockStore) GetBook(ctx context.Context, id int32) (db.Book, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Book), args.Error(1)
}

func (m *MockStore) ListBooks(ctx context.Context, arg db.ListBooksParams) ([]db.Book, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]db.Book), args.Error(1)
}

func (m *MockStore) UpdateBook(ctx context.Context, arg db.UpdateBookParams) (db.Book, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Book), args.Error(1)
}

func (m *MockStore) DeleteBook(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStore) ToggleOwned(ctx context.Context, id int32) (db.Book, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Book), args.Error(1)
}

func (m *MockStore) ToggleRead(ctx context.Context, id int32) (db.Book, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Book), args.Error(1)
}

func (m *MockStore) GetBookByISBN(ctx context.Context, isbn pgtype.Text) (db.Book, error) {
	args := m.Called(ctx, isbn)
	return args.Get(0).(db.Book), args.Error(1)
}

func (m *MockStore) CreateSeries(ctx context.Context, arg db.CreateSeriesParams) (db.Series, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Series), args.Error(1)
}

func (m *MockStore) GetSeries(ctx context.Context, id int32) (db.Series, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Series), args.Error(1)
}

func (m *MockStore) ListSeries(ctx context.Context, arg db.ListSeriesParams) ([]db.Series, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]db.Series), args.Error(1)
}

func (m *MockStore) UpdateSeries(ctx context.Context, arg db.UpdateSeriesParams) (db.Series, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Series), args.Error(1)
}

func (m *MockStore) DeleteSeries(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStore) AddBookToSeries(ctx context.Context, arg db.AddBookToSeriesParams) (db.SeriesBook, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.SeriesBook), args.Error(1)
}

func (m *MockStore) RemoveBookFromSeries(ctx context.Context, arg db.RemoveBookFromSeriesParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockStore) GetSeriesBooks(ctx context.Context, arg db.GetSeriesBooksParams) ([]db.GetSeriesBooksRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]db.GetSeriesBooksRow), args.Error(1)
}

func (m *MockStore) GetMissingBooks(ctx context.Context, arg db.GetMissingBooksParams) ([]db.Book, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]db.Book), args.Error(1)
}

func (m *MockStore) GetSeriesBooksByBookID(ctx context.Context, bookID int32) ([]db.SeriesBook, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).([]db.SeriesBook), args.Error(1)
}

func (m *MockStore) CreateTag(ctx context.Context, name string) (db.Tag, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(db.Tag), args.Error(1)
}

func (m *MockStore) GetTag(ctx context.Context, id int32) (db.Tag, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Tag), args.Error(1)
}

func (m *MockStore) ListTags(ctx context.Context, arg db.ListTagsParams) ([]db.Tag, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]db.Tag), args.Error(1)
}

func (m *MockStore) AddTagToBook(ctx context.Context, arg db.AddTagToBookParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockStore) RemoveTagFromBook(ctx context.Context, arg db.RemoveTagFromBookParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockStore) GetBookTags(ctx context.Context, bookID int32) ([]db.Tag, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).([]db.Tag), args.Error(1)
}

func (m *MockStore) WithTx(ctx context.Context, fn func(context.Context, db.Querier) error) error {
	args := m.Called(ctx, fn)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return fn(ctx, m)
}

type MockOLClient struct {
	mock.Mock
}

func (m *MockOLClient) LookupByISBN(isbn string) (*openlibrary.BookData, error) {
	args := m.Called(isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*openlibrary.BookData), args.Error(1)
}

func (m *MockOLClient) SearchByQuery(query string, limit int) (*openlibrary.SearchResponse, error) {
	args := m.Called(query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*openlibrary.SearchResponse), args.Error(1)
}

func (m *MockOLClient) GetBookByKey(key string) (*openlibrary.BookDetails, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*openlibrary.BookDetails), args.Error(1)
}

func testResolver(store *MockStore, olClient *MockOLClient) *Resolver {
	return &Resolver{
		Store:    store,
		OLClient: olClient,
	}
}

func strPtr(s string) *string {
	return &s
}

func intPtr(n int) *int {
	return &n
}

func boolPtr(b bool) *bool {
	return &b
}
