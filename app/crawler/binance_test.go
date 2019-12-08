package crawler

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_binance_Crawl(t *testing.T) {
	client := newTestClient(func(r *http.Request) *http.Response {
		symbol := r.URL.Query()["symbol"][0]
		var price string
		switch symbol {
		case "BTCUSDT":
			price = "1.1"
			break
		case "BTCRUB":
			price = "2.2"
			break
		default:
			assert.FailNow(t, "unexpected symbol")
		}
		rsp := fmt.Sprintf(`{"mins":5,"price":"%s"}`, price)
		return createResponse(200, rsp)

	})

	b := binance{httpClient: client}
	p1 := Pair{
		From: "BTC",
		To:   "USDT",
	}
	p2 := Pair{
		From: "BTC",
		To:   "RUB",
	}

	result, err := b.Crawl(context.Background(), []Pair{p1, p2})
	assert.Nil(t, err)
	assert.Equal(t, 1.1, result[p1])
	assert.Equal(t, 2.2, result[p2])
}

func Test_binance_Crawl_api_error(t *testing.T) {
	client := newTestClient(func(r *http.Request) *http.Response {
		symbol := r.URL.Query()["symbol"][0]
		assert.Equal(t, "BTCUSDT", symbol)
		rsp := fmt.Sprint(`{"code":-1121,"msg":"Invalid symbol."}`)
		return createResponse(200, rsp)

	})

	b := binance{httpClient: client}

	result, err := b.Crawl(context.Background(), []Pair{{
		From: "BTC",
		To:   "USDT",
	}})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func Test_binance_Exchange(t *testing.T) {
	b := binance{}
	assert.Equal(t, "Binance", b.Exchange())
}
