package microservices

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
