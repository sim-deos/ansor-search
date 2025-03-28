package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Searcher struct {
	client   *http.Client
	engineId string
	apiKey   string
}

func NewSearcher(apiKey string, engineId string) Searcher {
	if apiKey == "" {
		panic("apiKey cannot be empty")
	}
	if engineId == "" {
		panic("engineId cannot be empty")
	}

	return Searcher{
		client:   &http.Client{},
		engineId: engineId,
		apiKey:   apiKey,
	}
}

func NewSearcherWithClient(client *http.Client, apiKey string, engineId string) Searcher {
	if apiKey == "" {
		panic("apiKey cannot be empty")
	}
	if engineId == "" {
		panic("engineId cannot be empty")
	}
	if client == nil {
		panic("client cannot be nil")
	}
	return Searcher{
		client:   client,
		engineId: engineId,
		apiKey:   apiKey,
	}
}

func (schr Searcher) Search(query string) (*SearchResponse, error) {
	canonicalizedQuery := canonicalizeQuery(query)
	url := fmt.Sprintf(
		"https://customsearch.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s",
		schr.apiKey, schr.engineId, canonicalizedQuery,
	)

	resp, err := schr.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get seach results: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to perform search, got %d: %s", resp.StatusCode, string(bytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read response body. Length: %d, StatusCode: %d: %w",
			resp.ContentLength, resp.StatusCode, err,
		)
	}

	var sr SearchResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body:\n%s: %w", string(body), err)
	}

	return &sr, nil
}

// Canonicalizes search strings including
// - replacing instances of " " with "+"
// - ...
func canonicalizeQuery(query string) string {
	return strings.ReplaceAll(query, " ", "+")
}

type SearchResponse struct {
	Kind       string       `json:"kind"`
	URL        URL          `json:"url"`
	Queries    Queries      `json:"queries"`
	Context    Context      `json:"context"`
	SearchInfo SearchInfo   `json:"searchInformation"`
	Items      []SearchItem `json:"items"`
}

// URL represents the URL information in the response
type URL struct {
	Type     string `json:"type"`
	Template string `json:"template"`
}

// Queries represents the query information including request and nextPage
type Queries struct {
	Request  []Query `json:"request"`
	NextPage []Query `json:"nextPage"`
}

// Query represents a single query with its parameters
type Query struct {
	Title          string `json:"title"`
	TotalResults   string `json:"totalResults"`
	SearchTerms    string `json:"searchTerms"`
	Count          int    `json:"count"`
	StartIndex     int    `json:"startIndex"`
	InputEncoding  string `json:"inputEncoding"`
	OutputEncoding string `json:"outputEncoding"`
	Safe           string `json:"safe"`
	CX             string `json:"cx"`
}

// Context represents the context information
type Context struct {
	Title string `json:"title"`
}

// SearchInfo represents the search information
type SearchInfo struct {
	SearchTime            float64 `json:"searchTime"`
	FormattedSearchTime   string  `json:"formattedSearchTime"`
	TotalResults          string  `json:"totalResults"`
	FormattedTotalResults string  `json:"formattedTotalResults"`
}

// SearchItem represents a single search result item
type SearchItem struct {
	Kind             string                   `json:"kind"`
	Title            string                   `json:"title"`
	HTMLTitle        string                   `json:"htmlTitle"`
	Link             string                   `json:"link"`
	DisplayLink      string                   `json:"displayLink"`
	Snippet          string                   `json:"snippet"`
	HTMLSnippet      string                   `json:"htmlSnippet"`
	FormattedURL     string                   `json:"formattedUrl"`
	HTMLFormattedURL string                   `json:"htmlFormattedUrl"`
	Pagemap          map[string][]PagemapItem `json:"pagemap"`
	Mime             string                   `json:"mime,omitempty"`
	FileFormat       string                   `json:"fileFormat,omitempty"`
}

// PagemapItem represents an item in the pagemap
type PagemapItem struct {
	// Common fields
	Src    string `json:"src,omitempty"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`

	// Metatags fields
	OgImage            string `json:"og:image,omitempty"`
	OgType             string `json:"og:type,omitempty"`
	OgTitle            string `json:"og:title,omitempty"`
	OgDescription      string `json:"og:description,omitempty"`
	OgURL              string `json:"og:url,omitempty"`
	TwitterCard        string `json:"twitter:card,omitempty"`
	TwitterTitle       string `json:"twitter:title,omitempty"`
	TwitterDescription string `json:"twitter:description,omitempty"`
	TwitterImage       string `json:"twitter:image,omitempty"`
	Viewport           string `json:"viewport,omitempty"`
}
