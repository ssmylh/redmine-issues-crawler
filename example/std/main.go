package main

import (
	"fmt"
	"github.com/ssmylh/redmine-issues-crawler"
	"time"
)

type std struct{}

func (so *std) Output(issue crawler.Issue) error {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	fmt.Println(issue.Id, issue.Tracker.Name, issue.Status.Name, issue.Priority.Name,
		issue.Subject, issue.UpdatedOn.In(jst))
	return nil
}

func (so *std) Select(issues crawler.Issue) bool {
	ss := []string{"新規", "終了"}
	for _, s := range ss {
		if s == issues.Status.Name {
			return true
		}
	}
	return false
}

func main() {
	std := &std{}
	c := crawler.NewCrawler(
		"redmine project url",
		10,
		20,
		std,
		std,
	)
	err := c.Crawl(time.Now().Add(-30 * time.Minute))
	fmt.Println(err)
}
