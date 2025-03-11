package env

import (
	e "github.com/codingconcepts/env"
)

type Env struct {
	MongoDbUri     string `env:"URI" required:"true"`
	MongoDbDatabse string `env:"MONGO_DB" default:"syncmeet"`

	Port int `env:"PORT" default:"8080"`
}

func GetEnvOrDie() *Env {
	c := Env{}
	if err := e.Set(&c); err != nil {
		panic(err)
	}
	return &c
}
