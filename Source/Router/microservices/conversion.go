package microservices

import (
	"errors"
	"strings"

	coreV1 "k8s.io/api/core/v1"
)

var (
	ErrPodWasNil              = errors.New("pod was nil")
	ErrPodMissingTenant       = errors.New("pod is missing tenant-id annotation")
	ErrPodMissingApplication  = errors.New("pod is missing application-id annotation")
	ErrPodMissingEnvironment  = errors.New("pod is missing environment label")
	ErrPodMissingMicroservice = errors.New("pod is missing microservice-id annotation")
	ErrPodMissingIPAddress    = errors.New("pod is missing IP address")
)

func convertPodToMicroservice(pod *coreV1.Pod, config MicroserviceConfiguration) (Microservice, error) {
	if pod == nil {
		return Microservice{}, ErrPodWasNil
	}

	tenant := pod.Annotations[config.TenantIDAnnotation]
	if tenant == "" {
		return Microservice{}, ErrPodMissingTenant
	}

	application := pod.Annotations[config.ApplicationIDAnnotation]
	if application == "" {
		return Microservice{}, ErrPodMissingApplication
	}

	environment := pod.Labels[config.EnvironmentLabel]
	if environment == "" {
		return Microservice{}, ErrPodMissingEnvironment
	}

	microservice := pod.Annotations[config.MicroserviceIDAnnotation]
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

func ToIdentity(tenant, application, environment, microservice string) Identity {
	return Identity{
		Tenant:       strings.ToLower(tenant),
		Application:  strings.ToLower(application),
		Environment:  strings.ToLower(environment),
		Microservice: strings.ToLower(microservice),
	}
}
