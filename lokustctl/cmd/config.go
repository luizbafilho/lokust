package cmd

type TestResourceRequirements struct {
	Limits   map[string]string
	Requests map[string]string
}

type TestResources struct {
	Master  TestResourceRequirements `json:"master"`
	Workers TestResourceRequirements `json:"workers"`
}

type Config struct {
	Name      string        `json:"name,omitempty"`
	Namespace string        `json:"namespace,omitempty"`
	Replicas  int32         `json:"replicas,omitempty"`
	Resources TestResources `json:"resources,omitempty"`

	File string `json:"file,omitempty"`
}
