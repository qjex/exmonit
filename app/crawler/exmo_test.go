package crawler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPair_exmoFormat(t *testing.T) {
	assert.Equal(t, "BTC_USDT", (&Pair{
		From: "BTC",
		To:   "USDT",
	}).exmoFormat())
}

func Test_exmo_Crawl(t *testing.T) {
	client := newTestClient(func(r *http.Request) *http.Response {
		return createResponse(200, `
{
	"BTC_USD":{"buy_price":"7786.06","sell_price":"7798.72988289","updated":1575806999},
	"BTC_RUB":{"buy_price":"493300","sell_price":"494400","updated":1575806999},
	"BTC_EUR":{"buy_price":"6804.12473735","sell_price":"6817.12473735","updated":1575806998}
}`)

	})

	b := exmo{httpClient: client}
	p1 := Pair{
		From: "BTC",
		To:   "USD",
	}
	p2 := Pair{
		From: "BTC",
		To:   "RUB",
	}

	result, err := b.Crawl(context.Background(), []Pair{p1, p2})
	assert.Nil(t, err)
	assert.Equal(t, 7798.72988289, result[p1])
	assert.Equal(t, float64(494400), result[p2])
}

func Test_exmo_Crawl_no_symbol(t *testing.T) {
	client := newTestClient(func(r *http.Request) *http.Response {
		return createResponse(200, `
{
	"BTC_USD":{"buy_price":"7786.06","sell_price":"7798.72988289","updated":1575806999},
	"BTC_RUB":{"buy_price":"493300","updated":1575806999}
}`)

	})

	b := exmo{httpClient: client}

	result, err := b.Crawl(context.Background(), []Pair{
		{
			From: "BTC",
			To:   "USDT",
		},
		{
			From: "BTC",
			To:   "RUB",
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func Test_exmo_Crawl_api_error(t *testing.T) {
	client := newTestClient(func(r *http.Request) *http.Response {
		return createResponse(200, `{"result":false,"error":"Error 40015: API function do not exist "}`)

	})

	b := exmo{httpClient: client}

	_, err := b.Crawl(context.Background(), []Pair{
		{
			From: "BTC",
			To:   "USDT",
		},
	})
	assert.NotNil(t, err)
}

func Test_exmo_Exchange(t *testing.T) {
	assert.Equal(t, "Exmo", (&exmo{}).Exchange())
}
