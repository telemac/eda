package natsbroker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nats.go/micro"
	"github.com/telemac/eda/broker"
	"github.com/telemac/eda/event"
	"github.com/telemac/eda/pkg/registry"
	"log/slog"
	"time"
)

type NatsBroker struct {
	broker.Subscriber[event.Eventer]
	config Config
	nc     *nats.Conn
	js     jetstream.JetStream // jetstream context
}

type Config struct {
	Host   string
	Logger *slog.Logger
}

func NewNatsBroker(config Config) (*NatsBroker, error) {
	var err error
	natsBroker := &NatsBroker{config: config}
	if natsBroker.config.Logger == nil {
		natsBroker.config.Logger = slog.Default()
	}

	natsBroker.nc, err = nats.Connect(config.Host, nats.ErrorHandler(func(conn *nats.Conn, subscription *nats.Subscription, err error) {
		if err != nil {
			if subscription != nil {
				natsBroker.config.Logger.Error("nats error",
					"error", err,
					slog.String("subject", subscription.Subject),
					slog.String("queue", subscription.Queue),
				)
			} else {
				natsBroker.config.Logger.Error("nats error",
					"error", err,
				)
			}
		}
	}))

	if err != nil {
		slog.Error("connect to nats broker", "error", err, "host", config.Host)
		return nil, err
	}
	natsBroker.js, err = jetstream.New(natsBroker.nc)
	if err != nil {
		slog.Error("jetstream new", "error", err)
		natsBroker.nc.Close()
		return nil, err
	}

	return natsBroker, nil
}

func (natsBroker *NatsBroker) Close() {
	var err error
	natsBroker.config.Logger.Debug("close nats broker")
	err = natsBroker.nc.Flush()
	if err != nil {
		natsBroker.config.Logger.Error("flush nats broker", "error", err)
	}
	err = natsBroker.nc.Drain()
	if err != nil {
		natsBroker.config.Logger.Error("drain nats broker", "error", err)
	}

	for i := 0; i < 100; i++ {
		if natsBroker.nc.IsDraining() {
			natsBroker.config.Logger.Debug("is draining", "i", i)
			time.Sleep(time.Millisecond * 100)
		}
	}
	if natsBroker.nc.IsDraining() {
		natsBroker.config.Logger.Error("nats did not drain all messages")
	}
	natsBroker.nc.Close()
}

type pubOpts struct {
	Async  bool
	Logger *slog.Logger
	MsgID  string
}

type PublishOpt func(*pubOpts) error

func PubLogger(logger *slog.Logger) PublishOpt {
	return func(o *pubOpts) error {
		if logger == nil {
			return errors.New("nil logger")
		}
		o.Logger = logger
		return nil
	}
}

func PubAsync() PublishOpt {
	return func(o *pubOpts) error {
		o.Async = true
		return nil
	}
}

func PubMsgID(msgID string) PublishOpt {
	return func(o *pubOpts) error {
		o.MsgID = msgID
		return nil
	}
}

// Publish publishes a message to the broker
func (natsBroker *NatsBroker) Publish(ctx context.Context, topic string, eventer event.Eventer, opts ...PublishOpt) error {
	var err error
	var o pubOpts
	for _, opt := range opts {
		err = opt(&o)
		if err != nil {
			return fmt.Errorf("publish option : %w", err)
		}
	}

	if err := eventer.Validate(); err != nil {
		return fmt.Errorf("validate event %T : %w", eventer, err)
	}

	msg := nats.NewMsg(topic)
	msg.Header.Set(event.EDATypeHeader, eventer.Type())
	msg.Data, err = json.Marshal(eventer)
	if err != nil {
		return err
	}
	if o.MsgID != "" {
		msg.Header.Set("Nats-Msg-Id", o.MsgID)
	}
	if o.Async {
		_, err = natsBroker.js.PublishMsgAsync(msg)
	} else {
		_, err = natsBroker.js.PublishMsg(ctx, msg)
	}
	if err != nil {
		// if error is not a cancelled context, log it
		if o.Logger != nil {
			if !errors.Is(err, context.Canceled) {
				slog.Error("publish event", "error", err, "topic", msg.Subject, "payload", string(msg.Data))
			} else {
				slog.Warn("publish event", "error", err, "topic", msg.Subject, "payload", string(msg.Data))
			}
		}
		return err
	}
	slog.Debug("publish event", "topic", msg.Subject, "payload", string(msg.Data))
	return nil
}

