package redis

// Config - Config struct for redis
type config struct {
	Host                  string
	Port                  string
	Password              string
	MaxIdleConnections    int
	MaxActiveConnections  int
	IdleTimeout           int
	ConnectTimeoutSeconds int
	ReadTimeoutSeconds    int
	WriteTimeoutSeconds   int
	UseTLS                *bool
	TLSSkipVerify         *bool
}
