package storage

import (
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/pkg/errors"
)

type Rate struct {
	tableName struct{} `pg:"rate"`
	Pair      string   `pg:",pk"`
	Exchange  string   `pg:",pk"`
	Rate      float64
	Updated   time.Time
}

type Conf struct {
	Addr     string
	User     string
	Password string
	Database string
}

type Storage struct {
	db *pg.DB
}

func NewStorage(conf Conf) *Storage {
	db := pg.Connect(&pg.Options{
		Addr:         conf.Addr,
		User:         conf.User,
		Password:     conf.Password,
		Database:     conf.Database,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	return &Storage{db: db}
}

//SaveOrUpdate saved provided entity, if a rate with the same (pair, exchange) exists
// the entity is updated with new rate and updated time
func (s *Storage) SaveOrUpdate(rate *Rate) error {
	_, err := s.db.Model(rate).
		OnConflict("(pair, exchange) DO UPDATE").
		Set("rate = ?", rate.Rate).
		Set("updated = ?", rate.Updated).
		Insert()
	if err != nil {
		return errors.Wrap(err, "error on inserting or updating")
	}
	return nil
}

//FindAll returns all saved entities
func (s *Storage) FindAll() (rates []Rate, err error) {
	err = s.db.Model(&rates).Select()
	return rates, err
}
