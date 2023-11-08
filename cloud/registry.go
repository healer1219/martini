package cloud

type ServiceRegistry interface {
	Register(serviceInstance ServiceInstance) bool

	Deregister()
}
