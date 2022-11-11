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

type CaseConfig struct {
	// timeout in ms
	Timeout int

	// TODO: retry
	// TODO: repeat
}
