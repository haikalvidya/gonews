package news

import (
	"net/http"
	"net/url"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type ArticleStruct struct {
	Source struct {
		ID		interface{}	`json:"id"`
		Name	string		`json:"name"`
	} `json:"source"`
	Author		string		`json:"author"`
	Title		string		`json:"title"`
	Description	string		`json:"description"`
	URL			string		`json:"url"`
	URLToImage	string		`json:"urlToImage"`
	PublishedAt	time.Time	`json:"publishedAt"`
	Content		string		`json:"content"`
}

type ResultFromJson struct {
	Status			string 			`json:"status"`
	TotalResults	int				`json:"totalResults"`
	Articles		[]ArticleStruct	`json:"articles"`
}

type Client struct {
	http		*http.Client
	key			string
	PageSize	int
}

func (c *Client) FetchAll(query, page string) (*ResultFromJson, error) {
	endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en", url.QueryEscape(query), c.PageSize, page, c.key)
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &ResultFromJson{}
	return res, json.Unmarshal(body,res)
}

func NewClient(httpClient *http.Client, key string, pagesize int) *Client {
	if pagesize > 100 {
		pagesize = 100
	}

	return &Client{httpClient, key, pagesize}
}

func (a *ArticleStruct) FormatPublishedDate() string {
	year, month, day := a.PublishedAt.Date()
	return fmt.Sprintf("%v %d, %d", month, day, year)
}