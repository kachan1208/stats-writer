package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gocql/gocql"
	"github.com/nats-io/stan.go"

	"github.com/kachan1208/statsWriter/src/app"
	"github.com/kachan1208/statsWriter/src/cfg"
	"github.com/kachan1208/statsWriter/src/dao"
	transport "github.com/kachan1208/statsWriter/src/transport/http"
)

func main() {
	c, _ := json.Marshal(cfg.Config)
	fmt.Println(string(c[:]))
	//TODO: Add signal handler(graceful)
	logger := initLogger()
	level.Info(logger).Log("msg", "Service starting")

	gocqlSess := initSess()
	defer gocqlSess.Close()
	statsRepo := dao.NewStatsRepo(gocqlSess)
	level.Info(logger).Log("msg", "DB session created")

	natsConnection := initNats(logger)
	defer natsConnection.Close()
	stream := dao.NewStreamRepo(natsConnection, cfg.Config.ServiceName, cfg.Config.StreamBatchSize)

	level.Info(logger).Log("msg", "Nats connection created")

	service := app.NewService(stream, statsRepo, logger)
	initHTTPListener(logger)
	level.Info(logger).Log("msg", "Service started", "address:", cfg.Config.HTTPAddress)

	service.Run()
}

func initSess() *gocql.Session {
	cluster := gocql.NewCluster(cfg.Config.CassandraHost)
	if cfg.Config.CassandraLogin != "" && cfg.Config.CassandraPassword != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			cfg.Config.CassandraLogin,
			cfg.Config.CassandraPassword,
		}
	}

	cluster.Keyspace = cfg.Config.CassandraKeyspace
	cluster.Consistency = gocql.LocalQuorum
	cluster.Port = cfg.Config.CassandraPort
	sess, err := cluster.CreateSession()
	//TODO: Add waiting for connection process
	if err != nil {
		panic(err)
	}

	return sess
}

func initNats(logger log.Logger) stan.Conn {
	c, err := stan.Connect(cfg.Config.NatsCluster, cfg.Config.NatsClient,
		stan.NatsURL(cfg.Config.NatsAddress),
		stan.Pings(10, 5),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			level.Error(logger).Log("err", "Exit, connection to NATS lost")
		}))

	//TODO: Add waiting for connection process
	if err != nil {
		panic(err)
	}

	return c
}

func initLogger() log.Logger {
	var logger log.Logger

	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.NewSyncLogger(logger)
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(logger,
		"service", "stats-writer",
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)

	return logger
}

func initHTTPListener(logger log.Logger) {
	handler := transport.NewHandler(cfg.Config.HTTPAddress)

	httpServer := http.Server{
		Addr:    handler.Address,
		Handler: handler.Router,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			level.Error(logger).Log(err)
		}
	}()
}
