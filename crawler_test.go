package crawler

import (
	"strings"
	"testing"
	"time"
)

func TestBuildFetchUrl(t *testing.T) {
	url := "https://example.com/redmine"
	key := "aaabbbcccxxxyyyzzz"
	expected := url + "/" +
		"issues.json" +
		"?key=" + key +
		"&limit=" + "5" +
		"&sort=updated_on:desc,id:desc" +
		"&status_id=*"

	c := NewCrawler(url, key, 10, 5, nil)
	lastUpdate := time.Date(2015, 2, 20, 20, 30, 30, 0, time.UTC)
	actual := c.BuildFetchUrl(lastUpdate)

	if expected != actual {
		t.Errorf("fetch url is invalid. expected : %s, but actual : %s ", expected, actual)
	}
}

type testOutputter struct {
	Done []Issue
}

func (to *testOutputter) Output(issue *Issue) error {
	to.Done = append(to.Done, *issue)
	return nil
}

func TestOutput(t *testing.T) {
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

	to := &testOutputter{
		Done: make([]Issue, 0, 3),
	}
	c := NewCrawler("", "", 10, 20, to)

	c.Output(issues)

	if to.Done[0] != issue3 {
		t.Errorf("1st output element should be ID3, but ID%d", to.Done[0].Id)
	}
	if to.Done[1] != issue2 {
		t.Errorf("2nd output element should be ID2, but ID%d", to.Done[1].Id)
	}
	if to.Done[2] != issue1 {
		t.Errorf("3rd output element should be ID1, but ID%d", to.Done[2].Id)
	}
}

func TestFilterEmptyIssues(t *testing.T) {
	issues := make([]Issue, 0)
	filtered := Filter(issues, func(issue *Issue) bool {
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
	filtered := Filter(issues, func(issue *Issue) bool {
		return strings.HasPrefix(issue.Subject, "Foo")
	})

	if len(filtered) != 2 {
		t.Errorf("filter issues should return %d elements, but %d elements", 2, len(filtered))
	}

	if filtered[0] != issue1 {
		t.Errorf("1st element of filtered issues is not issue1")
	}

	if filtered[1] != issue3 {
		t.Errorf("2nd element of filtered issues is not issue3")
	}
}

func TestToUTCTime_RFC3339(t *testing.T) {
	s := "2015-02-24T15:58:38Z"
	_, err := ToUTCTime(s)
	if err != nil {
		t.Errorf("could not convert %s into to UTC Time", s)
	}
}

func TestToUTCTime_JST(t *testing.T) {
	s := "2015/02/25 01:02:03 +0900"
	_, err := ToUTCTime(s)
	if err != nil {
		t.Errorf("could not convert %s into to UTC Time", s)
	}
}

func TestFilterWithUpdatedOnAfter(t *testing.T) {
	lastUpdate := time.Date(2015, 2, 28, 8, 30, 30, 0, time.UTC)
	issue1 := Issue{
		Id:        1,
		Subject:   "Foo Subject 1",
		UpdatedOn: lastUpdate.Add(2 * time.Minute).Format(time.RFC3339),
	}
	issue2 := Issue{
		Id:        2,
		Subject:   "Boo Subject 2",
		UpdatedOn: lastUpdate.Add(1 * time.Minute).Format(time.RFC3339),
	}
	issue3 := Issue{
		Id:        3,
		Subject:   "Foo Subject 3",
		UpdatedOn: lastUpdate.Format(time.RFC3339),
	}
	issues := []Issue{issue1, issue2, issue3}
	c := NewCrawler("", "", 10, 20, nil)

	filtered := c.filterWithUpdatedOnAfter(issues, lastUpdate)

	if len(filtered) != 2 {
		t.Errorf("filtered issues should return %d elements, but %d elements", 2, len(filtered))
	}

	if filtered[0] != issue1 {
		t.Errorf("1st element of filtered issues is not issue1")
	}

	if filtered[1] != issue2 {
		t.Errorf("2nd element of filtered issues is not issue2")
	}
}
