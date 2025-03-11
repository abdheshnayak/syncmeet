package entities

import (
	"time"

	"github.com/kloudlite/api/pkg/repos"
)

type Email string

type Event struct {
	repos.BaseEntity `json:",inline"`
	Title            string     `bson:"title" json:"title"`
	Slots            []TimeSlot `bson:"slots" json:"slots"`
	Duration         int        `bson:"duration" json:"duration"`
}

type TimeSlot string

var EventIndexes = []repos.IndexField{
	{
		Field: []repos.IndexKey{
			{Key: "id", Value: repos.IndexAsc},
		},
		Unique: true,
	},
}

type TimeZone string

type TimeSlotParsed struct {
	StartTime time.Time `bson:"start_time" json:"start_time"`
	EndTime   time.Time `bson:"end_time" json:"end_time"`
	TimeZone  TimeZone  `bson:"time_zone" json:"time_zone"`
}

func (tRange TimeSlot) Parse() (date, startTime, endTime string, err error) {
	// Define the layout to parse the input with timezone
	layout := "2 Jan 2006, 3 - 4PM MST"

	// Parse the input time range
	t, err := time.Parse(layout, string(tRange))
	if err != nil {
		return "", "", "", err
	}

	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())

	// Extract date
	date = t.Format("2006-01-02")

	// Convert to UTC
	loc, _ := time.LoadLocation("UTC")
	startUTC := t.In(loc)
	endUTC := startUTC.Add(2 * time.Hour) // Since 2-4PM means 2 hours duration

	// Format times
	startTime = startUTC.Format("15:04:05 UTC")
	endTime = endUTC.Format("15:04:05 UTC")

	return date, startTime, endTime, nil
}

func (tRange TimeSlot) Validate() (bool, error) {
	// Define the layout to parse the input with timezone
	layout := "2 Jan 2006, 3 - 4PM MST"

	// Parse the input time range
	_, err := time.Parse(layout, string(tRange))
	if err != nil {
		return false, err
	}

	return true, nil
}
