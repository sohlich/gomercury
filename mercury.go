package gomercury

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const API_URL = "https://mercury.postlight.com/parser?url=%s"

type MercuryConfig struct {
	ApiKey string
	Log    *log.Logger
}

type MercuryClient struct {
	*MercuryConfig
	client *http.Client
}

type MercuryDocument struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	DatePublished time.Time `json:"date_published"`
	LeadImageURL  string    `json:"lead_image_url"`
	Dek           string    `json:"dek"`
	NextPageURL   string    `json:"next_page_url"`
	URL           string    `json:"url"`
	Domain        string    `json:"domain"`
	Excerpt       string    `json:"excerpt"`
	WordCount     int       `json:"word_count"`
	Direction     string    `json:"direction"`
	TotalPages    int       `json:"total_pages"`
	RenderedPages int       `json:"rendered_pages"`
}

func New(c *MercuryConfig) *MercuryClient {
	if c.Log == nil {
		c.Log = log.New(ioutil.Discard, "mercury:", 0)
	}
	return &MercuryClient{
		c,
		http.DefaultClient,
	}
}

func (c *MercuryClient) Parse(URL string) (*MercuryDocument, error) {
	fURL := formatUrl(URL)
	c.Log.Printf("Formated url: %s", fURL)
	req, _ := http.NewRequest("GET", fURL, nil)
	req.Header.Set("x-api-key", c.ApiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get response from service.")
	}
	c.Log.Println(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Cannot get response from service. Status code: %d", resp.StatusCode))
	}
	// defer resp.Body.Close()
	return decodeToDocument(resp.Body)
}

func formatUrl(URL string) string {
	return fmt.Sprintf(API_URL, url.QueryEscape(URL))
}

func decodeToDocument(r io.Reader) (*MercuryDocument, error) {
	dec := json.NewDecoder(r)
	doc := &MercuryDocument{}
	err := dec.Decode(doc)
	return doc, err
}
