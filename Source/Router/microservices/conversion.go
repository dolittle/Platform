package microservices

import (
	"errors"
	coreV1 "k8s.io/api/core/v1"
	"strings"
)

var (
	ErrPodWasNil              = errors.New("pod was nil")
	ErrPodMissingTenant       = errors.New("pod is missing tenant-id annotation")
	ErrPodMissingApplication  = errors.New("pod is missing application-id annotation")
	ErrPodMissingEnvironment  = errors.New("pod is missing environment label")
	ErrPodMissingMicroservice = errors.New("pod is missing microservice-id annotation")
	ErrPodMissingIPAddress    = errors.New("pod is missing IP address")
)

func convertPodToMicroservice(pod *coreV1.Pod) (Microservice, error) {
	if pod == nil {
		return Microservice{}, ErrPodWasNil
	}

	tenant := pod.Annotations["dolittle.io/tenant-id"]
	if tenant == "" {
		return Microservice{}, ErrPodMissingTenant
	}

	application := pod.Annotations["dolittle.io/tenant-id"]
	if application == "" {
		return Microservice{}, ErrPodMissingApplication
	}

	environment := pod.Labels["environment"]
	if environment == "" {
		return Microservice{}, ErrPodMissingEnvironment
	}

	microservice := pod.Annotations["dolittle.io/microservice-id"]
	if microservice == "" {
		return Microservice{}, ErrPodMissingMicroservice
	}

	ipAddress := pod.Status.PodIP
	if ipAddress == "" {
		return Microservice{}, ErrPodMissingIPAddress
	}

	return Microservice{
		Identity: Identity{
			Tenant:       TenantID(strings.ToLower(tenant)),
			Application:  ApplicationID(strings.ToLower(application)),
			Environment:  Environment(strings.ToLower(environment)),
			Microservice: MicroserviceID(strings.ToLower(microservice)),
		},
		IP: IPAddress(ipAddress),
	}, nil
}
