package app

import (
	"net/http"

	"github.com/abdheshnayak/syncmeet/syncmeet/domain"
	"github.com/abdheshnayak/syncmeet/syncmeet/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/kloudlite/api/pkg/repos"
	"go.uber.org/fx"
)

var Module = fx.Module("app",

	repos.NewFxMongoRepo[*entities.Event]("events", "evnt", entities.EventIndexes),
	repos.NewFxMongoRepo[*entities.User]("users", "usr", entities.EventIndexes),
	repos.NewFxMongoRepo[*entities.Participation]("participation", "usr", entities.EventIndexes),

	fx.Invoke(setupRoutes),
)

func setupRoutes(app *fiber.App, d domain.Domain) {
	app.Get("/healthy", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	events := app.Group("/events")
	events.Get(":id/recommended-slot", func(c *fiber.Ctx) error {
		eventId := c.Params("id")
		ts, err := d.GetRecommendedSlots(c.Context(), repos.ID(eventId))

		if err != nil {
			return err
		}

		return c.JSON(ts)
	})

	events.Post("/", func(c *fiber.Ctx) error {
		var e entities.Event

		if err := c.BodyParser(&e); err != nil {
			return err
		}

		resp, err := d.CreateEvent(c.Context(), &e)
		if err != nil {
			return err
		}

		return c.JSON(resp)
	})

	events.Put(":id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		type reqBody struct {
			Slots []entities.TimeSlot
		}
		var rb reqBody

		if err := c.BodyParser(&rb); err != nil {
			return err
		}

		e, err := d.UpdateEventSlots(c.Context(), repos.ID(id), rb.Slots)
		if err != nil {
			return err
		}

		return c.JSON(e)
	})

	events.Get(":id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		e, err := d.GetEvent(c.Context(), repos.ID(id))

		if err != nil {
			return err
		}

		return c.JSON(e)
	})

	events.Get("/", func(c *fiber.Ctx) error {
		e, err := d.ListEvents(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(e)
	})

	users := app.Group("/users")
	users.Get("/", func(c *fiber.Ctx) error {
		u, err := d.ListUsers(c.Context())
		if err != nil {
			return err
		}
		return c.JSON(u)
	})

	users.Post("/", func(c *fiber.Ctx) error {
		var u entities.User
		if err := c.BodyParser(&u); err != nil {
			return err
		}
		resp, err := d.CreateUser(c.Context(), &u)
		if err != nil {
			return err
		}
		return c.JSON(resp)
	})

	users.Put(":id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		type reqBody struct {
			UserId repos.ID            `json:"userId"`
			Slots  []entities.TimeSlot `json:"slots"`
		}
		var rb reqBody

		if err := c.BodyParser(&rb); err != nil {
			return err
		}

		resp, err := d.UpdateSlots(c.Context(), repos.ID(id), rb.UserId, rb.Slots)
		if err != nil {
			return err
		}

		return c.JSON(resp)
	})

	users.Get(":id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		u, err := d.GetUser(c.Context(), repos.ID(id))
		if err != nil {
			return err
		}

		return c.JSON(u)
	})

	participants := app.Group("/participants")

	participants.Get("/:id", func(c *fiber.Ctx) error {
		eventId := c.Params("id")
		p, err := d.ListParticipations(c.Context(), repos.ID(eventId))
		if err != nil {
			return err
		}
		return c.JSON(p)
	})

	participants.Post("/", func(c *fiber.Ctx) error {
		var p entities.Participation
		if err := c.BodyParser(&p); err != nil {
			return err
		}

		resp, err := d.CreateParticipation(c.Context(), &p)
		if err != nil {
			return err
		}

		return c.JSON(resp)
	})

	participants.Put("/", func(c *fiber.Ctx) error {
		var p entities.Participation
		if err := c.BodyParser(&p); err != nil {
			return err
		}

		resp, err := d.UpdateParticipation(c.Context(), &p)
		if err != nil {
			return err
		}

		return c.JSON(resp)
	})

	participants.Get(":id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		p, err := d.GetParticipation(c.Context(), repos.ID(id))
		if err != nil {
			return err
		}

		return c.JSON(p)
	})
}
