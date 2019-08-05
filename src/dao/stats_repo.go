package dao

import (
	"github.com/gocql/gocql"

	"github.com/kachan1208/statsWriter/src/model"
)

type StatsRepo struct {
	db *gocql.Session
}

func NewStatsRepo(db *gocql.Session) *StatsRepo {
	return &StatsRepo{
		db: db,
	}
}

func (t *StatsRepo) UpdateCounter(counter *model.Counter) error {
	return t.db.Query(`
		UPDATE counters 
		SET count = count + ?
		WHERE 
			account_id = ?
			AND action_id = ?`,
		counter.Count,
		counter.AccountID,
		counter.ActionID).Exec()
}
