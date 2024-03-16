# Event Driver Architecture (EDA)

## What is EDA ?
EDA is a set of helper libraries to build event driven applications.

## Objectives
- [ ] Broker abstraction
- [ ] Event registry by type / topic
- [ ] Error handling
- [ ] Logging with slog
- [ ] Metrics
- [ ] Tracing

## Events
- [x] Has a PublishTopic() using the nats convention (. separated)
- [x] Has a SubscribeTopic() using the nats convention (. separated, * and > wildcard)
- [ ] Has an ID

## usefil links
[slog: NATS handler](https://github.com/samber/slog-nats/tree/main)
[Oops - Error handling with context, assertion, stack trace and source fragments](https://github.com/samber/oops)
[go-coffeeshop](https://github.com/thangchung/go-coffeeshop/tree/main)