package app

type Inventory struct {
	Runners []Runner
}

type Runner struct {
	Name            string
	Repo            string
	OS              string
	Host            string
	Path            string
	Transport       string
	LocalState      string
	ServiceName     string
	ControlMode     string
	GitHubStatus    string
	Busy            bool
	Labels          []string
	Version         string
	QueueCount      int
	StaleQueueCount int
}

type runnerConfig struct {
	AgentName  string `json:"agentName"`
	GitHubURL  string `json:"gitHubUrl"`
	WorkFolder string `json:"workFolder"`
}
