package openlibrary

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func skipIfOffline(t *testing.T) {
	if os.Getenv("TEST_OL_API") == "" {
		t.Skip("set TEST_OL_API=1 to run Open Library integration tests")
	}
}

func TestClient_LookupByISBN_Real(t *testing.T) {
	skipIfOffline(t)

	client := NewClient("https://openlibrary.org")

	result, err := client.LookupByISBN("9780743273565")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Title)
	assert.NotEmpty(t, result.Authors)
}

func TestClient_SearchByQuery_Real(t *testing.T) {
	skipIfOffline(t)

	client := NewClient("https://openlibrary.org")

	result, err := client.SearchByQuery("dune frank herbert", 5)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.NumFound, 0)
	assert.NotEmpty(t, result.Docs)
}

func TestClient_GetBookByKey_Real(t *testing.T) {
	skipIfOffline(t)

	client := NewClient("https://openlibrary.org")

	result, err := client.GetBookByKey("/books/OL26298570M")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Title)
}

func TestClient_LookupByISBN_NotFound(t *testing.T) {
	skipIfOffline(t)

	client := NewClient("https://openlibrary.org")

	result, err := client.LookupByISBN("0000000000000")

	require.NoError(t, err)
	assert.Nil(t, result)
}
