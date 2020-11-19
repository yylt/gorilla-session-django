package gorilla_session_django

type Gsd_config struct {
	// session key in cookie, mostly is 'sessionId'
	CookieKey string

	//auth function to impl auth
	Auth func(interface{}) error

	// json serializer,use pickle if false
	JsonSerializer bool

	// prefix the key when operate memcache will as prefixkey
	PrefixMemcache string
}
