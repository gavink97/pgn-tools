package global

import (
	"log/slog"
)

var ProgramLevel = &slog.LevelVar{}
var Logger *slog.Logger
var Args []string

var VERSION = 1.0
