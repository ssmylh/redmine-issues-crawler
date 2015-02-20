package crawler

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestBuildFetchUrl(t *testing.T) {
	url := "https://example.com/redmine/projects/dummy"
	limit := 5
	lastUpdate := time.Date(2015, 2, 20, 20, 30, 30, 0, time.UTC)
	expected := url + "/" +
		"issues.json?sort=updated_on:desc&id:desc" +
		"&limit=" + strconv.Itoa(limit) +
		"&updated_on=%3E%3D" + lastUpdate.Add(1*time.Second).Format(time.RFC3339)

	c := NewCrawler(url, 10, limit, nil)
	actual := c.BuildFetchUrl(lastUpdate)

	if expected != actual {
		t.Errorf("fetch url is invalid. expected : %s, but actual : %s ", expected, actual)
	}
}

func TestFilterEmptyIssues(t *testing.T) {
	issues := make([]Issue, 0)
	filtered := Filter(issues, func(issue Issue) bool {
		return true
	})
	if len(filtered) != 0 {
		t.Errorf("filter empty issues should return empty issues, but length is %d", len(filtered))
	}
}

func TestFilterIssues(t *testing.T) {
	issue1 := Issue{
		Id:      1,
		Subject: "Foo Subject 1",
	}
	issue2 := Issue{
		Id:      2,
		Subject: "Boo Subject 2",
	}
	issue3 := Issue{
		Id:      3,
		Subject: "Foo Subject 3",
	}
	issues := []Issue{issue1, issue2, issue3}
	filtered := Filter(issues, func(issue Issue) bool {
		return strings.HasPrefix(issue.Subject, "Foo")
	})

	if len(filtered) != 2 {
		t.Errorf("filter issues should return %d elements, but %d elements", 2, len(filtered))
	}

	if issue1 != filtered[0] {
		t.Errorf("1st element of filtered issues is not issue1")
	}

	if issue3 != filtered[1] {
		t.Errorf("1st element of filtered issues is not issue3")
	}
}
