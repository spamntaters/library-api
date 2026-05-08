package graphql

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/user/library-api/internal/db"
	"github.com/user/library-api/internal/openlibrary"
)

func TestLookupByIsbn_Success(t *testing.T) {
	store := &MockStore{}
	olClient := &MockOLClient{}
	resolver := testResolver(store, olClient)

	olData := &openlibrary.BookData{
		Title: "The Great Gatsby",
		Authors: []openlibrary.AuthorOL{
			{Name: "F. Scott Fitzgerald"},
		},
		Identifiers: openlibrary.IdentifiersOL{
			ISBN13: []string{"9780743273565"},
		},
		Cover: openlibrary.CoverOL{
			Medium: "https://covers.openlibrary.org/b/id/12345-M.jpg",
		},
	}

	olClient.On("LookupByISBN", "9780743273565").Return(olData, nil)

	result, err := resolver.Query().LookupByIsbn(context.Background(), "9780743273565")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "The Great Gatsby", *result.Title)
	assert.Len(t, result.AuthorName, 1)
	assert.Equal(t, "F. Scott Fitzgerald", *result.AuthorName[0])
	assert.Equal(t, "https://covers.openlibrary.org/b/id/12345-M.jpg", *result.CoverID)

	olClient.AssertExpectations(t)
}

func TestLookupByIsbn_NotFound(t *testing.T) {
	store := &MockStore{}
	olClient := &MockOLClient{}
	resolver := testResolver(store, olClient)

	olClient.On("LookupByISBN", "0000000000000").Return((*openlibrary.BookData)(nil), nil)

	result, err := resolver.Query().LookupByIsbn(context.Background(), "0000000000000")

	require.NoError(t, err)
	assert.Nil(t, result)

	olClient.AssertExpectations(t)
}

func TestLookupByIsbn_APIError(t *testing.T) {
	store := &MockStore{}
	olClient := &MockOLClient{}
	resolver := testResolver(store, olClient)

	olClient.On("LookupByISBN", "bad-isbn").Return((*openlibrary.BookData)(nil), errors.New("network error"))

	result, err := resolver.Query().LookupByIsbn(context.Background(), "bad-isbn")

	require.Error(t, err)
	assert.Nil(t, result)

	olClient.AssertExpectations(t)
}

func TestScanAndAddBook_AlreadyExists(t *testing.T) {
	store := &MockStore{}
	olClient := &MockOLClient{}
	resolver := testResolver(store, olClient)

	existingBook := db.Book{
		ID:      1,
		Title:   "Existing Book",
		Owned:   true,
		Isbn:    pgtype.Text{String: "9780743273565", Valid: true},
		Isbn13:  pgtype.Text{String: "9780743273565", Valid: true},
	}

	store.On("GetBookByISBN", mock.Anything, pgtype.Text{String: "9780743273565", Valid: true}).Return(existingBook, nil)

	result, err := resolver.Mutation().ScanAndAddBook(context.Background(), "9780743273565")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, *result.AlreadyExists)
	assert.Equal(t, "Existing Book", result.Book.Title)
	assert.Equal(t, 1, result.Book.ID)

	// Should NOT call Open Library or create anything
	olClient.AssertNotCalled(t, "LookupByISBN")
	store.AssertNotCalled(t, "CreateAuthor")
	store.AssertNotCalled(t, "CreateBook")
}

func TestScanAndAddBook_NewBook(t *testing.T) {
	store := &MockStore{}
	olClient := &MockOLClient{}
	resolver := testResolver(store, olClient)

	olData := &openlibrary.BookData{
		Title:         "Dune",
		PublishDate:   "1965",
		NumberOfPages: 412,
		Authors: []openlibrary.AuthorOL{
			{Name: "Frank Herbert"},
		},
		Identifiers: openlibrary.IdentifiersOL{
			ISBN10: []string{"0441172717"},
			ISBN13: []string{"9780441172719"},
		},
		Cover: openlibrary.CoverOL{
			Medium: "https://covers.openlibrary.org/b/id/67890-M.jpg",
		},
	}

	store.On("GetBookByISBN", mock.Anything, pgtype.Text{String: "9780441172719", Valid: true}).Return(db.Book{}, errors.New("not found"))
	store.On("WithTx", mock.Anything, mock.Anything).Return(nil)

	store.On("GetAuthorByName", mock.Anything, "Frank Herbert").Return(db.Author{}, errors.New("not found"))

	store.On("CreateAuthor", mock.Anything, db.CreateAuthorParams{
		Name: "Frank Herbert",
		Bio:  pgtype.Text{Valid: false},
	}).Return(db.Author{
		ID:   1,
		Name: "Frank Herbert",
	}, nil)

	store.On("CreateBook", mock.Anything, mock.MatchedBy(func(arg db.CreateBookParams) bool {
		return arg.Title == "Dune" &&
			arg.AuthorID == 1 &&
			arg.Isbn.String == "0441172717" &&
			arg.Isbn13.String == "9780441172719" &&
			arg.Owned == true
	})).Return(db.Book{
		ID:       1,
		Title:    "Dune",
		AuthorID: 1,
		Isbn:     pgtype.Text{String: "0441172717", Valid: true},
		Isbn13:   pgtype.Text{String: "9780441172719", Valid: true},
		Owned:    true,
	}, nil)

	olClient.On("LookupByISBN", "9780441172719").Return(olData, nil)

	result, err := resolver.Mutation().ScanAndAddBook(context.Background(), "9780441172719")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, *result.AlreadyExists)
	require.NotNil(t, result.Book)
	assert.Equal(t, "Dune", result.Book.Title)
	assert.Equal(t, 1, result.Book.ID)
	assert.True(t, result.Book.Owned)
	require.NotNil(t, result.External)
	assert.Equal(t, "Dune", *result.External.Title)

	store.AssertExpectations(t)
	olClient.AssertExpectations(t)
}

