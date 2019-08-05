package app

import (
	"encoding/json"

	"github.com/nats-io/stan.go"

	"github.com/kachan1208/statsWriter/src/dao"
	"github.com/kachan1208/statsWriter/src/model"
)

type counterHandler struct{}

var (
	counter = counterHandler{}
)

func (c counterHandler) process(msg *stan.Msg, repo *dao.StatsRepo) error {
	var counter model.Counter
	err := json.Unmarshal(msg.Data, &counter)
	if err != nil {
		return err
	}

	return repo.UpdateCounter(&counter)
}