// PublishEvent publishes a message to the broker
func (natsBroker *NatsBroker) PublishEvent(ctx context.Context, event event.Eventer, opts ...PublishOpt) error {
	return natsBroker.Publish(ctx, event.PublishTopic(), event, opts...)
}

// Request publishes a message to the broker and waits for a response
func (natsBroker *NatsBroker) Request(topic string, eventer event.Eventer, timeout time.Duration) (*nats.Msg, error) {
	var err error
	msg := nats.NewMsg(topic)
	msg.Header.Set(event.EDATypeHeader, eventer.Type())
	msg.Data, err = json.Marshal(eventer)
	if err != nil {
		return nil, err
	}
	return natsBroker.nc.RequestMsg(msg, timeout)
}

// Subscribe subscribes to a topic
func (natsBroker *NatsBroker) Subscribe(ctx context.Context, subscribeTopic string, eventRegistry registry.Registry[event.Eventer], callback func(topic string, event any)) error {
	_, err := natsBroker.nc.Subscribe(subscribeTopic, func(msg *nats.Msg) {
		// TODO : make concrete event from event type/topic
		eventType := msg.Header.Get(event.EDATypeHeader)
		subscribedSubject := msg.Sub.Subject
		_ = subscribedSubject
		if eventType == "" {
			slog.Error("missing event type", "topic", msg.Subject, "payload", string(msg.Data))
			return
		}
		eventer, err := eventRegistry.New(eventType)
		err = json.Unmarshal(msg.Data, &eventer)
		if err != nil {
			slog.Error("unmarshal event", "error", err, "topic", subscribeTopic, "payload", string(msg.Data))
			return
		}
		//err = mapstructure.Decode(unknownEvent, &event)
		//if err != nil {
		//	slog.Error("mapstructure decode", "error", err, "topic", topic, "payload", string(msg.Data))
		//	return
		//}
		callback(subscribeTopic, eventer)
	})
	if err != nil {
		slog.Error("subscribe to topic", "error", err, "topic", subscribeTopic)
		return err
	}
	//go func() {
	//	<-ctx.Done()
	//	slog.Warn("context done", "topic", topic)
	//	err := sub.Drain()
	//	if err != nil {
	//		slog.Error("drain subscription", "error", err, "topic", topic)
	//	}
	//	err = sub.Unsubscribe()
	//	if err != nil {
	//		slog.Error("unsubscribe subscription", "error", err, "topic", topic)
	//	}
	//}()
	return nil
}

// Nc returns the nats connection
func (natsBroker *NatsBroker) Nc() *nats.Conn {
	return natsBroker.nc
}

// Js returns the jetstream context
func (natsBroker *NatsBroker) Js() jetstream.JetStream {
	return natsBroker.js
}

// HandlerFunc2EventHandler converts a micro.HandlerFunc to a broker.EventHandler
func HandlerFunc2EventHandler(eh broker.EventHandler, eventRegistry registry.Registry[event.Eventer]) micro.HandlerFunc {
	return func(req micro.Request) {
		slog.Info("HandlerFunc2EventHandler", "data", string(req.Data()))
		// get reqEvent type from headers
		eventType := req.Headers().Get(event.EDATypeHeader)
		if eventType == "" {
			slog.Error("HandlerFunc2EventHandler",
				"error", "missing event type",
				"request", req,
			)
			req.Error("missing event type", "error description", nil)
			return
		}

		// create reqEvent from req.Data
		reqEventer, err := eventRegistry.New(eventType)
		if err != nil {
			slog.Error("HandlerFunc2EventHandler", "error", err)
			req.Error(err.Error(), "unknown event type "+eventType, nil)
			return
		}

		data := req.Data()
		err = json.Unmarshal(data, reqEventer)
		if err != nil {
			slog.Error("HandlerFunc2EventHandler", "error", err)
			req.Error(err.Error(), "unmarshal payload", nil)
			return
		}

		// call event handler
		resultEventer, err := eh(reqEventer)
		if err != nil {
			slog.Error("HandlerFunc2EventHandler", "error", err)
			req.Error(err.Error(), "call event handler", nil)
			return
		}

		// respond with resultEventer
		var headers micro.Headers
		if resultEventer != nil {
			resultEventType := fmt.Sprintf("%T", resultEventer)
			headers = micro.Headers{
				event.EDATypeHeader: []string{resultEventType},
			}
		}
		req.RespondJSON(resultEventer, micro.WithHeaders(headers))
	}
}
