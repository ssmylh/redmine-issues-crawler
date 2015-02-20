package main

import (
	"encoding/json"
	"fmt"
	"github.com/ssmylh/redmine-issues-crawler"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type hipchatNotifier struct {
	Settings *settings
}

func (h *hipchatNotifier) Output(issue crawler.Issue) error {
	message, err := h.createMessage(issue)
	if err != nil {
		return err
	}

	req := &hipchat.NotificationRequest{
		Message:       message,
		Color:         h.Settings.NotificationColor,
		MessageFormat: "text",
	}

	c := hipchat.NewClient(h.Settings.RoomNotificationToken)
	resp, err := c.Room.Notification(h.Settings.RoomId, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (h *hipchatNotifier) createMessage(issue crawler.Issue) (string, error) {
	issueUrl := h.Settings.RedmineProjectUrl
	if !strings.HasSuffix(issueUrl, "/") {
		issueUrl += "/"
	}
	issueUrl += strconv.Itoa(issue.Id)

	return fmt.Sprintf("%s #%d (%s) : [%s / author : %s] assigned to : %s - %s",
		issue.Tracker.Name, issue.Id,
		issue.Status.Name, issue.Subject,
		issue.Author.Name, issue.AssignedTo.Name,
		issueUrl), nil
}

func (h *hipchatNotifier) Select(issues crawler.Issue) bool {
	return true
}

type settings struct {
	RoomId                string
	RoomNotificationToken string
	NotificationColor     string
	RedmineProjectUrl     string
	CrawlInterval         int
	FetchLimit            int
}

func readSettings() (*settings, error) {
	name := "settings.json"
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	settings := &settings{}
	err = json.Unmarshal(data, settings)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func main() {
	settings, err := readSettings()
	if err != nil {
		fmt.Println("could not read settings.")
		fmt.Println(err)
		os.Exit(1)
	}

	h := &hipchatNotifier{
		Settings: settings,
	}
	c := crawler.NewCrawler(
		settings.RedmineProjectUrl,
		settings.CrawlInterval,
		settings.FetchLimit,
		h,
	)
	err = c.Crawl(time.Now())
	if err != nil {
		fmt.Println("error occurred.")
		fmt.Println(err)
		os.Exit(1)
	}
}
