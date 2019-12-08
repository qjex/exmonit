package crawler

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type exmo struct {
	httpClient *http.Client
}

func (p *Pair) exmoFormat() string {
	return fmt.Sprintf("%s_%s", p.From, p.To)
}

func (e *exmo) Exchange() string {
	return "Exmo"
}

func (e *exmo) Crawl(ctx context.Context, pairs []Pair) (map[Pair]float64, error) {
	rq, err := http.NewRequestWithContext(ctx, "GET", "https://api.exmo.com/v1/ticker/", nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request")
	}

	data, err := getJson(e.httpClient, rq)
	if err != nil {
		return nil, err
	}

	if result, ok := data["result"].(bool); ok && !result {
		return nil, errors.Errorf("api error: %v", data)
	}

	res := make(map[Pair]float64)
	for _, pair := range pairs {
		pairData, ok := data[pair.exmoFormat()].(map[string]interface{})
		if !ok {
			log.Errorf("unexpected format for pair=%s", pair.String())
			continue
		}
		sellPrice, ok := pairData["sell_price"].(string)
		if !ok {
			log.Errorf("sell_price not found for pair=%s", pair.String())
			continue
		}
		price, err := strconv.ParseFloat(sellPrice, 64)
		if err != nil {
			log.Errorf("sell_price can't be parsed: %v", err)
			continue
		}
		res[pair] = price
	}
	return res, nil
}
