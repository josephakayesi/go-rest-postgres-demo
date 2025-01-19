package config

import (
	"context"
	"fmt"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	USER_CREATED_EVENT_SUBJECT         string = "user.created.event"
	USER_LOGGED_IN_EVENT_SUBJECT       string = "user.logged.in.event"
	USER_REFRESHED_TOKEN_EVENT_SUBJECT string = "user.refreshed.token.event"
)

type NatsQueue struct {
	UserCreatedEventPublisher        *NatsPublisher
	UserLoggedInEventPublisher       *NatsPublisher
	UserRefreshedTokenEventPublisher *NatsPublisher
}

type NatsPublisher struct {
	connection jetstream.JetStream
	subject    string
}

type NatsSubscriber struct {
	connection jetstream.JetStream
	subject    string
}

func (publisher *NatsPublisher) Publish(event []byte) error {
	ctx := context.Background()

	_, err := publisher.connection.Publish(ctx, publisher.subject, event)

	if err != nil {
		return err
	}
	return nil
}

func ConnectToNats(c *Config) (*nats.Conn, jetstream.JetStream, error) {
	fmt.Println("establishing connection to nats")
	fmt.Println("Which environment in nats connection", GetEnvironment())
	if GetEnvironment() == Development {
		return connectToNatsDevelopment(c)
	}

	return connectToNatsProduction(c)
}

func connectToNatsDevelopment(c *Config) (*nats.Conn, jetstream.JetStream, error) {
	fmt.Println("establishing connection to nats in development")

	conn, err := nats.Connect(c.NATS_URL)

	if err != nil {
		fmt.Println("error connecting to nats: ", err)
		return nil, nil, err
	}

	stan, err := jetstream.New(conn)
	if err != nil {
		fmt.Println("error connecting to jetstream: ", err)
		conn.Close()
		return nil, nil, err
	}

	fmt.Println("connection established to nats in development: ", conn.Opts.Url)
	fmt.Println("jetstream established to nats in development: ", stan)

	return conn, stan, nil
}

func connectToNatsProduction(c *Config) (*nats.Conn, jetstream.JetStream, error) {
	fmt.Println("establishing connection to nats in prodcution")

	conn, err := nats.Connect(c.NATS_URL, nats.UserJWTAndSeed(c.NATS_JWT, c.NATS_SEED),
		nats.Name("auth"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second))

	if err != nil {
		fmt.Println("error connecting to nats: ", err)
		return nil, nil, err
	}

	stan, err := jetstream.New(conn)
	if err != nil {
		conn.Close()
		fmt.Println("error connecting to jetstream: ", err)
		return nil, nil, err
	}

	fmt.Println("connection established to nats in prodcution: ", conn.Opts.Url)
	fmt.Println("jetstream established to nats in prodcution: ", stan)

	return conn, stan, nil
}

func NewNatsPublisher(stan jetstream.JetStream, subject string) *NatsPublisher {
	return &NatsPublisher{
		connection: stan,
		subject:    subject,
	}
}

func NewNatsSubscriber(stan jetstream.JetStream, topic string) *NatsSubscriber {
	return &NatsSubscriber{
		connection: stan,
		subject:    topic,
	}
}

// CreateOrUpdateStreams : Creates or updates streams for auth
func CreateOrUpdateStreams(ctx context.Context, stan jetstream.JetStream) error {

	_, err := stan.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "auth",
		Description: "messages for auth",
		Subjects: []string{
			USER_CREATED_EVENT_SUBJECT,
			USER_LOGGED_IN_EVENT_SUBJECT,
			USER_REFRESHED_TOKEN_EVENT_SUBJECT,
		},
		MaxBytes: 1024 * 1024 * 1024,
	})

	if err != nil {
		return err
	}

	return nil
}
