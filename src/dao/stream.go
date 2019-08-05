package dao

import (
	"github.com/nats-io/stan.go"
)

type StreamRepo struct {
	connection  stan.Conn
	streamName  string
	maxInFlight int
}

func NewStreamRepo(connection stan.Conn, streamName string, maxInFlight int) *StreamRepo {
	return &StreamRepo{
		connection:  connection,
		streamName:  streamName,
		maxInFlight: maxInFlight,
	}
}

func (s *StreamRepo) Subscribe(channel string, handler stan.MsgHandler) {
	//TODO: this params should be tested in cluster
	options := []stan.SubscriptionOption{
		stan.DurableName(s.streamName),
		stan.SetManualAckMode(),
		stan.MaxInflight(s.maxInFlight),
	}

	s.connection.Subscribe(channel, handler, options...)
}
