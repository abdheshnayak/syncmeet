package domain

import (
	"context"
	"fmt"

	"github.com/kloudlite/api/pkg/repos"

	"github.com/abdheshnayak/syncmeet/syncmeet/entities"
)

func validateSlots(slots []entities.TimeSlot) error {
	for _, ts := range slots {
		if _, err := ts.Validate(); err != nil {
			return err
		}

	}

	return nil
}

type impl struct {
	eventsRepo        repos.DbRepo[*entities.Event]
	usersRepo         repos.DbRepo[*entities.User]
	participationRepo repos.DbRepo[*entities.Participation]
}

// UpdateEventSlots implements Domain.
func (i *impl) UpdateEventSlots(ctx context.Context, eventId repos.ID, slots []entities.TimeSlot) (*entities.Event, error) {

	if err := validateSlots(slots); err != nil {
		return nil, err
	}

	i.eventsRepo.PatchById(ctx, eventId, repos.Document{
		"slots": slots,
	})

	return i.eventsRepo.FindById(ctx, eventId)
}

// GetRecommendedSlots implements Domain.
func (i *impl) GetRecommendedSlots(ctx context.Context, eventId repos.ID) ([]entities.TimeSlot, error) {
	participations, err := i.participationRepo.Find(ctx, repos.Query{
		Filter: repos.Filter{
			"event_id": eventId,
		},
	})
	if err != nil {
		return nil, err
	}

	slotsCount := make(map[string]int)
	slotsMap := make(map[string]entities.TimeSlot)

	for _, p := range participations {
		for _, ts := range p.Slots {
			date, startTime, endTime, err := ts.Parse()
			if err != nil {
				fmt.Println(ts)
				return nil, err
			}

			key := fmt.Sprintf("%s-%s-%s", date, startTime, endTime)
			slotsCount[key] += 1
			slotsMap[key] = ts
		}
	}

	maxSlots := 0
	maxSlotsKey := ""
	for k, v := range slotsCount {
		if v > maxSlots {
			maxSlots = v
			maxSlotsKey = k
		}
	}

	return []entities.TimeSlot{slotsMap[maxSlotsKey]}, nil
}

// CreateParticipation implements Domain.
func (i *impl) CreateParticipation(ctx context.Context, participation *entities.Participation) (*entities.Participation, error) {
	if err := validateSlots(participation.Slots); err != nil {
		return nil, err
	}

	return i.participationRepo.Create(ctx, participation)
}

// CreateUser implements Domain.
func (i *impl) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	return i.usersRepo.Create(ctx, user)
}

// GetParticipation implements Domain.
func (i *impl) GetParticipation(ctx context.Context, id repos.ID) (*entities.Participation, error) {
	return i.participationRepo.FindOne(ctx, repos.Filter{"id": id})
}

// GetUser implements Domain.
func (i *impl) GetUser(ctx context.Context, id repos.ID) (*entities.User, error) {
	return i.usersRepo.FindOne(ctx, repos.Filter{"id": id})
}

// ListParticipations implements Domain.
func (i *impl) ListParticipations(ctx context.Context, eventId repos.ID) ([]*entities.Participation, error) {
	return i.participationRepo.Find(ctx, repos.Query{
		Filter: repos.Filter{"event_id": eventId},
	})
}

// ListUsers implements Domain.
func (i *impl) ListUsers(ctx context.Context) ([]*entities.User, error) {
	return i.usersRepo.Find(ctx, repos.Query{})
}

// UpdateParticipation implements Domain.
func (i *impl) UpdateParticipation(ctx context.Context, participation *entities.Participation) (*entities.Participation, error) {

	if err := validateSlots(participation.Slots); err != nil {
		return nil, err
	}

	return i.participationRepo.PatchById(ctx, participation.Id, repos.Document{
		fmt.Sprintf("slots.%s", participation.UserId): participation.Slots,
	})

}

// UpdateSlots implements Domain.
func (i *impl) UpdateSlots(ctx context.Context, eventId repos.ID, userId repos.ID, slots []entities.TimeSlot) (*entities.User, error) {

	if err := validateSlots(slots); err != nil {
		return nil, err
	}

	return i.usersRepo.PatchById(ctx, eventId, repos.Document{
		fmt.Sprintf("default_time_slots.%s", userId): slots,
	})
}

// ListEvents implements Domain.
func (i *impl) ListEvents(ctx context.Context) ([]*entities.Event, error) {
	q := repos.Query{}

	return i.eventsRepo.Find(ctx, q)
}

func (i *impl) CreateEvent(ctx context.Context, event *entities.Event) (*entities.Event, error) {
	return i.eventsRepo.Create(ctx, event)
}

func (i *impl) GetEvent(ctx context.Context, id repos.ID) (*entities.Event, error) {
	return i.eventsRepo.FindOne(ctx, repos.Filter{"id": id})
}
