package crawler

import (
	"strings"
	"testing"
)

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
