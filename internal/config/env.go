package config

type Env string

func (e Env) String() string {
	return string(e)
}

const (
	Local Env = "local"
	Prod  Env = "prod"
)

func (e Env) Local() bool {
	return e == Local
}

func (e Env) Prod() bool {
	return e == Prod
}
