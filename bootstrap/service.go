package bootstrap

var serviceList []Service
var mainService MainService

func RegisterService(service Service) {
	serviceList = append(serviceList, service)
}

func SetMainService(service MainService) {
	mainService = service
}

type Bootstraper struct{}

func (Bootstraper) RegistServices() {
	for _, service := range serviceList {
		err := service.Setup()
		if err != nil {
			panic(err)
		}
	}
}

func (Bootstraper) Run() {
	mainService.Run()
}

func (Bootstraper) ShutdownServices() {
	for _, service := range serviceList {
		service.Shutdown()
	}
}
