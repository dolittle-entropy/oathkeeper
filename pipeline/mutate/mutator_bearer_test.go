package mutate_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ory/oathkeeper/driver/configuration"
	"github.com/ory/oathkeeper/internal"
	"github.com/ory/oathkeeper/pipeline"
	"github.com/ory/oathkeeper/pipeline/authn"
	"github.com/ory/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBearerMutator(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistry(conf)

	bearerMutator, err := reg.PipelineMutator("bearer")

	require.NoError(t, err)
	assert.Equal(t, "bearer", bearerMutator.GetID())

	t.Run("method=validate", func(t *testing.T) {
		type TestCase struct {
			enabled    bool
			json       string
			shouldPass bool
		}

		for k, testCase := range []TestCase{
			{enabled: false, shouldPass: false},
			{enabled: false, shouldPass: false, json: `{"token_from": {"file": "", "environment_variable": ""}}`},
			{enabled: true, shouldPass: false, json: `{"token_from": {"file": "/some/file", "environment_variable": "SOME_VAR"}}`},
			{enabled: true, shouldPass: false, json: `{"token_from": {"file": "", "environment_variable": ""}}`},
			{enabled: true, shouldPass: false, json: `{"token_from": {"file": ""}}`},
			{enabled: true, shouldPass: false, json: `{"token_from": {}}`},
			{enabled: true, shouldPass: false, json: "{}"},
			{enabled: true, shouldPass: true, json: `{"token_from": {"file": "/some/file"}}`},
			{enabled: true, shouldPass: true, json: `{"token_from": {"environment_variable": "SOME_VAR"}}`},
		} {
			t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
				viper.Reset()
				viper.Set(configuration.ViperKeyMutatorBearerIsEnabled, testCase.enabled)
				err := bearerMutator.Validate(json.RawMessage(testCase.json))
				if testCase.shouldPass {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
				}
			})
		}
	})
	t.Run("method=mutate", func(t *testing.T) {
		type TestCase struct {
			name       string
			shouldPass bool
			config     string
			token      string
		}
		for _, testCase := range []TestCase{
			{name: "Invalid config", shouldPass: false, config: ""},
			{name: "File does not exist", shouldPass: false, config: `{"token_from": {"file": "/some/file/that/should/not/exist.noop"}}`},
			{name: "From environment variable", shouldPass: true, config: `{"token_from": {"environment_variable": "SOME_VAR"}}`, token: "dfa90b04-fa2b-41be-844b-eab200396283"},
		} {
			os.Setenv("SOME_VAR", testCase.token)
			t.Run(fmt.Sprintf("case=%s", testCase.name), func(t *testing.T) {
				viper.Reset()
				viper.Set(configuration.ViperKeyMutatorBearerIsEnabled, true)
				var request http.Request
				var session authn.AuthenticationSession
				var rule pipeline.Rule
				err := bearerMutator.Mutate(&request, &session, json.RawMessage(testCase.config), rule)

				if testCase.shouldPass {
					require.NoError(t, err)
					assert.Equal(t, fmt.Sprintf("Bearer %s", testCase.token), session.Header.Get("Authorization"))
				} else {
					require.Error(t, err)
				}

			})
		}
	})
}
