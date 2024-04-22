package types

type CrawlerTask struct {
	// item name (base name without exterior)
	Name string `json:"name"`
	// exteriors to crawl (e.g. Factory New, Minimal Wear, etc.)
	Exteriors []string `json:"exteriors"`
	// markets to crawl
	Markets []string `json:"markets"`
	// in seconds
	RerunInterval int64 `json:"rerunInterval"`
	// all tasks to run for this item, each task has its own config
	TaskConfigs map[string]CrawlTaskConfig `json:"taskConfigs"`
}

type CrawlerSubTask struct {
	// full item name (with exterior)
	Name       string `json:"name"`
	Market     string `json:"market"`
	TaskName   string `json:"taskName"`
	TaskConfig CrawlTaskConfig

	RerunInterval int64 `json:"rerunInterval"`
}

type CrawlerTasks struct {
	Tasks []CrawlerTask `json:"tasks"`
}
