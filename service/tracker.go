package main

import (
	ga "github.com/OzqurYalcin/google-analytics/src"
	"strconv"
)

type Tracker struct {
	api        *ga.API
	trackingId string
}

func newTracker(trackingId string, appName string) *Tracker {
	api := ga.API{}
	api.UserAgent = appName
	api.ContentType = "application/x-www-form-urlencoded"

	return &Tracker{api: &api, trackingId: trackingId}
}

func (t *Tracker) trackRequest(userId string, numberOfCommits int) {
	if len(t.trackingId) == 0 {
		return
	}

	client := ga.Client{}
	client.ProtocolVersion = "1"
	client.UserID = userId
	client.TrackingID = t.trackingId
	client.HitType = "event"
	client.EventCategory = "pull-request"
	client.EventAction = "request"
	client.EventLabel = userId
	client.EventValue = strconv.Itoa(numberOfCommits)

	t.api.Send(&client)
}
