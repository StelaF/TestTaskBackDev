package userService

type User struct {
	cfg config
}

type config struct {
	UserServiceUrl string `env:"USER_SERVICE_URL"`
}
