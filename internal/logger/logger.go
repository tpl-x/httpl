package logger

import (
	"github.com/google/wire"
	"github.com/tpl-x/httpl/internal/config"
)

var ProviderSet = wire.NewSet(
	wire.FieldsOf(new(*config.AppConfig), "Log"),
	NewSlogLogger,
)
