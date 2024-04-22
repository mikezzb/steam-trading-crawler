package runner

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/mikezzb/steam-trading-crawler/crawler/buff"
	"github.com/mikezzb/steam-trading-crawler/crawler/igxe"
	"github.com/mikezzb/steam-trading-crawler/errors"
	"github.com/mikezzb/steam-trading-crawler/handler"
	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
)

type Runner struct {
	handlerFactory   handler.IHandlerFactory
	secretStore      *shared.JsonKvStore
	taskHistoryStore *shared.JsonKvStore
	crawlers         map[string]types.ICrawler
	rerunCounts      map[string]int
	taskTimers       map[string]*time.Timer
	maxReruns        int
	marketErrors     map[string]error
}

type RunnerConfig struct {
	LogFolder        string
	SecretStore      *shared.JsonKvStore
	TaskHistoryStore *shared.JsonKvStore
	HandlerFactory   handler.IHandlerFactory
	MaxReruns        int
	TaskHistoryPath  string
}

const (
	DEFAULT_TASK_HISTORY_PATH = "tmp/taskHistory.json"
)

func NewRunner(config *RunnerConfig) (*Runner, error) {
	// init log
	logName := time.Now().Format("2006-01-02") + ".log"
	logFile, err := os.OpenFile(path.Join(config.LogFolder, logName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(logFile)

	if err != nil {
		log.Fatalf("Failed to init log: %v", err)
		return nil, err
	}

	// run history store
	var taskHistoryStore *shared.JsonKvStore
	if config.TaskHistoryStore != nil {
		taskHistoryStore, _ = shared.NewJsonKvStore(config.TaskHistoryPath)
	} else {
		taskHistoryStore, _ = shared.NewJsonKvStore(DEFAULT_TASK_HISTORY_PATH)
	}

	runner := &Runner{
		secretStore:      config.SecretStore,
		handlerFactory:   config.HandlerFactory,
		crawlers:         make(map[string]types.ICrawler),
		rerunCounts:      make(map[string]int),
		taskTimers:       make(map[string]*time.Timer),
		maxReruns:        config.MaxReruns,
		taskHistoryStore: taskHistoryStore,
		marketErrors:     make(map[string]error),
	}

	return runner, nil
}

func (r *Runner) GetCrawler(marketName string) (types.ICrawler, error) {
	// if crawler already exists, return it
	if crawler, ok := r.crawlers[marketName]; ok {
		return crawler, nil
	}

	switch marketName {
	case "buff":
		buffSecret := r.secretStore.Get(utils.GetSecretName(marketName)).(string)
		crawler, err := buff.NewCrawler(buffSecret)
		if err != nil {
			return nil, err
		}
		r.crawlers[marketName] = crawler
		return crawler, nil
	case "igxe":
		igxeSecret := r.secretStore.Get(utils.GetSecretName(marketName)).(string)
		crawler, err := igxe.NewCrawler(igxeSecret)
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
		utils.UpdateSecrets(crawler, r.secretStore, marketSecret)
	}
}

func (r *Runner) cleanup() {
	// save secrets
	r.saveSecrets()
	r.taskHistoryStore.Save()
}

func (r *Runner) RunMarketSubTasks(market string, subTasks []types.CrawlerSubTask) {
	for _, subTask := range subTasks {
		err := r.runSubTask(subTask)
		// if has error, skip all subtasks for this market
		if err != nil {
			return
		}
	}
}

func (r *Runner) Run(tasks []types.CrawlerTask) {
	log.Printf("[START] running %v tasks", len(tasks))
	defer r.cleanup()

	// group tasks by market to run in parallel
	marketSubTasks := make(map[string][]types.CrawlerSubTask)
	for _, task := range tasks {
		for _, exterior := range task.Exteriors {
			for _, market := range task.Markets {
				for taskName, taskConfig := range task.TaskConfigs {
					fullItemName := shared.FormatItemName(task.Name, exterior, false)
					marketSubTasks[market] = append(marketSubTasks[market], types.CrawlerSubTask{
						Name:          fullItemName,
						Market:        market,
						TaskName:      taskName,
						TaskConfig:    taskConfig,
						RerunInterval: task.RerunInterval,
					})
				}
			}
		}
	}

	// run subtasks for each market
	for market, subTasks := range marketSubTasks {
		// run in parallel
		go r.RunMarketSubTasks(market, subTasks)
	}

	// enter loop to rerun tasks until no more reruns
	tick := time.NewTicker(1 * time.Minute)
	defer tick.Stop()

	for range tick.C {
		allDone := true
		// check if all subtasks are done
		for _, subTasks := range marketSubTasks {
			for _, subTask := range subTasks {
				subTaskId := getSubTaskId(&subTask)
				count, ok := r.rerunCounts[subTaskId]
				// if subtask is not done or has not reached max reruns, continue
				if !ok || (count < r.maxReruns && r.marketErrors[subTask.Market] == nil) {
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

func getSubTaskId(subtask *types.CrawlerSubTask) string {
	return fmt.Sprintf("%s_%s_%s", subtask.Market, subtask.Name, subtask.TaskName)
}

func (r *Runner) recordSubTaskRun(subTaskId string) {
	r.taskHistoryStore.Set(subTaskId, shared.GetUnixFloat())
	if _, ok := r.rerunCounts[subTaskId]; !ok {
		r.rerunCounts[subTaskId] = 0
	}
	r.rerunCounts[subTaskId]++
	r.saveSecrets()
}

func (r *Runner) runSubTask(subtask types.CrawlerSubTask) error {
	log.Printf("[%s] Running subtask %v", subtask.Market, subtask)

	subTaskId := getSubTaskId(&subtask)

	var exec = func(subtask *types.CrawlerSubTask) error {
		// check if market has error
		if err, ok := r.marketErrors[subtask.Market]; ok {
			log.Printf("Market %s has error: %v, skipping subtask %s", subtask.Market, err, subtask.TaskName)
			return nil
		}

		// check if last task run within the rerun interval
		if r.taskHistoryStore.Get(subTaskId) != nil {
			// json number loaded as float64
			lastRunTime := r.taskHistoryStore.Get(subTaskId).(float64)
			now := shared.GetUnixNow()
			if now-int64(lastRunTime) < subtask.RerunInterval {
				log.Printf("Subtask %s already run within %d seconds, skipping", subtask.TaskName, subtask.RerunInterval)
				return nil
			}
		}

		// run subtask
		err := r.Crawl(subtask.Market, subtask.Name, subtask.TaskName, subtask.TaskConfig)
		return err
	}

	// run
	err := exec(&subtask)
	r.recordSubTaskRun(subTaskId)
	// if error, record error
	if err != nil {
		log.Printf("[%s] Failed to crawl %s for sub task %s: %v", subtask.Market, subtask.TaskName, subtask.Name, err)
		r.marketErrors[subtask.Market] = err
		return err
	}

	// rerun subtask
	if r.rerunCounts[subTaskId] < r.maxReruns {
		waitDuration := time.Duration(subtask.RerunInterval) * time.Second
		log.Printf("Scheduling rerun %d for subtask %s in %v", r.rerunCounts[subTaskId], subTaskId, waitDuration)
		r.taskTimers[subTaskId] = time.AfterFunc(waitDuration, func() {
			r.runSubTask(subtask)
		})
	} else {
		log.Printf("Finished subtask %s after %d runs", subTaskId, r.maxReruns)
	}

	return nil
}

func (r *Runner) Crawl(market string, itemName, taskName string, config types.CrawlTaskConfig) error {
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
