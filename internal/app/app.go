package app

import (
	"TradeBot/internal/config"
	"TradeBot/internal/domain"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/skwizi4/lib/ErrChan"
	logger "github.com/skwizi4/lib/logs"
	tradermade "github.com/tradermade/Go-SDK/rest"
	"os"
)

type App struct {
	appName        string
	logger         logger.GoLogger
	cfg            config.Config
	errChan        *ErrChan.ErrorChannel
	validator      *validator.Validate
	restClient     *tradermade.RESTClient
	timeSeriesRate *tradermade.TimeSeriesRate
	currency       string
	startDate      string
	endDate        string
	interval       string
}

func New(appName string) App {
	return App{
		appName: appName,
		logger:  logger.InitLogger(),
	}
}

func (a *App) Run(ctx context.Context, testRange string, interval string) {
	a.InitValidator()
	a.PopulateConfig()
	a.InitErrHandlers(ctx)
	a.InitRestClient(a.cfg.CFD.ApiKey)
	a.ParseBacktestSettings(testRange, interval)
	a.FetchData()

}

func (a *App) InitErrHandlers(ctx context.Context) {
	a.errChan = ErrChan.InitErrChan(10, a.logger)
	go func() {
		for {
			select {
			case <-ctx.Done():
				a.errChan.Close()
				return
			}
		}
	}()
	a.errChan.Start()
	a.logger.InfoFrmt("InitErrorHandler-Successfully")
}
func (a *App) InitValidator() {
	a.validator = validator.New()
	a.logger.InfoFrmt("initValidator-Successfully")
}

func (a *App) PopulateConfig() {
	cfg, err := config.ParseConfig("/Users/allright6789012gmail.com/Desktop/TradeBot/tsconfig.json")
	if err != nil {
		a.logger.ErrorFrmt("error in parsing config: %s", err)
	}

	err = cfg.ValidateConfig(a.validator)
	if err != nil {
		a.logger.ErrorFrmt("error in config validation: %s", err)
	}
	a.cfg = *cfg
	a.logger.InfoFrmt("InitConfig-Successfully")
}
func (a *App) InitRestClient(apiKey string) *tradermade.RESTClient {
	client := tradermade.NewRESTClient(apiKey)
	if client == nil {
		a.logger.ErrorFrmt("error in creating rest client")
		os.Exit(1)
	}
	return client
}

func (a *App) ParseBacktestSettings(testRange string, interval string) {
	a.startDate = "" // todo: implement getting start date from testRange
	a.endDate = ""   // todo: implement getting end date from testRange
	a.interval = ""  // todo: implement getting interval from str
}
func (a *App) FetchData() {
	var err error
	a.timeSeriesRate, err = a.restClient.GetTimeSeriesData(a.currency, a.startDate, a.endDate, a.interval, domain.Period)
	if err != nil {
		a.logger.ErrorFrmt("error in fetching data: %s", err)
		os.Exit(1)
	}
}
