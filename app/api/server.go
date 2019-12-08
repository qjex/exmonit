package api

import (
	"exmonit/crawler"
	"exmonit/storage"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type UpdateTime time.Time

func (j UpdateTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02 15:04:05.000"))), nil
}

type RateView struct {
	Pair     string     `json:"pair"`
	Exchange string     `json:"exchange"`
	Rate     float64    `json:"rate"`
	Updated  UpdateTime `json:"updated"`
}

type Server struct {
	Conf       *crawler.Conf
	Storage    *storage.Storage
	httpServer *http.Server
}

func (s *Server) Serve() {
	router := chi.NewRouter()
	router.Use(middleware.Throttle(1000), middleware.Timeout(60*time.Second))
	s.httpServer = &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
	router.Use(middleware.Heartbeat("/status"))
	router.Get("/get_rates", s.getRates)
	err := s.httpServer.ListenAndServe()
	log.Warningf("http server exited with %v", err)
}

func (s *Server) Close() {
	log.Debug("closing http server")
	_ = s.httpServer.Close()
}

func (s *Server) getRates(w http.ResponseWriter, r *http.Request) {
	saved, err := s.Storage.FindAll()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	result := make([]RateView, 0)
	for _, p := range saved {
		result = append(result, RateView{
			Pair:     p.Pair,
			Exchange: p.Exchange,
			Rate:     p.Rate,
			Updated:  UpdateTime(p.Updated),
		})
	}

	render.JSON(w, r, result)
}
