package crawler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Issue struct {
	Id          int
	Project     Term
	Tracker     Term
	Status      Term
	Priority    Term
	Author      Term
	Category    Term
	Subject     string
	Description string
	AssignedTo  Term      `json:"assigned_to"`
	CreatedOn   time.Time `json:"created_on"`
	UpdatedOn   time.Time `json:"updated_on"`
}

type issuesResponse struct {
	Issues     []Issue
	TotalCount int `json:"total_count"`
	Offset     int
	Limit      int
}

type Term struct {
	Id   int
	Name string
}

type Outputter interface {
	Output(issue Issue) error
}

type Selector interface {
	Select(issue Issue) bool
}

type Crawler struct {
	Url       string
	Interval  int
	Limit     int
	Outputter Outputter
	Selector  Selector
}

// NewCrawler returs a new Crawler.
// The url is Redmines project url.
// The interval is interval of crawling.
// The limit is limit on the number of per fetch.
// The outputter is how this Crawler outputs fetched issues.
func NewCrawler(url string, interval int, limit int, outputter Outputter) *Crawler {
	if interval < 10 {
		interval = 10
	}
	c := &Crawler{
		Interval:  interval,
		Url:       url,
		Limit:     limit,
		Outputter: outputter,
	}
	return c
}

func (c *Crawler) Crawl(startTime time.Time) error {
	url := c.Url
	if !strings.HasSuffix(c.Url, "/") {
		url += "/"
	}
	url += "issues.json"
	query := "?sort=updated_on:desc&id:desc" + "&limit=" + strconv.Itoa(c.Limit)

	lastUpdate := startTime
	for _ = range time.Tick(time.Duration(c.Interval) * time.Second) {
		fetchUrl := url + query + "&updated_on=%3E%3D" + lastUpdate.Add(1*time.Second).UTC().Format(time.RFC3339)
		issuesResp, err := Fetch(fetchUrl)
		if err != nil {
			return err
		}

		issues := issuesResp.Issues
		if len(issues) == 0 {
			continue
		}
		lastUpdate = issues[0].UpdatedOn

		issues = Filter(issues, c.Selector.Select)
		for _, issue := range issues {
			err = c.Outputter.Output(issue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Fetch(url string) (*issuesResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	issuesResp := &issuesResponse{}
	err = dec.Decode(issuesResp)
	if err != nil {
		return nil, err
	}
	return issuesResp, nil
}

func Filter(issues []Issue, predicate func(Issue) bool) []Issue {
	if len(issues) == 0 {
		return issues
	}

	capacity := (len(issues) + 1) / 2
	filtered := make([]Issue, 0, capacity)
	for _, issue := range issues {
		if predicate(issue) {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}
