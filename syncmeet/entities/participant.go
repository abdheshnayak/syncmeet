package entities

import "github.com/kloudlite/api/pkg/repos"

type Participation struct {
	repos.BaseEntity `json:",inline"`
	EventId          string     `json:"event_id"`
	UserId           string     `json:"user_id"`
	Slots            []TimeSlot `json:"slots"`
}

var ParticipationIndexes = []repos.IndexField{
	{
		Field: []repos.IndexKey{
			{Key: "id", Value: repos.IndexAsc},
		},
		Unique: true,
	},
	{
		Field: []repos.IndexKey{
			{Key: "user_id", Value: repos.IndexAsc},
			{Key: "event_id", Value: repos.IndexAsc},
		},
		Unique: true,
	},
}
