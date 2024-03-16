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
	eventRegistry := event.NewEventRegistry()
	//eventRegistry.Register(events.UserCreationRequested{})
	eventRegistry.Register(event.Factory[events.UserCreationRequested]())
	eventRegistry.Register(&events.UserCreationDone{})

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

		userCreateHandler2 := func(requestEvent event.Eventer) (any, error) {
			requestEvent, err := eventRegistry.New(requestEvent.Type())
			if err != nil {
				return nil, err
			}
			ucr, ok := requestEvent.(*events.UserCreationRequested)
			if !ok {
				return nil, fmt.Errorf("invalid event type %T", requestEvent)
			}

			ucr.Password.Password = ucr.Password.Password + " modified"

			return ucr, nil
		}

		userCreateHandler := func(req micro.Request) {
			eventType := req.Headers().Get(event.EDATypeHeader)
			logger.Info("create user",
				slog.String("data", string(req.Data())),
				slog.String("subject", req.Subject()),
				slog.String("event-type", eventType),
				slog.Any("headers", req.Headers()),
				slog.Int("count", count),
			)
			userCreationDone := events.UserCreationDone{
				UserCreationRequested: events.UserCreationRequested{},
				Uuid:                  edaentities.Uuid{UUID: "12345678-1234-1234-1234-123456789012"},
			}
			headers := micro.Headers{
				event.EDATypeHeader: []string{userCreationDone.Type()},
			}
			req.RespondJSON(userCreationDone, micro.WithHeaders(headers))
			count++
		}
		_ = userCreateHandler

		//m := svc.AddGroup("user")
		svc.AddEndpoint("create",
			//micro.HandlerFunc(userCreateHandler),
			natsbroker.HandlerFunc2EventHandler(userCreateHandler2, eventRegistry),
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
