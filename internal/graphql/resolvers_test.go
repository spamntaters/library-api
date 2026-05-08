package graphql

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
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

	result, err := resolver.Query().Books(context.Background(), nil, nil, nil, nil, nil)

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

	result, err := resolver.Query().Books(context.Background(), boolPtr(true), nil, nil, nil, nil)

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

	result, err := resolver.Query().Books(context.Background(), nil, intPtr(5), nil, nil, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestQuery_Books_FilterByTag(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	books := []db.Book{
		{ID: 1, Title: "Tagged Book", AuthorID: 1, Owned: true},
	}

	store.On("ListBooks", mock.Anything, mock.MatchedBy(func(arg db.ListBooksParams) bool {
		return arg.TagID.Valid && arg.TagID.Int32 == 3
	})).Return(books, nil)

	result, err := resolver.Query().Books(context.Background(), nil, nil, intPtr(3), nil, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Tagged Book", result[0].Title)
}

func TestQuery_Books_Error(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("ListBooks", mock.Anything, mock.Anything).Return([]db.Book{}, errors.New("db error"))

	result, err := resolver.Query().Books(context.Background(), nil, nil, nil, nil, nil)

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

	store.On("ListAuthors", mock.Anything, mock.Anything).Return(authors, nil)

	result, err := resolver.Query().Authors(context.Background(), nil, nil)

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

	store.On("ListSeries", mock.Anything, mock.Anything).Return(series, nil)

	result, err := resolver.Query().Series(context.Background(), nil, nil)

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

	store.On("ListTags", mock.Anything, mock.Anything).Return(tags, nil)

	result, err := resolver.Query().Tags(context.Background(), nil, nil)

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

func TestValidation_CreateSeries_DuplicateName(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("CreateSeries", mock.Anything, mock.Anything).Return(db.Series{}, &pgconn.PgError{Code: "23505"})

	input := CreateSeriesInput{
		Name: "Existing Series",
	}

	result, err := resolver.Mutation().CreateSeries(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsDuplicateError(err))
	assert.Contains(t, err.Error(), "name")
	assert.Contains(t, err.Error(), "Existing Series")
}

func TestValidation_CreateTag_DuplicateName(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("CreateTag", mock.Anything, "existing-tag").Return(db.Tag{}, &pgconn.PgError{Code: "23505"})

	result, err := resolver.Mutation().CreateTag(context.Background(), "existing-tag")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsDuplicateError(err))
	assert.Contains(t, err.Error(), "name")
	assert.Contains(t, err.Error(), "existing-tag")
}

func TestValidation_UpdateSeries_DuplicateName(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	newName := "Duplicate Series"
	store.On("UpdateSeries", mock.Anything, mock.Anything).Return(db.Series{}, &pgconn.PgError{Code: "23505"})

	input := UpdateSeriesInput{
		Name: &newName,
	}

	result, err := resolver.Mutation().UpdateSeries(context.Background(), 1, input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsDuplicateError(err))
	assert.Contains(t, err.Error(), "Duplicate Series")
}

func TestValidation_CreateBook_EmptyTitle(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	input := CreateBookInput{
		Title:    "",
		AuthorID: 1,
	}

	result, err := resolver.Mutation().CreateBook(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsValidationError(err))
	assert.Contains(t, err.Error(), "title")
}

func TestValidation_CreateBook_ZeroAuthorID(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	input := CreateBookInput{
		Title:    "Valid Title",
		AuthorID: 0,
	}

	result, err := resolver.Mutation().CreateBook(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsValidationError(err))
	assert.Contains(t, err.Error(), "authorID")
}

func TestValidation_CreateAuthor_EmptyName(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	input := CreateAuthorInput{
		Name: "",
	}

	result, err := resolver.Mutation().CreateAuthor(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsValidationError(err))
	assert.Contains(t, err.Error(), "name")
}

func TestValidation_CreateSeries_EmptyName(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	input := CreateSeriesInput{
		Name: "   ",
	}

	result, err := resolver.Mutation().CreateSeries(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsValidationError(err))
	assert.Contains(t, err.Error(), "name")
}

func TestValidation_CreateTag_EmptyName(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	result, err := resolver.Mutation().CreateTag(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsValidationError(err))
	assert.Contains(t, err.Error(), "name")
}

func TestValidation_DeleteBook_ZeroID(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	result, err := resolver.Mutation().DeleteBook(context.Background(), 0)

	require.Error(t, err)
	assert.False(t, result)
	assert.True(t, IsValidationError(err))
}

func TestValidation_UpdateBook_NegativeID(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	input := UpdateBookInput{
		Title: strPtr("Updated"),
	}

	result, err := resolver.Mutation().UpdateBook(context.Background(), -1, input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsValidationError(err))
}

func TestValidation_AddBookToSeries_ZeroBookID(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	result, err := resolver.Mutation().AddBookToSeries(context.Background(), 0, 1, 1)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsValidationError(err))
	assert.Contains(t, err.Error(), "bookID")
}

func TestValidation_AddTagToBook_ZeroTagID(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	result, err := resolver.Mutation().AddTagToBook(context.Background(), 1, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsValidationError(err))
	assert.Contains(t, err.Error(), "tagID")
}

func TestValidation_UpdateBook_EmptyTitle(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	emptyTitle := ""
	input := UpdateBookInput{
		Title: &emptyTitle,
	}

	result, err := resolver.Mutation().UpdateBook(context.Background(), 1, input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, IsValidationError(err))
	assert.Contains(t, err.Error(), "title")
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

	store.On("GetBook", mock.Anything, int32(2)).Return(db.Book{
		ID:    2,
		Title: "Linked Book",
		Owned: true,
	}, nil)

	store.On("GetSeries", mock.Anything, int32(1)).Return(db.Series{
		ID:   1,
		Name: "Linked Series",
	}, nil)

	result, err := resolver.Mutation().AddBookToSeries(context.Background(), 2, 1, 3)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 3, result.Position)
	assert.Equal(t, "Linked Book", result.Book.Title)
	assert.Equal(t, "Linked Series", result.Series.Name)
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

func TestQuery_Books_WithPagination(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	books := []db.Book{
		{ID: 1, Title: "Book A", Owned: true},
		{ID: 2, Title: "Book B", Owned: false},
	}

	store.On("ListBooks", mock.Anything, mock.MatchedBy(func(arg db.ListBooksParams) bool {
		return arg.Limit.Valid && arg.Limit.Int32 == 10 && arg.Offset.Valid && arg.Offset.Int32 == 5
	})).Return(books, nil)

	result, err := resolver.Query().Books(context.Background(), nil, nil, nil, intPtr(10), intPtr(5))

	require.NoError(t, err)
	require.Len(t, result, 2)
}

func TestQuery_Authors_WithPagination(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	authors := []db.Author{
		{ID: 1, Name: "Author A"},
	}

	store.On("ListAuthors", mock.Anything, mock.MatchedBy(func(arg db.ListAuthorsParams) bool {
		return arg.Limit.Valid && arg.Limit.Int32 == 5 && arg.Offset.Valid && arg.Offset.Int32 == 0
	})).Return(authors, nil)

	result, err := resolver.Query().Authors(context.Background(), intPtr(5), intPtr(0))

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestQuery_Series_WithPagination(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	series := []db.Series{
		{ID: 1, Name: "Series A"},
	}

	store.On("ListSeries", mock.Anything, mock.MatchedBy(func(arg db.ListSeriesParams) bool {
		return arg.Limit.Valid && arg.Limit.Int32 == 3 && arg.Offset.Valid && arg.Offset.Int32 == 1
	})).Return(series, nil)

	result, err := resolver.Query().Series(context.Background(), intPtr(3), intPtr(1))

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestQuery_Tags_WithPagination(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	tags := []db.Tag{
		{ID: 1, Name: "fiction"},
	}

	store.On("ListTags", mock.Anything, mock.MatchedBy(func(arg db.ListTagsParams) bool {
		return arg.Limit.Valid && arg.Limit.Int32 == 20 && !arg.Offset.Valid
	})).Return(tags, nil)

	result, err := resolver.Query().Tags(context.Background(), intPtr(20), nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestEntity_Books_Author(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("GetAuthor", mock.Anything, int32(1)).Return(db.Author{
		ID:   1,
		Name: "Test Author",
	}, nil)

	book := &Book{ID: 1, Title: "Test Book", AuthorID: 1}

	result, err := resolver.Book().Author(context.Background(), book)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Test Author", result.Name)
}

func TestEntity_Books_Tags(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	tags := []db.Tag{
		{ID: 1, Name: "fiction"},
		{ID: 2, Name: "sci-fi"},
	}

	store.On("GetBookTags", mock.Anything, int32(1)).Return(tags, nil)

	book := &Book{ID: 1, Title: "Test Book", AuthorID: 1}

	result, err := resolver.Book().Tags(context.Background(), book)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "fiction", result[0].Name)
}

func TestEntity_Books_Series(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	seriesBooks := []db.SeriesBook{
		{SeriesID: 1, BookID: 1, Position: 1},
	}

	store.On("GetSeriesBooksByBookID", mock.Anything, int32(1)).Return(seriesBooks, nil)
	store.On("GetSeries", mock.Anything, int32(1)).Return(db.Series{
		ID:   1,
		Name: "Test Series",
	}, nil)

	book := &Book{ID: 1, Title: "Test Book", AuthorID: 1}

	result, err := resolver.Book().Series(context.Background(), book)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 1, result[0].Position)
	assert.Equal(t, "Test Series", result[0].Series.Name)
}

func TestEntity_Author_Books(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	books := []db.Book{
		{ID: 1, Title: "Book A", AuthorID: 1, Owned: true},
		{ID: 2, Title: "Book B", AuthorID: 1, Owned: false},
	}

	store.On("GetBooksByAuthor", mock.Anything, int32(1)).Return(books, nil)

	author := &Author{ID: 1, Name: "Test Author"}

	result, err := resolver.Author().Books(context.Background(), author)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Book A", result[0].Title)
}

func TestEntity_Author_Series(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	series := []db.Series{
		{ID: 1, Name: "Series A"},
	}

	store.On("GetSeriesByAuthor", mock.Anything, int32(1)).Return(series, nil)

	author := &Author{ID: 1, Name: "Test Author"}

	result, err := resolver.Author().Series(context.Background(), author)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Series A", result[0].Name)
}

func TestEntity_Series_Books(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	seriesBooks := []db.GetSeriesBooksRow{
		{SeriesID: 1, BookID: 1, Position: 1},
	}

	store.On("GetSeriesBooks", mock.Anything, int32(1)).Return(seriesBooks, nil)
	store.On("GetBook", mock.Anything, int32(1)).Return(db.Book{
		ID:       1,
		Title:    "Test Book",
		AuthorID: 1,
		Owned:    true,
	}, nil)

	series := &Series{ID: 1, Name: "Test Series"}

	result, err := resolver.Series().Books(context.Background(), series)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Test Book", result[0].Book.Title)
	assert.Equal(t, 1, result[0].Position)
}

func TestEntity_SeriesBook_Book(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("GetBook", mock.Anything, int32(1)).Return(db.Book{
		ID:       1,
		Title:    "Test Book",
		AuthorID: 1,
		Owned:    true,
	}, nil)

	sb := &SeriesBook{BookID: 1, SeriesID: 1, Position: 1}

	result, err := resolver.SeriesBook().Book(context.Background(), sb)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Test Book", result.Title)
}

func TestEntity_SeriesBook_Series(t *testing.T) {
	store := &MockStore{}
	resolver := testResolver(store, &MockOLClient{})

	store.On("GetSeries", mock.Anything, int32(1)).Return(db.Series{
		ID:   1,
		Name: "Test Series",
	}, nil)

	sb := &SeriesBook{BookID: 1, SeriesID: 1, Position: 1}

	result, err := resolver.SeriesBook().Series(context.Background(), sb)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Test Series", result.Name)
}
