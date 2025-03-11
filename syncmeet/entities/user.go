package entities

import "github.com/kloudlite/api/pkg/repos"

type User struct {
	repos.BaseEntity `json:",inline"`
	Email            string     `json:"email"`
	Name             string     `json:"name"`
	DefaultTimeSlots []TimeSlot `json:"default_time_slots"`
}

var UserIndexes = []repos.IndexField{
	{
		Field: []repos.IndexKey{
			{Key: "id", Value: repos.IndexAsc},
		},
		Unique: true,
	},
	{
		Field: []repos.IndexKey{
			{Key: "email", Value: repos.IndexAsc},
		},
		Unique: true,
	},
}
