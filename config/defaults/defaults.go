package defaults

import (
	"time"
)

// server constants
const (
	Address         = "localhost"
	Port            = 8080
	ReadTimeout     = time.Second * 5
	WriteTimeout    = time.Second * 5
	ShutdownTimeout = time.Second * 30
	IdleTimeout     = time.Second * 60
)
