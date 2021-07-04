package config

type Request struct {
	Method string
	URL    string
}

type Assert struct {
	Status     string
	StatusCode int
	// TODO: lt/lte/gt/gte

	StatusCodeIn []int `mapstructure:"statusCode_in"`

	ContentLength int64
	// TODO: lt/lte/gt/gte

	Body string

	BodyContains    string `mapstructure:"body_contains"`
	BodyNotContains string `mapstructure:"body_not_contains"`
	BodyStartsWith  string `mapstructure:"body_startswith"`
	BodyEndsWith    string `mapstructure:"body_endswith"`
}

type Case struct {
	Title       string
	Description string

	Request Request
	Assert  Assert
}
