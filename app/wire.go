//+build wireinject

package app

import (
	"github.com/google/wire"
)

func InitApp(config Config) *App {
	wire.Build(set)
	return &App{}
}
