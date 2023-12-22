package cloud

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/healer1219/martini/config"
	"github.com/healer1219/martini/result"
	"log"
	"strings"
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
	defaultHealthCheckUrl        = "actuator/health"
	healthCheckUrlFmt            = "%s://%s:%d/%s"
	http                         = "http"
	https                        = "https"
	defaultHealthCheckTimeOut    = "5s"
	defaultHealthCheckInterval   = "10s"
	defaultHealthCheckDeregister = "20s"
	hostPortFormat               = "%s:%d"
)

// DefaultHealthCheck default consul health check router
func DefaultHealthCheck(engine *gin.Engine) {
	engine.GET(
		"/actuator/health",
		func(ctx *gin.Context) {
			ctx.JSON(result.SuccessResult("success"))
		},
	)
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

	var healthCheckUrl string
	if serviceInstance.GetHealthCheckUrl() != "" {
		healthCheckUrl, _ = strings.CutPrefix(
			serviceInstance.GetHealthCheckUrl(),
			"/",
		)
	} else {
		healthCheckUrl = defaultHealthCheckUrl
	}

	if isSecure {
		agent.HTTP = fmt.Sprintf(healthCheckUrlFmt, https, registration.Address, registration.Port, healthCheckUrl)
	} else {
		agent.HTTP = fmt.Sprintf(healthCheckUrlFmt, http, registration.Address, registration.Port, healthCheckUrl)
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

func NewDefaultConsulServiceRegistry(registryConf *config.Registry) (*ConsulServiceRegistry, error) {
	if len(registryConf.Ip) < 3 {
		return nil, errors.New("host is illegal")
	}
	if registryConf.Port <= 0 || registryConf.Port > 65535 {
		return nil, errors.New("port is illegal, port should between 1 and 65535")
	}
	consulConfig := api.DefaultConfig()
	consulConfig.Address = fmt.Sprintf(hostPortFormat, registryConf.Ip, registryConf.Port)
	consulConfig.Token = registryConf.Token
	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Println("err creat client", err)
		return nil, err
	}
	return &ConsulServiceRegistry{
		client: client,
	}, nil
}
