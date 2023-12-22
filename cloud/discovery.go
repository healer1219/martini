package cloud

import (
	"fmt"
	"github.com/healer1219/martini/global"
	"github.com/healer1219/martini/utils"
	"strconv"
)

type ServiceInstance interface {

	// GetInstanceId return The unique instance ID as registered.
	GetInstanceId() string

	// GetServiceId return The service ID as registered.
	GetServiceId() string

	// GetHost return The hostname of the registered service instance.
	GetHost() string

	// GetPort return The port of the registered service instance.
	GetPort() int

	// IsSecure return Whether the port of the registered service instance uses HTTPS.
	IsSecure() bool

	// GetHealthCheckUrl return The health check URL of the registered service instance.
	GetHealthCheckUrl() string

	// GetMetadata return The key / value pair metadata associated with the service instance.
	GetMetadata() map[string]string
}

type DefaultServiceInstance struct {
	InstanceId     string
	ServiceId      string
	Host           string
	Port           int
	Secure         bool
	HealthCheckUrl string
	Metadata       map[string]string
}

const defaultInstanceIdTemplate = "%s-%s-%s"

func NewDefaultServiceInstance() (*DefaultServiceInstance, error) {
	app := global.Config().App
	cloud := global.Config().Cloud
	ip, err := utils.GetLocalIp()
	if err != nil {
		return nil, err
	}
	return NewServiceInstance(app.Name, ip, app.Port, false, cloud.HealthCheckPath, nil, "")
}

func NewServiceInstance(serviceId string, host string, port int, secure bool, healthCheckUrl string, metadata map[string]string, instanceId string) (*DefaultServiceInstance, error) {
	if len(host) == 0 {
		localIp, err := utils.GetLocalIp()
		if err != nil {
			return nil, err
		}
		host = localIp
	}

	if len(instanceId) == 0 {
		instanceId = fmt.Sprintf(defaultInstanceIdTemplate, serviceId, host, strconv.Itoa(port))
	}

	return &DefaultServiceInstance{
		InstanceId:     instanceId,
		ServiceId:      serviceId,
		Host:           host,
		Port:           port,
		Secure:         secure,
		HealthCheckUrl: healthCheckUrl,
		Metadata:       metadata,
	}, nil
}

func (serviceInstance *DefaultServiceInstance) GetInstanceId() string {
	return serviceInstance.InstanceId
}

func (serviceInstance *DefaultServiceInstance) GetServiceId() string {
	return serviceInstance.ServiceId
}

func (serviceInstance *DefaultServiceInstance) GetHost() string {
	return serviceInstance.Host
}

func (serviceInstance *DefaultServiceInstance) GetPort() int {
	return serviceInstance.Port
}

func (serviceInstance *DefaultServiceInstance) IsSecure() bool {
	return serviceInstance.Secure
}

func (serviceInstance *DefaultServiceInstance) GetMetadata() map[string]string {
	return serviceInstance.Metadata
}

func (serviceInstance *DefaultServiceInstance) GetHealthCheckUrl() string {
	return serviceInstance.HealthCheckUrl
}
