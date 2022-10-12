package microservices

import (
	"errors"

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
	microservice := Microservice{}
	if pod == nil {
		return microservice, ErrPodWasNil
	}
	tenantID := pod.Annotations[config.TenantIDAnnotation]
	if tenantID == "" {
		return microservice, ErrPodMissingTenant
	}

	microservice.Identity.Tenant = identityFromString[TenantID](tenantID)
	applicationID := pod.Annotations[config.ApplicationIDAnnotation]
	if applicationID == "" {
		return microservice, ErrPodMissingApplication
	}

	microservice.Identity.Application = identityFromString[ApplicationID](applicationID)
	environment := pod.Labels[config.EnvironmentLabel]
	if environment == "" {
		return microservice, ErrPodMissingEnvironment
	}

	microservice.Identity.Environment = identityFromString[Environment](environment)
	microserviceID := pod.Annotations[config.MicroserviceIDAnnotation]
	if microserviceID == "" {
		return microservice, ErrPodMissingMicroservice
	}

	microservice.Identity.Microservice = identityFromString[MicroserviceID](microserviceID)
	ipAddress := pod.Status.PodIP
	if ipAddress == "" {
		return microservice, ErrPodMissingIPAddress
	}
	microservice.IP = IPAddress(ipAddress)

	microservice.Ports = make(map[Port]int32)
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			microservice.Ports[Port{
				Container: container.Name,
				Port:      port.Name,
			}] = port.ContainerPort
		}
	}

	return microservice, nil
}
