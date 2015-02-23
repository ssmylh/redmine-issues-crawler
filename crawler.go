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

type IssuesUrl struct {
	Endpoint string
	Offset   int
	Limit    int
	Sort     string
	StatusId string
}

// String builds url for issues and returns it.
// In addtion to IssuesUrl properties, it appends updated_on(UTC, RFC3339), too.
// More precisely, it set ">=" + add 1 second to lastUpdate.
// The reason why adding 1 second is that it can set ">=" to updated_on, but can not set ">".
func (iu *IssuesUrl) String(lastUpdate time.Time) string {
	url := iu.Endpoint
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "issues.json"
	url += "?updated_on=%3E%3D" + lastUpdate.Add(1*time.Second).UTC().Format(time.RFC3339)

	if iu.Offset > 0 {
		url += "offset=" + strconv.Itoa(iu.Offset)
	}
	if iu.Limit > 0 {
		url += "&limit=" + strconv.Itoa(iu.Limit)
	}
	if iu.Sort != "" {
		url += "&sort=" + iu.Sort
	}
	if iu.StatusId != "" {
		url += "&status_id=" + iu.StatusId
	}
	return url
}

type Outputter interface {
	Output(issue *Issue) error
}

type Selector interface {
	Select(issue *Issue) bool
}

type Crawler struct {
	Interval  int
	Url       *IssuesUrl
	Outputter Outputter
	Selector  Selector
}

// NewCrawler returns a new Crawler.
// The endpoint is Redmine's endpoint(Redmine's home URL).
// The interval is interval of crawling.
// The limit is limit on the number of per fetch.
// The outputter is how this Crawler outputs fetched issues.
func NewCrawler(endpoint string, interval int, limit int, outputter Outputter) *Crawler {
	if interval < 10 {
		interval = 10
	}

	url := &IssuesUrl{
		Endpoint: endpoint,
		Limit:    limit,
		Sort:     "updated_on:desc,id:desc",
		StatusId: "*",
	}
	c := &Crawler{
		Url:       url,
		Interval:  interval,
		Outputter: outputter,
	}
	return c
}

func (c *Crawler) Crawl(startTime time.Time) error {
	lastUpdate := startTime
	for _ = range time.Tick(time.Duration(c.Interval) * time.Second) {
		fetchUrl := c.BuildFetchUrl(lastUpdate)
		issuesResp, err := Fetch(fetchUrl)
		if err != nil {
			return err
		}

		issues := issuesResp.Issues
		if len(issues) == 0 {
			continue
		}
		lastUpdate = issues[0].UpdatedOn

		if c.Selector != nil {
			issues = Filter(issues, c.Selector.Select)
		}

		err = c.Output(issues)
		if err != nil {
			return err
		}
	}
	return nil
}

// BuildFetchUrl builds a fetct url from Crawler.Url and returns it.
func (c *Crawler) BuildFetchUrl(lastUpdate time.Time) string {
	return c.Url.String(lastUpdate)
}

// Output outputs fetched issues(sorted by updated_on in desc) in reverse order(updated_on in asc),
// following Outputter implementation.
func (c *Crawler) Output(issues []Issue) error {
	for i := len(issues) - 1; i >= 0; i-- {
		err := c.Outputter.Output(&issues[i])
		if err != nil {
			return err
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

func Filter(issues []Issue, predicate func(*Issue) bool) []Issue {
	if len(issues) == 0 {
		return issues
	}

	capacity := (len(issues) + 1) / 2
	filtered := make([]Issue, 0, capacity)
	for _, issue := range issues {
		if predicate(&issue) {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}
