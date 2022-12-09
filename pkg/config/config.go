package config

type RunConfig struct {
	Debug  bool
	Render bool
	Env    map[string]interface{}

	// exit 1 if got one assert fail
	FailFast bool

	// timeout in ms
	Timeout int

	Order []Order
}

type Order struct {
	Pattern string
	// TODO: not supported yet
	Parallel bool
}

type Retry struct {
	Enable bool
	// retry times
	Count int
	// in ms
	Interval int
	// match the status_codes will do retry
	StatusCodes []int
}

type CaseConfig struct {
	// timeout in ms
	Timeout int

	Retry Retry

	Repeat int
}
