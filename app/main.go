package main

import (
	"context"
	"exmonit/api"
	"exmonit/crawler"
	"exmonit/storage"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var cfg struct {
	Host     string `long:"h" env:"DB_HOST" description:"db host:port"`
	User     string `long:"u" env:"DB_USER" description:"db user"`
	Password string `long:"p" env:"DB_PASSWORD" description:"db password"`
	Database string `long:"d" env:"DB_DATABASE" description:"database name" default:"exmonit"`
	Cfg      string `short:"f" env:"CONFIG" default:"config.yml" description:"config file"`
	Debug    bool   `short:"d" env:"DEBUG" description:"debug mode"`
}

func main() {
	if _, err := flags.Parse(&cfg); err != nil {
		os.Exit(1)
	}
	setupLog(cfg.Debug)
	conf, err := loadConfig(cfg.Cfg)
	if err != nil {
		log.Fatalf("can't load config %v", err)
	}

	s := storage.NewStorage(storage.Conf{
		Addr:     cfg.Host,
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
	})

	metricsServer := api.MetricsServer{}

	server := api.Server{
		Conf:    conf,
		Storage: s,
	}
	go server.Serve()
	go metricsServer.Serve()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Debug("SIGTERM caught")
		server.Close()
		metricsServer.Close()
		cancel()
	}()

	u := crawler.NewUpdater(s, prepare(conf.Pairs), conf.UpdateInterval)
	u.Run(ctx)
}

func setupLog(debug bool) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

func prepare(pairs []crawler.Pair) []crawler.Pair {
	for _, p := range pairs {
		p.To = strings.ToUpper(p.To)
		p.From = strings.ToUpper(p.From)
	}
	return pairs
}

func loadConfig(name string) (res *crawler.Conf, err error) {
	res = &crawler.Conf{}
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, res); err != nil {
		return nil, err
	}

	return res, nil
}
