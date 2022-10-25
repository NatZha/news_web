package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"net/http"
	"time"
)

type Article struct {
	Source struct {
		ID   interface{} `json:"id"`
		Name string      `json:"name"`
	} `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

// saves articles into its own struct
type Results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

// represent the client for working with the News API
type Client struct {
	http		*http.Client
	key 		string 
	PageSize 	int 
}


// FetchEverything appends query and string to request URL as will as API key
// and page size
func (c *Client) FetchEverything(query, page string) (*Results, error) {
	endpoint := fmt.Sprintf(
		"https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en", 
		url.QueryEscape(query), c.PageSize, page, c.key)

	resp, err := c.http.Get(endpoint) 
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() 

	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil {
		return nil, err 
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf((string(body)))
	}

	res := &Results{} 
	return res, json.Unmarshal(body, res)


}


/* 
 * Create and return a new client instance for each request  tothe News API
 * :httpClient : *http.Client pointer to the HTTP client used to make requests
 * :key :string  holds the API key
 * :pageSize :int holds number fo results to return per page
 */
func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{httpClient, key, pageSize}
}