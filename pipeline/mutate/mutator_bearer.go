package mutate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ory/oathkeeper/driver/configuration"
	"github.com/ory/oathkeeper/pipeline"
	"github.com/ory/oathkeeper/pipeline/authn"
)

/*

mutators:
  bearer:
    enabled: true

---

mutators:
  - handler: bearer
	config:
	  token_from:
		file: "/var/lib/kubernetes.io/serviceaccount/token"
		environment_variable: "TOKEN"
*/

type MutatorBearerTokenFromConfig struct {
	File                string `json:"file"`
	EnvironmentVariable string `json:"environment_variable"`
}

type MutatorBearerConfig struct {
	TokenFrom MutatorBearerTokenFromConfig `json:"token_from"`
}

type MutatorBearer struct {
	c configuration.Provider
}

func NewMutatorBearer(c configuration.Provider) *MutatorBearer {
	return &MutatorBearer{c: c}
}

func (a *MutatorBearer) GetID() string {
	return "bearer"
}

func (a *MutatorBearer) Mutate(r *http.Request, session *authn.AuthenticationSession, config json.RawMessage, rl pipeline.Rule) error {

	cfg, err := a.config(config)
	if err != nil {
		return err
	}

	token := ""
	if cfg.TokenFrom.EnvironmentVariable != "" {
		token = os.Getenv(cfg.TokenFrom.EnvironmentVariable)
	} else if cfg.TokenFrom.File != "" {
		data, err := ioutil.ReadFile(cfg.TokenFrom.File)
		if err != nil {
			return err
		}
		token = string(data)
	}

	session.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))

	return nil
}

func (a *MutatorBearer) Validate(config json.RawMessage) error {
	if !a.c.MutatorIsEnabled(a.GetID()) {
		return NewErrMutatorNotEnabled(a)
	}

	mutatorConfig, err := a.config(config)
	if err != nil {
		return err
	}

	if mutatorConfig.TokenFrom.File != "" && mutatorConfig.TokenFrom.EnvironmentVariable != "" {
		return NewErrMutatorMisconfigured(a, errors.New("multiple token sources configured"))
	} else if mutatorConfig.TokenFrom.File == "" && mutatorConfig.TokenFrom.EnvironmentVariable == "" {
		return NewErrMutatorMisconfigured(a, errors.New("no token sources configured"))
	}

	return nil
}

func (a *MutatorBearer) config(config json.RawMessage) (*MutatorBearerConfig, error) {
	var c MutatorBearerConfig
	if err := a.c.MutatorConfig(a.GetID(), config, &c); err != nil {
		return nil, NewErrMutatorMisconfigured(a, err)
	}

	return &c, nil
}
