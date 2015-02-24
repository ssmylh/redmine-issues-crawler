package main

import (
	"fmt"
	"github.com/ssmylh/redmine-issues-crawler"
	"time"
)

type std struct{}

func (so *std) Output(issue *crawler.Issue) error {
	fmt.Println(issue.Id, issue.Tracker.Name, issue.Status.Name, issue.Priority.Name,
		issue.Subject, issue.UpdatedOn)
	return nil
}

func (so *std) Select(issue *crawler.Issue) bool {
	ss := []string{"新規", "終了"}
	for _, s := range ss {
		if s == issue.Status.Name {
			return true
		}
	}
	return false
}

func main() {
	std := &std{}
	c := crawler.NewCrawler(
		"your redmine's endpoint(redmine's home url)",
		10, // interval
		20, // fetch limit
		std,
	)
	c.Selector = std
	err := c.Crawl(time.Now())
	fmt.Println(err)
}