func TestScanAndAddBook_ExistingAuthor(t *testing.T) {
	store := &MockStore{}
	olClient := &MockOLClient{}
	resolver := testResolver(store, olClient)

	olData := &openlibrary.BookData{
		Title:         "Dune Messiah",
		PublishDate:   "1969",
		NumberOfPages: 256,
		Authors: []openlibrary.AuthorOL{
			{Name: "Frank Herbert"},
		},
		Identifiers: openlibrary.IdentifiersOL{
			ISBN13: []string{"9780441172726"},
		},
		Cover: openlibrary.CoverOL{
			Medium: "https://covers.openlibrary.org/b/id/67891-M.jpg",
		},
	}

	store.On("GetBookByISBN", mock.Anything, pgtype.Text{String: "9780441172726", Valid: true}).Return(db.Book{}, errors.New("not found"))
	store.On("WithTx", mock.Anything, mock.Anything).Return(nil)

	store.On("GetAuthorByName", mock.Anything, "Frank Herbert").Return(db.Author{
		ID:   1,
		Name: "Frank Herbert",
	}, nil)

	store.On("CreateBook", mock.Anything, mock.MatchedBy(func(arg db.CreateBookParams) bool {
		return arg.Title == "Dune Messiah" &&
			arg.AuthorID == 1 &&
			arg.Isbn13.String == "9780441172726" &&
			arg.Owned == true
	})).Return(db.Book{
		ID:       2,
		Title:    "Dune Messiah",
		AuthorID: 1,
		Isbn13:   pgtype.Text{String: "9780441172726", Valid: true},
		Owned:    true,
	}, nil)

	olClient.On("LookupByISBN", "9780441172726").Return(olData, nil)

	result, err := resolver.Mutation().ScanAndAddBook(context.Background(), "9780441172726")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, *result.AlreadyExists)
	require.NotNil(t, result.Book)
	assert.Equal(t, "Dune Messiah", result.Book.Title)
	assert.Equal(t, 2, result.Book.ID)

	// Should NOT call CreateAuthor since author already exists
	store.AssertNotCalled(t, "CreateAuthor")
	store.AssertExpectations(t)
	olClient.AssertExpectations(t)
}

func TestScanAndAddBook_TxRollback(t *testing.T) {
	store := &MockStore{}
	olClient := &MockOLClient{}
	resolver := testResolver(store, olClient)

	olData := &openlibrary.BookData{
		Title:   "Failing Book",
		Authors: []openlibrary.AuthorOL{{Name: "Test Author"}},
		Identifiers: openlibrary.IdentifiersOL{
			ISBN13: []string{"9780000000000"},
		},
	}

	store.On("GetBookByISBN", mock.Anything, pgtype.Text{String: "9780000000000", Valid: true}).Return(db.Book{}, errors.New("not found"))
	store.On("WithTx", mock.Anything, mock.Anything).Return(nil)
	olClient.On("LookupByISBN", "9780000000000").Return(olData, nil)
	store.On("GetAuthorByName", mock.Anything, "Test Author").Return(db.Author{}, errors.New("not found"))
	store.On("CreateAuthor", mock.Anything, db.CreateAuthorParams{
		Name: "Test Author",
		Bio:  pgtype.Text{Valid: false},
	}).Return(db.Author{ID: 1, Name: "Test Author"}, nil)
	store.On("CreateBook", mock.Anything, mock.Anything).Return(db.Book{}, errors.New("db error"))

	result, err := resolver.Mutation().ScanAndAddBook(context.Background(), "9780000000000")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "db error", err.Error())

	store.AssertExpectations(t)
	olClient.AssertExpectations(t)
}

func TestScanAndAddBook_NotFoundInOpenLibrary(t *testing.T) {
	store := &MockStore{}
	olClient := &MockOLClient{}
	resolver := testResolver(store, olClient)

	store.On("GetBookByISBN", mock.Anything, mock.Anything).Return(db.Book{}, errors.New("not found"))
	olClient.On("LookupByISBN", "0000000000000").Return((*openlibrary.BookData)(nil), nil)

	result, err := resolver.Mutation().ScanAndAddBook(context.Background(), "0000000000000")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, *result.AlreadyExists)
	assert.Nil(t, result.Book)
	assert.Nil(t, result.External)

	// Should NOT create anything
	store.AssertNotCalled(t, "CreateAuthor")
	store.AssertNotCalled(t, "CreateBook")
}

func TestScanAndAddBook_OLClientError(t *testing.T) {
	store := &MockStore{}
	olClient := &MockOLClient{}
	resolver := testResolver(store, olClient)

	store.On("GetBookByISBN", mock.Anything, mock.Anything).Return(db.Book{}, errors.New("not found"))
	olClient.On("LookupByISBN", "bad-isbn").Return((*openlibrary.BookData)(nil), errors.New("API error"))

	result, err := resolver.Mutation().ScanAndAddBook(context.Background(), "bad-isbn")

	require.Error(t, err)
	assert.Nil(t, result)

	store.AssertNotCalled(t, "CreateAuthor")
	store.AssertNotCalled(t, "CreateBook")
}
