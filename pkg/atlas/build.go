package atlas

import (
	"get.porter.sh/porter/pkg/exec/builder"
	yaml "gopkg.in/yaml.v2"
	"fmt"
)

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the atlas mixin in porter.yaml
// mixins:
// - atlas:
//	  clientVersion: "v0.0.0"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`
}

// This is an example. Replace the following with whatever steps are needed to
// install required components into
const dockerfileLines = `RUN apt-get update && \
apt-get -y install gnupg wget && \
wget -qO - https://pgp.mongodb.com/server-5.0.asc | apt-key add - && \
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/5.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-5.0.list && \
apt-get update && \
apt-get install -y mongodb-atlas-cli
`

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build() error {

	// Create new Builder.
	var input BuildInput

	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	suppliedClientVersion := input.Config.ClientVersion

	if suppliedClientVersion != "" {
		m.ClientVersion = suppliedClientVersion
	}

	fmt.Fprintf(m.Out, dockerfileLines)

	// Example of pulling and defining a client version for your mixin
	// fmt.Fprintf(m.Out, "\nRUN curl https://get.helm.sh/helm-%s-linux-amd64.tar.gz --output helm3.tar.gz", m.ClientVersion)

	return nil
}
