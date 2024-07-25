package image

const (
	DefaultName      = "docker.io"
	DefaultNamespace = "k8s.io"
)

var (
	DefaultUserName = ""
	DefaultPassword = ""
)

type Registry struct {
	Name         string
	Namespace    string
	UserName     string
	Password     string
	Repository   string
	LoginAddress string
}

func NewRegistry(name, repo, userName, password string) *Registry {
	registry := &Registry{
		Name:         checkEmpty(name, DefaultName),
		Namespace:    DefaultNamespace,
		UserName:     checkEmpty(userName, DefaultUserName),
		Password:     checkEmpty(password, DefaultPassword),
		Repository:   repo,
		LoginAddress: checkEmpty(name, DefaultName),
	}
	if registry.LoginAddress == "docker.io" {
		registry.LoginAddress = ""
	}
	return registry
}

func (r Registry) GetImageRef(image string) string {
	if r.Repository == "" {
		return r.Name + "/" + r.UserName + "/" + image
	}
	return r.Name + "/" + r.Repository + "/" + r.UserName + "/" + image
}

func checkEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
