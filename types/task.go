package types

type CrawlerTask struct {
	Name    string   `json:"name"`
	Markets []string `json:"markets"`
	// in seconds
	RerunInterval int           `json:"rerun_interval"`
	Config        CrawlerConfig `json:"config"`
}

type CrawlerTasks struct {
	Tasks []CrawlerTask `json:"tasks"`
}
