package microservices

type Identity struct {
	Tenant       string
	Application  string
	Environment  string
	Microservice string
}

type Port struct {
	Container string
	Port      string
}

type Ports map[Port]int

type Microservice struct {
	Identity Identity
	IP       string
	Ports    Ports
}
