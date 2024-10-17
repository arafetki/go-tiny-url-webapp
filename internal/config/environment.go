package config

type Environment int

const (
	DEVELOPMENT Environment = iota + 1
	STAGING
	PRODUCTION
)

func (e Environment) String() string {
	return [...]string{"development", "staging", "production"}[e-1]
}
