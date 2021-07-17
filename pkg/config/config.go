package config

type RunConfig struct {
	Debug  bool
	Render bool
	Env    map[string]interface{}
}
