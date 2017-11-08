package api

type Config struct {
	Name	string
	Params	map[string]interface{}
}

type Event struct {
	Tag		string
	Payload	map[string]interface{}
}

type Input interface {
	Init(c *Config) error
	Start() (<-chan *Event, error)
	Stop() error
}

type Filter interface {
	Init(c *Config) error
	Start(in <-chan *Event) (<-chan *Event, error)
	Stop() error
}

type Output interface {
	Init(c *Config) error
	Start(in <-chan *Event) error
	Stop() error
}

