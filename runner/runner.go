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
	rerunCounts    map[string]int
	taskTimers     map[string]*time.Timer
	maxReruns      int
}

type RunnerConfig struct {
	LogFolder      string
	SecretPath     string
	HandlerFactory utils.HandlerFactoryInterface
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
		rerunCounts:    make(map[string]int),
		taskTimers:     make(map[string]*time.Timer),
		maxReruns:      config.MaxReruns,
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

func (r *Runner) saveSecrets() {
	for key, crawler := range r.crawlers {
		marketSecret := utils.GetSecretName(key)
		utils.UpdateSecrets(crawler, *r.secretStore, marketSecret)
	}
}

func (r *Runner) cleanup() {
	// save secrets
	r.saveSecrets()
}

func (r *Runner) Run(tasks []types.CrawlerTask) {
	log.Printf("[START] running %v tasks: %v", len(tasks), tasks)
	defer r.cleanup()

	// run all tasks once at start
	for _, task := range tasks {
		r.RunTask(task)
	}

	// enter loop to rerun tasks until no more reruns
	tick := time.NewTicker(1 * time.Minute)
	defer tick.Stop()

	for range tick.C {
		allDone := true
		log.Printf("Checking tasks")
		for _, task := range tasks {
			if count, ok := r.rerunCounts[task.Name]; ok {
				log.Printf("Current task %s | run count: %d/%d", task.Name, count, r.maxReruns)
				if count < r.maxReruns {
					allDone = false
					break
				}
			}
		}

		// if all tasks are done, exit
		if allDone {
			log.Printf("All tasks are done")
			return
		}
	}
}

func (r *Runner) OnError(err error) {
	r.stopTimers()
}

func (r *Runner) stopTimers() {
	// stop all timers to prevent reruns
	for _, timer := range r.taskTimers {
		timer.Stop()
	}
}

// run crawling task
func (r *Runner) RunTask(task types.CrawlerTask) {
	// if not run before, init rerun count
	if _, ok := r.rerunCounts[task.Name]; !ok {
		r.rerunCounts[task.Name] = 0
	}

	var execTask = func(task types.CrawlerTask) error {
		for _, market := range task.Markets {
			// for each market run crawl tasks
			for taskName, taskConfig := range task.TaskConfigs {
				err := r.Crawl(market, task.Name, taskName, taskConfig)
				r.saveSecrets()
				if err != nil {
					log.Printf("Failed to crawl %s: %v", task.Name, err)
					r.OnError(err)
					return err
				}
			}
		}

		return nil
	}

	// attempt to run task and scedule reruns if needed
	r.rerunCounts[task.Name]++
	err := execTask(task)

	// if no error, reset rerun count
	if err == nil {
		// rerun logics
		if r.rerunCounts[task.Name] < r.maxReruns {
			waitDuration := time.Duration(task.RerunInterval) * time.Second
			log.Printf("Scheduling rerun %d for task %s in %v", r.rerunCounts[task.Name], task.Name, waitDuration)
			r.taskTimers[task.Name] = time.AfterFunc(waitDuration, func() {
				r.RunTask(task)
			})
		} else {
			log.Printf("Failed to run task %s after %d reruns", task.Name, r.maxReruns)
		}
	}
}

func (r *Runner) Crawl(market string, itemName, taskName string, config types.CrawlerConfig) error {
	crawler, err := r.GetCrawler(market)
	if err != nil {
		return err
	}

	switch taskName {
	case "listings":
		// crawl item listings
		listingHandler := r.handlerFactory.GetListingsHandler()
		err = crawler.CrawlItemListings(itemName, listingHandler, &config)
		if err != nil {
			return err
		}
	case "transactions":
		// crawl item transactions
		transactionHandler := r.handlerFactory.GetTransactionHandler()

		err = crawler.CrawlItemTransactions(itemName, transactionHandler, &config)
		if err != nil {
			return err
		}
	default:
		return errors.ErrTaskNotFound
	}

	return nil
}
