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

type CrawlerTasks struct {
	Tasks []CrawlerTask `json:"tasks"`
}
