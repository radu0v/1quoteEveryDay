package config

import "github.com/alexedwards/scs/v2"

type AppConfig struct {
	InProduction bool
	Session      *scs.SessionManager
}
