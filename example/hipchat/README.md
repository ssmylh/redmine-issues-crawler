# Hipchat Example

This example is notify hipchat of fetched issues.

## Dependency

```go
import "github.com/tbruyelle/hipchat-go/hipchat"
```

## Settings

As follows, create `settings.json` into the directory where executable file exists.

```
{
    "RoomId": "Room id or room name",
    "RoomNotificationToken": "Room Notification Tokens",
    "NotificationColor": "Background color for message. Choices are yellow, green, red, purple, gray and random.",
    "RedmineEndpoint": "Your Redmine's endpoint(Redmine's home url)",
    "RedmineAPIKey" : "Your Redmines's API Key. If required, please set."
    "CrawlInterval": "The interval(sec) of crawl",
    "FetchLimit": 20 "The number of per fetch."
}
```