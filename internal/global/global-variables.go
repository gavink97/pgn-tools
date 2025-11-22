package global

import (
	"log/slog"
)

var ProgramLevel = &slog.LevelVar{}
var Logger *slog.Logger

var VERSION = "1.0.3"
var Output = ""

var AllowExperimental = false
