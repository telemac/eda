package main

import (
	"github.com/nats-io/nats.go/jetstream"
	"github.com/telemac/eda"
	natsbroker "github.com/telemac/eda/brokers/nats"
	"github.com/telemac/eda/edaentities"
	"github.com/telemac/eda/events"
	"github.com/telemac/eda/examples/nats-publisher/entities"
	"github.com/telemac/eda/internal/slogutils"
	"github.com/telemac/eda/pkg/registry"
	"github.com/telemac/goutils/task"
	"log/slog"
	"strconv"
	"time"
)

const COUNT = 1_000_000_000

func main() {
	ctx, cancel := task.NewCancellableContext(time.Second * 10)
	defer cancel()

	logger := slogutils.NewTintLogger(slog.LevelInfo, false, false, true)

	logger.Info("nats publisher example started")
	defer logger.Info("nats publisher example stopped")

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

	_, err = natsBroker.Js().CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "TEST",
		Subjects: []string{"test.>"},
		Storage:  jetstream.MemoryStorage,
		//Retention: jetstream.LimitsPolicy,
		//Discard:   jetstream.DiscardOld,
		MaxMsgs: 1_000_000,
		MaxAge:  time.Hour * 24,
	})
	if err != nil {
		logger.Error("create or update stream", "error", err)
		return
	}
	eventRegistry := registry.New()
	registry.Register[events.UserCreationRequested](eventRegistry)
	registry.Register[events.UserCreationDone](eventRegistry)
	registry.Register[entities.NamedCounter](eventRegistry)

	var lastCount int = -1

	natsBroker.Subscribe(ctx, "test.>", eventRegistry, func(topic string, event any) {
		//attrs := []any{
		//	slog.String("type", event.Type()),
		//	slog.String("PublishTopic", event.PublishTopic()),
		//	slog.String("SubscribeTopic", event.SubscribeTopic()),
		//	slog.Any("validate", event.Validate()),
		//	slog.Any("event", event),
		//}
		//logger.Info("received event", attrs...)
		switch event.(type) {
		case *entities.NamedCounter:
			counter := event.(*entities.NamedCounter)
			if counter.Count%1000 == 0 || counter.Count > COUNT-100 {
				logger.Info("received counter", "count", counter.Count)
			}
			if counter.Count-1 != lastCount {
				logger.Error("last count does not match",
					"count", counter.Count,
					"last-count", lastCount,
				)
				time.Sleep(time.Millisecond * 10)
			} else {
				lastCount = counter.Count
			}
		}

	})

	// create a user
	user := edaentities.User{
		FirstName: "Alexandre",
		LastName:  "HEIM",
	}
	// validate the user
	err = eda.Validate(user)
	if err != nil {
		logger.Error("validate user", "error", err)
	}

	// create a user creation requested event
	userCreationRequestedEvent := events.UserCreationRequested{
		User: user,
		Password: edaentities.Password{
			Password: "123456",
		},
	}

	err = eda.Validate(userCreationRequestedEvent)
	if err != nil {
		logger.Error("validate user creation requested event", "error", err)
	}

	err = natsBroker.Publish(ctx, "test."+userCreationRequestedEvent.PublishTopic(), userCreationRequestedEvent)
	if err != nil {
		logger.Error("publish user creation requested event", "error", err)
	}

	userCreationDoneEvent := events.UserCreationDone{
		UserCreationRequested: &userCreationRequestedEvent,
		Uuid:                  edaentities.Uuid{UUID: "005C6A52-85EB-4CA5-826A-3FE71864EBE0"},
	}
	err = eda.ValidateAll(userCreationDoneEvent)
	if err != nil {
		logger.Error("validate user creation done event", "error", err)
	}

	err = natsBroker.Publish(ctx, "test."+userCreationDoneEvent.Topic(), userCreationDoneEvent)
	if err != nil {
		logger.Error("publish user creation done event", "error", err)
	}

	counter := entities.NamedCounter{
		Name: "test-counter",
	}

	for i := range COUNT {
	retry:
		if task.IsCancelled(ctx) {
			break
		}
		counter.Count = i
		msgID := strconv.Itoa(i)
		_ = msgID
		err = natsBroker.PublishEvent(ctx, counter,
			natsbroker.PubAsync(),
			//natsbroker.PubMsgID(msgID),
			natsbroker.PubLogger(logger),
		)
		if err != nil {
			logger.Error("publish counter event", "error", err)
			goto retry
		} else {
			//logger.Info("publish counter event", "i", i)
		}
		//if i%1000 == 0 {
		//	err = natsBroker.Nc().Flush()
		//	if err != nil {
		//		logger.Error("flush nats", "error", err)
		//	}
		//}
	}

	logger.Warn("publish done")

	//<-ctx.Done()

}
