package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"golang.ngrok.com/ngrok/internal/tunnel/proto"
	po "golang.ngrok.com/ngrok/policy"
)

func testPolicy[T tunnelConfigPrivate, O any, OT any](t *testing.T,
	makeOpts func(...OT) Tunnel,
	getPolicies func(*O) any,
) {
	optsFunc := func(opts ...any) Tunnel {
		return makeOpts(assertSlice[OT](opts)...)
	}
	cases := testCases[T, O]{
		{
			name: "absent",
			opts: optsFunc(),
			expectOpts: func(t *testing.T, opts *O) {
				actual := getPolicies(opts)
				require.Nil(t, actual)
			},
		},
		{
			name: "with policy",
			opts: optsFunc(
				WithPolicy(
					po.Policy{
						Inbound: []po.Rule{
							{
								Name:        "denyPUT",
								Expressions: []string{"req.Method == 'PUT'"},
								Actions: []po.Action{
									{Type: "deny"},
								},
							},
							{
								Name:        "logFooHeader",
								Expressions: []string{"'foo' in req.Headers"},
								Actions: []po.Action{
									{
										Type:   "log",
										Config: map[string]any{"metadata": map[string]any{"key": "val"}},
									},
								},
							},
						},
						Outbound: []po.Rule{
							{
								Name: "InternalErrorWhenFailed",
								Expressions: []string{
									"res.StatusCode <= '0'",
									"res.StatusCode >= '300'",
								},
								Actions: []po.Action{
									{
										Type:   "custom-response",
										Config: map[string]any{"status_code": 500},
									},
								},
							},
						},
					},
				),
			),
			expectOpts: func(t *testing.T, opts *O) {
				actualAny := getPolicies(opts)
				require.NotNil(t, actualAny)

				actual := actualAny.(po.Policy)

				require.Len(t, actual.Inbound, 2)
				require.Equal(t, "denyPUT", actual.Inbound[0].Name)
				require.Equal(t, actual.Inbound[0].Actions, []po.Action{{Type: "deny"}})
				require.Len(t, actual.Outbound, 1)
				require.Len(t, actual.Outbound[0].Expressions, 2)
			},
		},
		{
			name: "with policy string",
			opts: optsFunc(
				WithPolicyString(`
					{
						"inbound":[
							{
								"name":"denyPut",
								"expressions":["req.Method == 'PUT'"],
								"actions":[{"type":"deny"}]
							},
							{
								"name":"logFooHeader",
								"expressions":["'foo' in req.Headers"],
								"actions":[
									{"type":"log","config":{"metadata":{"key":"val"}}}
								]
							}
						],
						"outbound":[
							{
								"name":"500ForFailures",
								"expressions":["res.StatusCode <= 0", "res.StatusCode >= 300"],
								"actions":[{"type":"custom-response", "config":{"status_code":500}}]
							}
						]
					}`)),
			expectOpts: func(t *testing.T, opts *O) {
				actualAny := getPolicies(opts)
				actualSer, err := json.Marshal(actualAny)
				require.NoError(t, err)

				var actual po.Policy
				require.NoError(t, json.Unmarshal(actualSer, &actual))

				require.NotNil(t, actual)
				require.Len(t, actual.Inbound, 2)
				require.Equal(t, "denyPut", actual.Inbound[0].Name)
				require.Equal(t, []po.Action{{Type: "deny"}}, actual.Inbound[0].Actions)
				require.Len(t, actual.Outbound, 1)
				require.Len(t, actual.Outbound[0].Expressions, 2)
				require.Equal(t, map[string]any{"status_code": 500.}, actual.Outbound[0].Actions[0].Config)
			},
		},
	}

	cases.runAll(t)
}

func TestPolicy(t *testing.T) {
	testPolicy[*httpOptions](t, HTTPEndpoint,
		func(h *proto.HTTPEndpoint) any {
			return h.Policy
		})
	testPolicy[*tcpOptions](t, TCPEndpoint,
		func(h *proto.TCPEndpoint) any {
			return h.Policy
		})
	testPolicy[*tlsOptions](t, TLSEndpoint,
		func(h *proto.TLSEndpoint) any {
			return h.Policy
		})
}
