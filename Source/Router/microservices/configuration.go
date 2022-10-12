package microservices

type Configuration struct {
	Microservice MicroserviceConfiguration `koanf:"microservice"`
	Proxy        ProxyConfiguration        `koanf:"proxy"`
}

type MicroserviceConfiguration struct {
	TenantIDAnnotation       string `koanf:"tenant"`
	ApplicationIDAnnotation  string `koanf:"application"`
	MicroserviceIDAnnotation string `koanf:"microservice"`
	EnvironmentLabel         string `koanf:"environment"`
}

type ProxyConfiguration struct {
	TenantIDHeader string     `koanf:"tenant-header"`
	Paths          ProxyPaths `koanf:"paths"`
}

type ProxyPaths map[string]ContainerAndPort

type ContainerAndPort struct {
	Container string `koanf:"container"`
	Port      string `koanf:"port"`
}
