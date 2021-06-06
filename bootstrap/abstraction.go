package bootstrap

type MainService interface {
	Run() error
}

type Service interface {
	Setup() error
	Shutdown()
}
