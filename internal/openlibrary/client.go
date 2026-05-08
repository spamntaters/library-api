package openlibrary

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client interface {
	LookupByISBN(isbn string) (*BookData, error)
	SearchByQuery(query string, limit int) (*SearchResponse, error)
	GetBookByKey(key string) (*BookDetails, error)
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) Client {
	return &client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *client) LookupByISBN(isbn string) (*BookData, error) {
	params := url.Values{}
	params.Set("bibkeys", "ISBN:"+isbn)
	params.Set("format", "json")
	params.Set("jscmd", "data")

	reqURL := fmt.Sprintf("%s/api/books?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("openlibrary request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openlibrary returned %d: %s", resp.StatusCode, string(body))
	}

	var raw map[string]*BookData
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode openlibrary response: %w", err)
	}

	key := "ISBN:" + isbn
	book, ok := raw[key]
	if !ok {
		return nil, nil
	}

	return book, nil
}

func (c *client) SearchByQuery(query string, limit int) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("limit", fmt.Sprintf("%d", limit))

	reqURL := fmt.Sprintf("%s/search.json?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("openlibrary search failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openlibrary returned %d: %s", resp.StatusCode, string(body))
	}

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode openlibrary search response: %w", err)
	}

	return &result, nil
}

func (c *client) GetBookByKey(key string) (*BookDetails, error) {
	reqURL := fmt.Sprintf("%s%s.json", c.baseURL, key)

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("openlibrary book lookup failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openlibrary returned %d: %s", resp.StatusCode, string(body))
	}

	var result BookDetails
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode openlibrary book response: %w", err)
	}

	return &result, nil
}

type BookData struct {
	Title        string     `json:"title"`
	Authors      []AuthorOL `json:"authors"`
	Publishers   []PublisherOL `json:"publishers"`
	PublishDate  string     `json:"publish_date"`
	NumberOfPages int       `json:"number_of_pages"`
	Cover        CoverOL    `json:"cover"`
	URL          string     `json:"url"`
	Identifiers  IdentifiersOL `json:"identifiers"`
}

type AuthorOL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PublisherOL struct {
	Name string `json:"name"`
}

type CoverOL struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type IdentifiersOL struct {
	ISBN  []string `json:"isbn"`
	ISBN10 []string `json:"isbn_10"`
	ISBN13 []string `json:"isbn_13"`
	LCCN  []string `json:"lccn"`
	OCLC  []string `json:"oclc"`
}

type SearchResponse struct {
	NumFound int     `json:"numFound"`
	Docs     []Doc   `json:"docs"`
}

type Doc struct {
	Title            string   `json:"title"`
	AuthorName       []string `json:"author_name"`
	ISBN             []string `json:"isbn"`
	FirstPublishYear int      `json:"first_publish_year"`
	CoverID          string   `json:"cover_i"`
	Key              string   `json:"key"`
}

type BookDetails struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	PublishDate string   `json:"publish_date"`
	ISBNs       []string `json:"isbn_13"`
	Pages       int      `json:"number_of_pages"`
	CoverID     string   `json:"covers"`
}
