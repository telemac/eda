package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	natsbroker "github.com/telemac/eda/brokers/nats"
	"github.com/telemac/eda/edaentities"
	"github.com/telemac/eda/event"
	"github.com/telemac/eda/events"
	"github.com/telemac/eda/internal/slogutils"
	"github.com/telemac/eda/pkg/registry"
	"github.com/telemac/goutils/cli"
	"github.com/telemac/goutils/task"
	"log/slog"
	"os"
	"time"
)

type CLI struct {
	Request bool `help:"send requests" default:"false"`
	Service bool `help:"run service" default:"false"`
}

func main() {
	var commandLine CLI
	command := cli.Parse(&commandLine)
	_ = command

	ctx, cancel := task.NewCancellableContext(time.Second * 10)
	defer cancel()

	logger := slogutils.NewTintLogger(slog.LevelInfo, false, false, true)

	logger.Info("nats micro example started")
	defer logger.Info("nats micro example stopped")

	if !commandLine.Service && !commandLine.Request {
		logger.Error("service or request must be true")
		return
	}

	// create a nats publisher broker
	natsBroker, err := natsbroker.NewNatsBroker(natsbroker.Config{
		//Host: "nats://megalarm:megalarm@192.168.1.210:4222",
		Host: "localhost",
		//Logger: logger,
	})
	if err != nil {
		logger.Error("create nats broker", "error", err)
		return
	}
	defer natsBroker.Close()

	// create event registry
	eventRegistry := registry.New()
	// register an event
	err = eventRegistry.Register(registry.RegistryEntry{
		EventType:    events.UserCreationRequested{}.Type(),
		EventFactory: func() any { return new(events.UserCreationRequested) },
	})
	if err != nil {
		logger.Error("register events.UserCreationRequested event", "error", err)
		return
	}
	// register another event using a helper function
	err = registry.Register[events.UserCreationDone](eventRegistry)
	if err != nil {
		logger.Error("register events.UserCreationDone event", "error", err)
		return
	}

	hostname, _ := os.Hostname()
	svc, err := micro.AddService(natsBroker.Nc(), micro.Config{
		Name:        "user",
		Version:     "0.0.2",
		Description: "A Simple user managment service",
		Metadata: map[string]string{
			"hostname": hostname,
		},
		ErrorHandler: func(service micro.Service, natsError *micro.NATSError) {
			logger.Error("nats micro service error", slog.Any("service", service), "error", natsError)
		},
	})
	if err != nil {
		logger.Error("create echo service", "error", err)
		return
	}
	defer svc.Stop()

	count := 1

	if commandLine.Service {

		userCreateHandler := func(req micro.Request) {
			eventType := req.Headers().Get(event.EDATypeHeader)

			logger.Info("create user",
				slog.String("data", string(req.Data())),
				slog.String("subject", req.Subject()),
				slog.String("event-type", eventType),
				slog.Any("headers", req.Headers()),
				slog.Int("count", count),
			)
			ucr, err := registry.UnmarshalEvent[events.UserCreationRequested](eventRegistry, eventType, req.Data())
			if err != nil {
				logger.Error("registry.UnmarshalEvent", "error", err)
				req.Error(err.Error(), "unmarshal event", nil)
				return
			}

			userCreationDone := events.UserCreationDone{
				UserCreationRequested: ucr,
				Uuid:                  edaentities.Uuid{UUID: "12345678-1234-1234-1234-123456789012"},
			}
			isRegistered := eventRegistry.IsRegistered(&userCreationDone)
			if !isRegistered {
				logger.Error("eventRegistry.IsRegistered", "error", "event not registered")
				req.Error("event not registered", "event not registered", nil)
				return
			}
			headers := micro.Headers{
				event.EDATypeHeader: []string{userCreationDone.Type()},
			}
			req.RespondJSON(userCreationDone, micro.WithHeaders(headers))
			count++
		}

		svc.AddEndpoint("create",
			micro.HandlerFunc(userCreateHandler),
			micro.WithEndpointMetadata(map[string]string{
				"description": "echo",
				"format":      "application/json",
			}),
			micro.WithEndpointSubject(events.UserCreationRequested{}.SubscribeTopic()),
		)
	}

	if commandLine.Request {
		var userCreationRequested events.UserCreationRequested
		userCreationRequested.User.FirstName = "Alexandre"
		userCreationRequested.User.LastName = "HEIM"

		var msg *nats.Msg
		for i := 1; i <= 1_000_000 && !task.IsCancelled(ctx); i++ {
			userCreationRequested.Password.Password = fmt.Sprintf("svc.echo.%d", i)
			msg, err = natsBroker.Request(userCreationRequested.PublishTopic(), userCreationRequested, time.Second)
			if err != nil {
				logger.Error("request", "error", err)
			} else {
				logger.Warn("request",
					"msg", string(msg.Data),
					"headers", msg.Header,
				)
			}
		}
	}

	<-ctx.Done()
}
