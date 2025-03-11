package domain

import (
	"context"

	"github.com/abdheshnayak/syncmeet/syncmeet/entities"
	"github.com/kloudlite/api/pkg/repos"

	"go.uber.org/fx"
)

type Domain interface {
	CreateEvent(ctx context.Context, event *entities.Event) (*entities.Event, error)
	GetEvent(ctx context.Context, id repos.ID) (*entities.Event, error)
	ListEvents(ctx context.Context) ([]*entities.Event, error)
	UpdateEventSlots(ctx context.Context, eventId repos.ID, slots []entities.TimeSlot) (*entities.Event, error)

	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUser(ctx context.Context, id repos.ID) (*entities.User, error)
	UpdateSlots(ctx context.Context, eventId repos.ID, userId repos.ID, slots []entities.TimeSlot) (*entities.User, error)
	ListUsers(ctx context.Context) ([]*entities.User, error)

	CreateParticipation(ctx context.Context, participation *entities.Participation) (*entities.Participation, error)
	GetParticipation(ctx context.Context, id repos.ID) (*entities.Participation, error)
	UpdateParticipation(ctx context.Context, participation *entities.Participation) (*entities.Participation, error) //nolint:lll
	ListParticipations(ctx context.Context, eventId repos.ID) ([]*entities.Participation, error)

	GetRecommendedSlots(ctx context.Context, eventId repos.ID) ([]entities.TimeSlot, error)
}

var Module = fx.Module("domain",
	fx.Provide(func(
		eventsRepo repos.DbRepo[*entities.Event],
		usersRepo repos.DbRepo[*entities.User],
		participationRepo repos.DbRepo[*entities.Participation],
	) Domain {
		return &impl{
			eventsRepo:        eventsRepo,
			usersRepo:         usersRepo,
			participationRepo: participationRepo,
		}
	}),
)
