package types

type CrawlerTask struct {
	// item name
	Name string `json:"name"`
	// markets to crawl
	Markets []string `json:"markets"`
	// in seconds
	RerunInterval int64                    `json:"rerunInterval"`
	TaskConfigs   map[string]CrawlerConfig `json:"taskConfigs"`
}

type CrawlerSubTask struct {
	Name       string `json:"name"`
	Market     string `json:"market"`
	TaskName   string `json:"taskName"`
	TaskConfig CrawlerConfig

	RerunInterval int64 `json:"rerunInterval"`
}

type CrawlerTasks struct {
	Tasks []CrawlerTask `json:"tasks"`
}
