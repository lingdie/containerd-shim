package image

type RegistryOptions struct {
	RegistryAddr string
	UserName     string
	Password     string
	Repository   string
}

type Registry struct {
	ContainerdNamespace string

	RegistryAddr string
	UserName     string
	Password     string
	Repository   string
	LoginAddress string
}

func NewRegistry(globalRegistry RegistryOptions, envRegistry RegistryOptions, containerNamespace string) *Registry {
	registry := &Registry{
		ContainerdNamespace: containerNamespace,

		RegistryAddr: checkEmpty(envRegistry.RegistryAddr, globalRegistry.RegistryAddr),
		UserName:     checkEmpty(envRegistry.UserName, globalRegistry.UserName),
		// do not expose password
		Password:     checkEmpty(envRegistry.Password, ""),
		Repository:   checkEmpty(envRegistry.Repository, globalRegistry.Repository),
		LoginAddress: checkEmpty(envRegistry.RegistryAddr, globalRegistry.RegistryAddr),
	}
	if registry.LoginAddress == "docker.io" {
		registry.LoginAddress = ""
	}
	return registry
}

func (r Registry) GetImageRef(image string) string {
	if r.Repository == "" {
		return r.RegistryAddr + "/" + r.UserName + "/" + image
	}
	return r.RegistryAddr + "/" + r.Repository + "/" + r.UserName + "/" + image
}

func checkEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
