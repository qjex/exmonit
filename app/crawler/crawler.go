package crawler

import (
	"context"
	"exmonit/storage"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type Conf struct {
	Pairs          []Pair        `yaml:"pairs"`
	UpdateInterval time.Duration `yaml:"update"`
}

type Pair struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

func (p *Pair) String() string {
	return fmt.Sprintf("%s/%s", p.From, p.To)
}

type Crawler interface {
	Crawl(ctx context.Context, pairs []Pair) (map[Pair]float64, error)
	Exchange() string
}

type Updater struct {
	storage  *storage.Storage
	pairs    []Pair
	crawlers []Crawler
	interval time.Duration
}

func NewUpdater(storage *storage.Storage, pairs []Pair, interval time.Duration) *Updater {
	client := &http.Client{Timeout: 5 * time.Second}
	return &Updater{
		storage: storage,
		pairs:   pairs,
		crawlers: []Crawler{
			&binance{httpClient: client},
			&exmo{httpClient: client},
		},
		interval: interval,
	}
}

//Run runs Crawl method for each crawler in separate gorutines, saves crawled items it storage,
//and sleeps for "update" interval specified in config
func (u *Updater) Run(ctx context.Context) {
	for {
		log.Info("update started")

		var wg sync.WaitGroup
		wg.Add(len(u.crawlers))
		for _, c := range u.crawlers {
			go func(c Crawler) {
				defer wg.Done()
				log.Infof("starting crawling %s", c.Exchange())
				res, err := c.Crawl(ctx, u.pairs)
				if err != nil {
					log.Errorf("crawling for %s failed: %v", c.Exchange(), err)
					return
				}
				log.Infof("crawling %s finished", c.Exchange())

				select {
				case <-ctx.Done():
					return
				default:
				}

				u.store(c.Exchange(), res)

			}(c)

		}

		wg.Wait()
		log.Info("update finished")
		select {
		case <-ctx.Done():
			log.Debug("stopping update loop")
			return
		case <-time.After(u.interval):
		}

	}
}

func (u *Updater) store(exchange string, res map[Pair]float64) {
	for pair, rate := range res {
		rate := storage.Rate{
			Pair:     pair.String(),
			Exchange: exchange,
			Rate:     rate,
			Updated:  time.Now(),
		}
		err := u.storage.SaveOrUpdate(&rate)
		if err != nil {
			log.Errorf("error saving pair=%s,exchange=%s,rate=%f, %v", rate.Pair, rate.Exchange, rate.Rate, err)
			continue
		}
		log.Infof("pair=%s,exchange=%s,rate=%f saved", rate.Pair, rate.Exchange, rate.Rate)
	}
}