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
)

func TestQuery_Books_NoFilters(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	books := []db.Book{
		{ID: 1, Title: "Book A", Owned: true},
		{ID: 2, Title: "Book B", Owned: false},
	}

	store.On("ListBooks", mock.Anything, mock.Anything).Return(books, nil)

	result, err := resolver.Query().Books(context.Background(), nil, nil, nil)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Book A", result[0].Title)
	assert.Equal(t, "Book B", result[1].Title)
}

func TestQuery_Books_FilterOwned(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	books := []db.Book{
		{ID: 1, Title: "Owned Book", Owned: true},
	}

	store.On("ListBooks", mock.Anything, mock.MatchedBy(func(arg db.ListBooksParams) bool {
		return arg.Owned.Valid && arg.Owned.Bool
	})).Return(books, nil)

	result, err := resolver.Query().Books(context.Background(), boolPtr(true), nil, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.True(t, result[0].Owned)
}

func TestQuery_Books_FilterByAuthor(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	books := []db.Book{
		{ID: 1, Title: "Author's Book", AuthorID: 5, Owned: true},
	}

	store.On("ListBooks", mock.Anything, mock.MatchedBy(func(arg db.ListBooksParams) bool {
		return arg.AuthorID.Valid && arg.AuthorID.Int32 == 5
	})).Return(books, nil)

	result, err := resolver.Query().Books(context.Background(), nil, intPtr(5), nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestQuery_Books_Error(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("ListBooks", mock.Anything, mock.Anything).Return([]db.Book{}, errors.New("db error"))

	result, err := resolver.Query().Books(context.Background(), nil, nil, nil)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestQuery_Book_ByID(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("GetBook", mock.Anything, int32(1)).Return(db.Book{
		ID:    1,
		Title: "Test Book",
		Owned: true,
	}, nil)

	result, err := resolver.Query().Book(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Test Book", result.Title)
	assert.Equal(t, 1, result.ID)
}

func TestQuery_Book_NotFound(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("GetBook", mock.Anything, int32(999)).Return(db.Book{}, errors.New("not found"))

	result, err := resolver.Query().Book(context.Background(), 999)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestQuery_Authors(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	authors := []db.Author{
		{ID: 1, Name: "Author A"},
		{ID: 2, Name: "Author B"},
	}

	store.On("ListAuthors", mock.Anything).Return(authors, nil)

	result, err := resolver.Query().Authors(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Author A", result[0].Name)
}

func TestQuery_Author_ByID(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("GetAuthor", mock.Anything, int32(1)).Return(db.Author{
		ID:   1,
		Name: "Test Author",
	}, nil)

	result, err := resolver.Query().Author(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Test Author", result.Name)
}

func TestQuery_Author_NotFound(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("GetAuthor", mock.Anything, int32(999)).Return(db.Author{}, errors.New("not found"))

	result, err := resolver.Query().Author(context.Background(), 999)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestQuery_Series(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	series := []db.Series{
		{ID: 1, Name: "Series A"},
		{ID: 2, Name: "Series B"},
	}

	store.On("ListSeries", mock.Anything).Return(series, nil)

	result, err := resolver.Query().Series(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Series A", result[0].Name)
}

func TestQuery_SeriesByID(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("GetSeries", mock.Anything, int32(1)).Return(db.Series{
		ID:   1,
		Name: "Test Series",
	}, nil)

	result, err := resolver.Query().SeriesByID(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Test Series", result.Name)
}

func TestQuery_SeriesByID_NotFound(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("GetSeries", mock.Anything, int32(999)).Return(db.Series{}, errors.New("not found"))

	result, err := resolver.Query().SeriesByID(context.Background(), 999)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestQuery_SeriesMissingBooks(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	missing := []db.Book{
		{ID: 1, Title: "Missing Book 1", Owned: false},
		{ID: 2, Title: "Missing Book 2", Owned: false},
	}

	store.On("GetMissingBooks", mock.Anything, int32(1)).Return(missing, nil)

	result, err := resolver.Query().SeriesMissingBooks(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.False(t, result[0].Owned)
}

func TestQuery_SeriesMissingBooks_Empty(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("GetMissingBooks", mock.Anything, int32(1)).Return([]db.Book{}, nil)

	result, err := resolver.Query().SeriesMissingBooks(context.Background(), 1)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestQuery_Tags(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	tags := []db.Tag{
		{ID: 1, Name: "fiction"},
		{ID: 2, Name: "sci-fi"},
	}

	store.On("ListTags", mock.Anything).Return(tags, nil)

	result, err := resolver.Query().Tags(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "fiction", result[0].Name)
}

func TestMutation_CreateBook(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("CreateBook", mock.Anything, mock.MatchedBy(func(arg db.CreateBookParams) bool {
		return arg.Title == "New Book" && arg.AuthorID == 1
	})).Return(db.Book{
		ID:       1,
		Title:    "New Book",
		AuthorID: 1,
		Owned:    false,
	}, nil)

	input := CreateBookInput{
		Title:    "New Book",
		AuthorID: 1,
	}

	result, err := resolver.Mutation().CreateBook(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "New Book", result.Title)
	assert.False(t, result.Owned)
}

func TestMutation_UpdateBook(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("UpdateBook", mock.Anything, mock.Anything).Return(db.Book{
		ID:    1,
		Title: "Updated Title",
		Owned: true,
	}, nil)

	input := UpdateBookInput{
		Title: strPtr("Updated Title"),
	}

	result, err := resolver.Mutation().UpdateBook(context.Background(), 1, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Title", result.Title)
}

func TestMutation_DeleteBook(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("DeleteBook", mock.Anything, int32(1)).Return(nil)

	result, err := resolver.Mutation().DeleteBook(context.Background(), 1)

	require.NoError(t, err)
	assert.True(t, result)
}

func TestMutation_ToggleOwned(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("ToggleOwned", mock.Anything, int32(1)).Return(db.Book{
		ID:    1,
		Title: "Book",
		Owned: true,
	}, nil)

	result, err := resolver.Mutation().ToggleOwned(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Owned)
}

func TestMutation_CreateAuthor(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("CreateAuthor", mock.Anything, db.CreateAuthorParams{
		Name: "New Author",
		Bio:  pgtype.Text{Valid: false},
	}).Return(db.Author{
		ID:   1,
		Name: "New Author",
	}, nil)

	input := CreateAuthorInput{
		Name: "New Author",
	}

	result, err := resolver.Mutation().CreateAuthor(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "New Author", result.Name)
}

func TestMutation_UpdateAuthor(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("UpdateAuthor", mock.Anything, mock.Anything).Return(db.Author{
		ID:   1,
		Name: "Updated Name",
	}, nil)

	input := UpdateAuthorInput{
		Name: strPtr("Updated Name"),
	}

	result, err := resolver.Mutation().UpdateAuthor(context.Background(), 1, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Name", result.Name)
}

func TestMutation_DeleteAuthor(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("DeleteAuthor", mock.Anything, int32(1)).Return(nil)

	result, err := resolver.Mutation().DeleteAuthor(context.Background(), 1)

	require.NoError(t, err)
	assert.True(t, result)
}

func TestMutation_CreateSeries(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("CreateSeries", mock.Anything, db.CreateSeriesParams{
		Name:        "New Series",
		Description: pgtype.Text{Valid: false},
	}).Return(db.Series{
		ID:   1,
		Name: "New Series",
	}, nil)

	input := CreateSeriesInput{
		Name: "New Series",
	}

	result, err := resolver.Mutation().CreateSeries(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "New Series", result.Name)
}

func TestMutation_UpdateSeries(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("UpdateSeries", mock.Anything, mock.Anything).Return(db.Series{
		ID:   1,
		Name: "Updated Series",
	}, nil)

	input := UpdateSeriesInput{
		Name: strPtr("Updated Series"),
	}

	result, err := resolver.Mutation().UpdateSeries(context.Background(), 1, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Series", result.Name)
}

func TestMutation_DeleteSeries(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("DeleteSeries", mock.Anything, int32(1)).Return(nil)

	result, err := resolver.Mutation().DeleteSeries(context.Background(), 1)

	require.NoError(t, err)
	assert.True(t, result)
}

func TestMutation_AddBookToSeries(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("AddBookToSeries", mock.Anything, db.AddBookToSeriesParams{
		SeriesID: 1,
		BookID:   2,
		Position: 3,
	}).Return(db.SeriesBook{
		SeriesID: 1,
		BookID:   2,
		Position: 3,
	}, nil)

	result, err := resolver.Mutation().AddBookToSeries(context.Background(), 2, 1, 3)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 3, result.Position)
}

func TestMutation_RemoveBookFromSeries(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("RemoveBookFromSeries", mock.Anything, db.RemoveBookFromSeriesParams{
		SeriesID: 1,
		BookID:   2,
	}).Return(nil)

	result, err := resolver.Mutation().RemoveBookFromSeries(context.Background(), 2, 1)

	require.NoError(t, err)
	assert.True(t, result)
}

func TestMutation_CreateTag(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("CreateTag", mock.Anything, "sci-fi").Return(db.Tag{
		ID:   1,
		Name: "sci-fi",
	}, nil)

	result, err := resolver.Mutation().CreateTag(context.Background(), "sci-fi")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "sci-fi", result.Name)
}

func TestMutation_AddTagToBook(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("AddTagToBook", mock.Anything, db.AddTagToBookParams{
		BookID: 1,
		TagID:  2,
	}).Return(nil)

	store.On("GetBook", mock.Anything, int32(1)).Return(db.Book{
		ID:    1,
		Title: "Tagged Book",
		Owned: true,
	}, nil)

	result, err := resolver.Mutation().AddTagToBook(context.Background(), 1, 2)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Tagged Book", result.Title)
}

func TestMutation_RemoveTagFromBook(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("RemoveTagFromBook", mock.Anything, db.RemoveTagFromBookParams{
		BookID: 1,
		TagID:  2,
	}).Return(nil)

	store.On("GetBook", mock.Anything, int32(1)).Return(db.Book{
		ID:    1,
		Title: "Untagged Book",
		Owned: true,
	}, nil)

	result, err := resolver.Mutation().RemoveTagFromBook(context.Background(), 1, 2)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Untagged Book", result.Title)
}
