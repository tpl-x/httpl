package logger

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewSlogLogger,
)