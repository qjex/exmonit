package storage

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/pkg/errors"
)

const dbRequestsTotalMetric = "db_requests_total"
const dbDurationMetric = "db_requests_duration"

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
	db      *pg.DB
	metrics metrics
}

type metrics struct {
	requests *prometheus.CounterVec
	duration prometheus.Summary
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
	m := metrics{
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: dbRequestsTotalMetric,
			},
			[]string{"type", "status"},
		),
		prometheus.NewSummary(
			prometheus.SummaryOpts{
				Name:       dbDurationMetric,
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			}),
	}
	prometheus.MustRegister(m.requests, m.duration)
	return &Storage{
		db:      db,
		metrics: m,
	}
}

//SaveOrUpdate saved provided entity, if a rate with the same (pair, exchange) exists
// the entity is updated with new rate and updated time
func (s *Storage) SaveOrUpdate(rate *Rate) error {
	timer := prometheus.NewTimer(s.metrics.duration)
	defer timer.ObserveDuration()

	_, err := s.db.Model(rate).
		OnConflict("(pair, exchange) DO UPDATE").
		Set("rate = ?", rate.Rate).
		Set("updated = ?", rate.Updated).
		Insert()
	if err != nil {
		s.metrics.requests.WithLabelValues("save", "fail").Inc()
		return errors.Wrap(err, "error on inserting or updating")
	}
	s.metrics.requests.WithLabelValues("save", "success").Inc()
	return nil
}

//FindAll returns all saved entities
func (s *Storage) FindAll() (rates []Rate, err error) {
	timer := prometheus.NewTimer(s.metrics.duration)
	defer timer.ObserveDuration()

	err = s.db.Model(&rates).Select()
	if err != nil {
		s.metrics.requests.WithLabelValues("findAll", "fail").Inc()
		return nil, err
	}
	s.metrics.requests.WithLabelValues("findAll", "success").Inc()
	return rates, nil
}
