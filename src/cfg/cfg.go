package cfg

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	HTTPAddress       string `envconfig:"HTTP_ADDRESS" required:"true"`
	CassandraHost     string `envconfig:"CASSANDRA_HOST" required:"true"`
	CassandraPort     int    `envconfig:"CASSANDRA_PORT" required:"true"`
	CassandraLogin    string `envconfig:"CASSANDRA_LOGIN"`
	CassandraPassword string `envconfig:"CASSANDRA_PASSWORD"`
	CassandraKeyspace string `envconfig:"CASSANDRA_KEYSPACE" required:"true"`
	NatsAddress       string `envconfig:"NATS_HOST" required:"true"`
	NatsClient        string `envconfig:"NATS_CLIENT" required:"true"`
	NatsCluster       string `envconfig:"NATS_CLUSTER" required:"true"`
	StreamBatchSize   int    `envconfig:"STREAM_BATCH_SIZE" default:"1024"`
	ServiceName       string
}

var Config config

func init() {
	err := envconfig.Process("STATS", &Config)
	if err != nil {
		log.Fatalf("Can't load config: %s ", err.Error())
	}

	Config.ServiceName = "stats-writer"
}
