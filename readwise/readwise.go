package readwise

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type ListResult struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []ReadwiseItem `json:"results"`
}

type ReadwiseItem struct {
	Id                int    `json:"id"`
	Title             string `json:"title"`
	Author            string `json:"author"`
	Category          string `json:"category"`
	Source            string `json:"source"`
	NumHighlights     int    `json:"num_highlights"`
	LastHighlightedAt string `json:"last_highlight_at"`
	Updated           string `json:"updated"`
	CoverImageUrl     string `json:"cover_image_url"`
	HighlightsUrl     string `json:"highlights_url"`
	SourceUrl         string `json:"source_url"`
	Asin              string `json:"asin"`
	Tags              []Tag  `json:"tags"`
}

type Tag struct {
	Id       int64  `json:"id"`
	UserBook int64  `json:"user_book"`
	Name     string `json:"name"`
}

type ExportResult struct {
	Count          int             `json:"count"`
	NextPageCursor int             `json:"nextPageCursor"`
	Results        []HighlightItem `json:"results"`
}

type HighlightItem struct {
	UserBookId    int         `json:"user_book_id"`
	Title         string      `json:"title"`
	Author        string      `json:"author"`
	ReadableTitle string      `json:"readable_title"`
	Source        string      `json:"source"`
	CoverImageUrl string      `json:"cover_image_url"`
	UniqueUrl     string      `json:"unique_url"`
	BookTags      []Tag       `json:"book_tags"`
	Category      string      `json:"category"`
	ReadwiseUrl   string      `json:"readwise_url"`
	SourceUrl     string      `json:"source_url"`
	Asin          string      `json:"asin"`
	Highlights    []Highlight `json:"highlights"`
}

type Highlight struct {
	Id            int    `json:"id"`
	Text          string `json:"text"`
	Location      int    `json:"location"`
	LocationType  string `json:"location_type"`
	Note          string `json:"note"`
	Color         string `json:"color"`
	HighlightedAt string `json:"highlighted_at"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	ExternalId    int    `json:"external_id"`
	EndLocation   int    `json:"end_location"`
	Url           string `json:"url"`
	BookId        int    `json:"book_id"`
	Tags          []Tag  `json:"tags"`
	IsFavorite    bool   `json:"is_favorite"`
	IsDiscard     bool   `json:"is_discard"`
	ReadwiseUrl   string `json:"readwise_url"`
}

// Lists books based on specified filterting criteria.
func List(token string, category string) error {
	client := resty.New()
	url := "https://readwise.io/api/v2/books/"
	ret := []ReadwiseItem{}

	for {
		result := ListResult{}

		_, err := client.R().
			SetHeader("Authorization", fmt.Sprintf("Token %s", token)).
			SetResult(&result).
			SetQueryParams(map[string]string{
				"category": category,
			}).
			Get(url)

		if err != nil {
			return err
		}

		ret = append(ret, result.Results...)

		if result.Next != "" {
			url = result.Next
		} else {
			break
		}
	}

	for _, item := range ret {
		j, _ := json.MarshalIndent(item, "", "\t")
		fmt.Println(string(j))
	}

	return nil
}

// Lists highlights from readwise using the Readwise Export API. Highlights can
// be limited by passing in additional filtering conditions such as 'updatedAfter'
// which limits the highlights fetched from past 'x' days, and 'ids' which
// fetches highlights from only specified book ids.
func ListHighlights(token string, updatedAfter int, ids []string) error {
	client := resty.New()
	url := "https://readwise.io/api/v2/export/"
	ret := []HighlightItem{}
	params := map[string]string{}

	if updatedAfter > 0 {
		params["updatedAfter"] = string(time.Now().AddDate(0, 0, -updatedAfter).Format(time.RFC3339))
	}

	if len(ids) > 0 {
		params["ids"] = strings.Join(ids[:], ",")
	}

	for {
		result := ExportResult{}

		_, err := client.R().
			SetHeader("Authorization", fmt.Sprintf("Token %s", token)).
			SetResult(&result).
			SetQueryParams(params).
			Get(url)

		// fmt.Printf("%+v\n", string(req.Body()))

		if err != nil {
			return err
		}

		ret = append(ret, result.Results...)

		if result.NextPageCursor > 0 {
			params["pageCursor"] = strconv.Itoa(result.NextPageCursor)
		} else {
			break
		}
	}

	fmt.Println("%%tana%%")
	formatHighlights(ret)

	return nil
}

// Formats highlights based on Tana Paste format.
func formatHighlights(data []HighlightItem) {
	for _, item := range data {
		fmt.Printf("- %s #readwise\n", item.Title)

		if item.SourceUrl != "" && strings.HasPrefix(item.SourceUrl, "https://") {
			fmt.Printf("  - Source URL:: %s\n", item.SourceUrl)
		}

		fmt.Printf("  - Type:: %s\n", getCategory(item.Category))
		fmt.Printf("  - Author:: %s\n", item.Author)
		fmt.Printf("  - Readwise URL:: %s\n", item.ReadwiseUrl)

		if len(item.Highlights) > 0 {
			fmt.Printf("  - Highlights\n")
			sort.Slice(item.Highlights, func(i, j int) bool {
				return item.Highlights[i].Location < item.Highlights[j].Location
			})

			for _, highlight := range item.Highlights {
				processHighlight(highlight)
			}
		}
	}
}

// Process a individual highlight
func processHighlight(highlight Highlight) {
	headers := []string{".h1", ".h2", ".h3"}
	lines := strings.Split(highlight.Text, "\n")
	re := regexp.MustCompile(`•\s+`)

	for _, line := range lines {
		cleanedLine := re.ReplaceAllString(strings.TrimSpace(line), "")

		if len(cleanedLine) > 0 {
			if highlight.Note != "" && contains(highlight.Note, headers) {
				fmt.Printf("    - %s %s\n", highlight.Note, cleanedLine)
			} else {
				fmt.Printf("    - %s\n", cleanedLine)
			}
		}
	}
}

func contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getCategory(cat string) string {
	switch cat {
	case "books":
		return "Book"
	case "articles":
		return "Article"
	case "tweets":
		return "Tweet"
	case "podcasts":
		return "Podcast"
	default:
		return ""
	}
}
