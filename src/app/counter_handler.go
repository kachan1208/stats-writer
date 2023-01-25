package app

import (
	"encoding/json"

	"github.com/kachan1208/statsWriter/src/dao"
	"github.com/kachan1208/statsWriter/src/model"
)

type counterHandler struct{}

var (
	counter = counterHandler{}
)

func (c counterHandler) process(accID string, msg []byte, repo *dao.StatsRepo) error {
	var counter model.Counter
	err := json.Unmarshal(msg, &counter)
	if err != nil {
		return err
	}

	counter.AccountID = accID

	return repo.UpdateCounter(&counter)
}
