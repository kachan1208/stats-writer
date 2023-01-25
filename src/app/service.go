package app

import (
	"encoding/json"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nats-io/stan.go"

	"github.com/kachan1208/statsWriter/src/dao"
	"github.com/kachan1208/statsWriter/src/model"
)

type MsgHandler interface {
	process(string, []byte, *dao.StatsRepo) error
}

type Service struct {
	stream      *dao.StreamRepo
	stats       *dao.StatsRepo
	log         log.Logger
	errs        chan error
	msgHandlers map[string]MsgHandler
}

func NewService(stream *dao.StreamRepo, stats *dao.StatsRepo, log log.Logger) *Service {
	s := &Service{
		stream:      stream,
		stats:       stats,
		log:         log,
		errs:        make(chan error),
		msgHandlers: make(map[string]MsgHandler),
	}

	s.registerHandler("counter", counter)

	return s
}

func (s *Service) Run() {
	s.stream.Subscribe("stats", s.handle)

	//lock
	level.Error(s.log).Log("err", <-s.errs)
}

func (s *Service) handle(m *stan.Msg) {
	var (
		err error
		msg model.Message
	)
	err = json.Unmarshal(m.Data, &msg)
	if err != nil {
		level.Error(s.log).Log(err)
	}

	if h, ok := s.msgHandlers[msg.Subject]; ok {
		err = h.process(msg.AccountID, msg.Body, s.stats)
	} else {
		level.Error(s.log).Log("err", "Handler for this type is not registered %s", m.Subject)
	}

	if err != nil {
		level.Error(s.log).Log(err)
	}

	m.Ack()
}

func (s *Service) registerHandler(name string, handler MsgHandler) {
	if _, ok := s.msgHandlers[name]; ok {
		level.Error(s.log).Log("msg", "Can't register same handler two times %s", name)
		return
	}

	s.msgHandlers[name] = handler
}
