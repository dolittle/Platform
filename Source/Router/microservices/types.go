package microservices

import "strings"

type identityString interface {
	TenantID | ApplicationID | Environment | MicroserviceID
}

func identityFromString[T identityString](value string) T {
	return T(strings.ToLower(value))
}

type TenantID string
type ApplicationID string
type Environment string
type MicroserviceID string

type IPAddress string

type Identity struct {
	Tenant       TenantID
	Application  ApplicationID
	Environment  Environment
	Microservice MicroserviceID
}

type Microservice struct {
	Identity Identity
	IP       IPAddress
}
