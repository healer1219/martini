package cloud

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/healer1219/martini/result"
	"log"
)

type ConsulServiceRegistry struct {
	serviceInstances     map[string]map[string]ServiceInstance
	client               *api.Client
	localServiceInstance ServiceInstance
}

const (
	needSecure                   = "secure=true"
	unSecure                     = "secure=false"
	tagFormat                    = "%s=%s"
	healthCheckUrl               = "%s://%s:%d/actuator/health"
	http                         = "http"
	https                        = "https"
	defaultHealthCheckTimeOut    = "5s"
	defaultHealthCheckInterval   = "10s"
	defaultHealthCheckDeregister = "20s"
	hostPortFormat               = "%s:%d"
)

// DefaultHealthCheck default consul health check router
func DefaultHealthCheck(ctx *gin.Context) {
	ctx.JSON(result.SuccessResult("success"))
}

func (c *ConsulServiceRegistry) Register(serviceInstance ServiceInstance) bool {
	registration := new(api.AgentServiceRegistration)
	registration.ID = serviceInstance.GetInstanceId()
	registration.Name = serviceInstance.GetServiceId()
	registration.Address = serviceInstance.GetHost()
	registration.Port = serviceInstance.GetPort()

	tags := make([]string, 0)
	isSecure := serviceInstance.IsSecure()
	if isSecure {
		tags = append(tags, needSecure)
	} else {
		tags = append(tags, unSecure)
	}
	metadata := serviceInstance.GetMetadata()
	if metadata != nil {
		for k, v := range metadata {
			tags = append(tags, fmt.Sprintf(tagFormat, k, v))
		}
	}
	registration.Tags = tags

	agent := new(api.AgentServiceCheck)

	if isSecure {
		agent.HTTP = fmt.Sprintf(healthCheckUrl, https, registration.Address, registration.Port)
	} else {
		agent.HTTP = fmt.Sprintf(healthCheckUrl, http, registration.Address, registration.Port)
	}

	agent.Timeout = defaultHealthCheckTimeOut
	agent.Interval = defaultHealthCheckInterval
	// todo 暂时不下线
	//agent.DeregisterCriticalServiceAfter = defaultHealthCheckDeregister

	registration.Check = agent

	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Println("register to consul failed! ", err)
		return false
	}

	if c.serviceInstances == nil {
		c.serviceInstances = map[string]map[string]ServiceInstance{}
	}

	services := c.serviceInstances[serviceInstance.GetServiceId()]

	if services == nil {
		services = map[string]ServiceInstance{}
	}

	services[serviceInstance.GetInstanceId()] = serviceInstance
	c.serviceInstances[serviceInstance.GetServiceId()] = services
	c.localServiceInstance = serviceInstance

	return true
}

// Deregister 从consul中移除服务
func (c *ConsulServiceRegistry) Deregister() {
	if c.serviceInstances == nil {
		return
	}
	localServiceId := c.localServiceInstance.GetServiceId()
	services := c.serviceInstances[localServiceId]
	if services == nil {
		return
	}
	localInstanceId := c.localServiceInstance.GetInstanceId()
	delete(services, localInstanceId)
	if len(services) == 0 {
		delete(c.serviceInstances, localServiceId)
	}
	err := c.client.Agent().ServiceDeregister(localInstanceId)
	if err != nil {
		log.Println("deregister service failed!", err)
	}
	c.localServiceInstance = nil
}

func NewDefaultConsulServiceRegistry(host string, port int, token string) (*ConsulServiceRegistry, error) {
	if len(host) < 3 {
		return nil, errors.New("host is illegal")
	}
	if port <= 0 || port > 65535 {
		return nil, errors.New("port is illegal, port should between 1 and 65535")
	}
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf(hostPortFormat, host, port)
	config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		log.Println("err creat client", err)
		return nil, err
	}
	return &ConsulServiceRegistry{
		client: client,
	}, nil
}
