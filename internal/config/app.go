package config

type App struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		Image string `yaml:"image"`
		Port  int32  `yaml:"port"`
		Env   []Env  `yaml:"env"`
		Scale struct {
			Min *int32 `yaml:"min"`
			Max *int32 `yaml:"max"`
		} `yaml:"scale"`
	} `yaml:"spec"`
}

type Env struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
