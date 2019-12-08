package crawler

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type binance struct {
	httpClient *http.Client
}

func (e *binance) Exchange() string {
	return "Binance"
}

func (e *binance) Crawl(ctx context.Context, pairs []Pair) (map[Pair]float64, error) {
	res := make(map[Pair]float64)

	for _, pair := range pairs {
		logger := log.WithFields(log.Fields{
			"exchange": e.Exchange(),
			"pair":     pair.String(),
		})
		url := fmt.Sprintf("https://api.binance.com/api/v3/avgPrice?symbol=%s%s", pair.From, pair.To)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, errors.Wrap(err, "error creating request")
		}

		data, err := getJson(e.httpClient, req)
		if err != nil {
			logger.Error(err)
			continue
		}

		select {
		case <-ctx.Done():
			logger.Debug("crawling stopped by context cancellation")
			return nil, ctx.Err()
		default:
		}

		if _, ok := data["code"]; ok {
			logger.Errorf("api error: %v", data)
			continue
		}

		if price, ok := data["price"].(string); ok {
			p, err := strconv.ParseFloat(price, 64)
			if err != nil {
				logger.Error("price can't be parsed as float")
			}
			res[pair] = p
		}
	}

	return res, nil
}
