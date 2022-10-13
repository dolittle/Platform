package microservices

import (
	"context"
	"errors"
	"github.com/dolittle/platform-router/config"
	"github.com/rs/zerolog/log"
	"strings"

	coreV1 "k8s.io/api/core/v1"
)

var (
	ErrPodWasNil              = errors.New("pod was nil")
	ErrPodMissingTenant       = errors.New("pod is missing tenant field")
	ErrPodMissingApplication  = errors.New("pod is missing application field")
	ErrPodMissingEnvironment  = errors.New("pod is missing environment field")
	ErrPodMissingMicroservice = errors.New("pod is missing microservice field")
	ErrPodMissingIPAddress    = errors.New("pod is missing IP address")
)

type FieldSpecifier struct {
	Annotations []string
	Labels      []string
}

type FieldSpecifiers struct {
	Tenant       FieldSpecifier
	Application  FieldSpecifier
	Environment  FieldSpecifier
	Microservice FieldSpecifier
}

type Converter struct {
	Config *config.Config
	fields FieldSpecifiers
}

func (c *Converter) WatchConfig(ctx context.Context) {
	for {
		changed := c.Config.Changed()

		newFields := FieldSpecifiers{}
		if err := c.Config.Unmarshal("kubernetes.fields", &newFields); err != nil {
			log.Warn().Err(err).Msg("Failed to unmarshal kubernetes fields")
		} else {
			c.fields = newFields
		}

		select {
		case <-changed:
			continue
		case <-ctx.Done():
			break
		}
	}
}

func (c *Converter) ConvertPodToMicroservice(pod *coreV1.Pod) (Microservice, error) {
	if pod == nil {
		return Microservice{}, ErrPodWasNil
	}

	tenant := c.getPodField(pod, c.fields.Tenant)
	if tenant == "" {
		return Microservice{}, ErrPodMissingTenant
	}

	application := c.getPodField(pod, c.fields.Application)
	if application == "" {
		return Microservice{}, ErrPodMissingApplication
	}

	environment := c.getPodField(pod, c.fields.Environment)
	if environment == "" {
		return Microservice{}, ErrPodMissingEnvironment
	}

	microservice := c.getPodField(pod, c.fields.Microservice)
	if microservice == "" {
		return Microservice{}, ErrPodMissingMicroservice
	}

	ipAddress := pod.Status.PodIP
	if ipAddress == "" {
		return Microservice{}, ErrPodMissingIPAddress
	}

	info := Microservice{
		Identity: ToIdentity(tenant, application, environment, microservice),
		IP:       ipAddress,
		Ports:    make(map[Port]int),
	}

	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			info.Ports[Port{
				Container: container.Name,
				Port:      port.Name,
			}] = int(port.ContainerPort)
		}
	}

	return info, nil

}

func (c *Converter) getPodField(pod *coreV1.Pod, field FieldSpecifier) string {
	for _, annotation := range field.Annotations {
		if value := pod.Annotations[annotation]; value != "" {
			return value
		}
	}
	for _, label := range field.Labels {
		if value := pod.Labels[label]; value != "" {
			return value
		}
	}
	return ""
}

func ToIdentity(tenant, application, environment, microservice string) Identity {
	return Identity{
		Tenant:       strings.ToLower(tenant),
		Application:  strings.ToLower(application),
		Environment:  strings.ToLower(environment),
		Microservice: strings.ToLower(microservice),
	}
}
