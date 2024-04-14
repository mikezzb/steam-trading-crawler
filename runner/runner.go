package runner

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/mikezzb/steam-trading-crawler/buff"
	"github.com/mikezzb/steam-trading-crawler/errors"
	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
)

type Runner struct {
	handlerFactory utils.HandlerFactoryInterface
	secretStore    *shared.PersisitedStore
	crawlers       map[string]types.Crawler
}

type RunnerConfig struct {
	LogFolder      string
	SecretPath     string
	HandlerFactory utils.HandlerFactoryInterface
	RerunInterval  int // in seconds
	MaxReruns      int
}

func NewRunner(config *RunnerConfig) (*Runner, error) {
	// init log
	logName := time.Now().Format("2006-01-02") + ".log"
	logFile, err := os.OpenFile(path.Join(config.LogFolder, logName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(logFile)

	if err != nil {
		log.Fatalf("Failed to init log: %v", err)
		return nil, err
	}

	// init secret store
	secretStore, err := shared.NewPersisitedStore(config.SecretPath)

	if err != nil {
		log.Fatalf("Failed to init secret store: %v", err)
		return nil, err
	}

	runner := &Runner{
		secretStore:    secretStore,
		handlerFactory: config.HandlerFactory,
		crawlers:       make(map[string]types.Crawler),
	}

	return runner, nil
}

func (r *Runner) GetCrawler(marketName string) (types.Crawler, error) {
	// if crawler already exists, return it
	if crawler, ok := r.crawlers[marketName]; ok {
		return crawler, nil
	}

	switch marketName {
	case "buff":
		crawler := &buff.BuffCrawler{}
		buffSecret := r.secretStore.Get(utils.GetSecretName(marketName)).(string)
		err := crawler.Init(buffSecret)
		if err != nil {
			return nil, err
		}
		r.crawlers[marketName] = crawler
		return crawler, nil
	default:
		return nil, errors.ErrMarketNotFound
	}
}

func (r *Runner) cleanup() {
	// save secrets
	for key, crawler := range r.crawlers {
		marketSecret := utils.GetSecretName(key)
		utils.UpdateSecrets(crawler, *r.secretStore, marketSecret)
	}
}

func (r *Runner) Run(tasks []types.CrawlerTask) {
	log.Printf("[%v] Start running %v tasks: %v", time.Now(), len(tasks), tasks)
	defer r.cleanup()

	for idx, task := range tasks {
		r.RunTask(task)

		if idx < len(tasks)-1 {
			shared.RandomSleep(40, 80)
		}
	}
}

// run crawling task
func (r *Runner) RunTask(task types.CrawlerTask) {
	for _, market := range task.Markets {
		// TODO: can run in parallel since market is independent
		err := r.Crawl(market, task.Name, task.Config)

		if err != nil {
			log.Printf("Failed to crawl %s: %v", task.Name, err)
		}
	}
}

func (r *Runner) Crawl(market string, name string, config types.CrawlerConfig) error {
	crawler, err := r.GetCrawler(market)
	if err != nil {
		return err
	}
	// crawl item transactions
	transactionHandler := r.handlerFactory.GetTransactionHandler()

	err = crawler.CrawlItemTransactions(name, transactionHandler, &config)
	if err != nil {
		return err
	}

	shared.RandomSleep(10, 30)

	// crawl item listings
	listingHandler := r.handlerFactory.GetListingsHandler()
	err = crawler.CrawlItemListings(name, listingHandler, &config)
	if err != nil {
		return err
	}

	return nil
}
