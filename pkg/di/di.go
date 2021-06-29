package di

import (
	"github.com/golobby/container"
	"github.com/profiralex/go-bootstrap-redis/pkg/config"
)

func RegisterDependencies() {
	//Config
	container.Singleton(func() config.Config {
		return config.GetConfig()
	})
}

func UnregisterDependencies() {
	container.Reset()
}

func Make(receiver interface{}) {
	container.Make(receiver)
}

func MakeAll(receivers ...interface{}) {
	for _, receiver := range receivers {
		Make(receiver)
	}
}
