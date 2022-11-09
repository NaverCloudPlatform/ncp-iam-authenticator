package kubeconfig

type KubectlConfig struct {
	Kind           string                    `yaml:"kind"`
	ApiVersion     string                    `yaml:"apiVersion"`
	CurrentContext string                    `yaml:"current-context"`
	Clusters       []*KubectlClusterWithName `yaml:"clusters"`
	Contexts       []*KubectlContextWithName `yaml:"contexts"`
	Users          []*KubectlUserWithName    `yaml:"users"`
}

type KubectlClusterWithName struct {
	Name    string         `yaml:"name"`
	Cluster KubectlCluster `yaml:"cluster"`
}

type KubectlCluster struct {
	Server                   string `yaml:"server,omitempty"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
}

type KubectlContextWithName struct {
	Name    string         `yaml:"name"`
	Context KubectlContext `yaml:"context"`
}

type KubectlContext struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type KubectlUserWithName struct {
	Name string      `yaml:"name"`
	User KubectlUser `yaml:"user"`
}

type KubectlUser struct {
	ClientCertificateData string     `yaml:"client-certificate-data,omitempty"`
	ClientKeyData         string     `yaml:"client-key-data,omitempty"`
	Password              string     `yaml:"password,omitempty"`
	Username              string     `yaml:"username,omitempty"`
	Token                 string     `yaml:"token,omitempty"`
	Exec                  ExecConfig `yaml:"exec,omitempty"`
}

type ExecConfig struct {
	// Command to execute.
	Command string `yaml:"command"`
	// Arguments to pass to the command when executing it.
	// +optional
	Args []string `yaml:"args"`
	// Env defines additional environment variables to expose to the process. These
	// are unioned with the host's environment, as well as variables client-go uses
	// to pass argument to the plugin.
	// +optional
	Env []ExecEnvVar `yaml:"env"`

	// Preferred input version of the ExecInfo. The returned ExecCredentials MUST use
	// the same encoding version as the input.
	APIVersion string `yaml:"apiVersion,omitempty"`
}

type ExecEnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
