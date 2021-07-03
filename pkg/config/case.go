package config


type Request struct {
	Method string
	URL string
}

type Assert struct {
	Status string
	StatusCode int
}

type Case struct {
	Title string
	Description string

	Request Request
	Assert Assert
}

