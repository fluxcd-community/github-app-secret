package main

import (
	"testing"

	gas "github.com/darkowlzz/github-app-secret"
)

func TestValidateInput(t *testing.T) {
	cases := []struct {
		name    string
		flags   *cmdFlags
		wantErr bool
	}{
		{
			name: "no app ID",
			flags: &cmdFlags{
				installationID: 22,
				secretName:     "foo",
				privateKeyPath: "some-path",
				secretType:     gas.SecretGit,
			},
			wantErr: true,
		},
		{
			name: "no installation ID",
			flags: &cmdFlags{
				appID:          11,
				secretName:     "foo",
				privateKeyPath: "some-path",
				secretType:     gas.SecretGit,
			},
			wantErr: true,
		},
		{
			name: "no secret name",
			flags: &cmdFlags{
				appID:          11,
				installationID: 22,
				privateKeyPath: "some-path",
				secretType:     gas.SecretGit,
			},
			wantErr: true,
		},
		{
			name: "no private key",
			flags: &cmdFlags{
				appID:          11,
				installationID: 22,
				secretName:     "foo",
				secretType:     gas.SecretGit,
			},
			wantErr: true,
		},
		{
			name: "unknown secret type",
			flags: &cmdFlags{
				appID:          11,
				installationID: 22,
				secretName:     "foo",
				privateKeyPath: "some-path",
				secretType:     "bar",
			},
			wantErr: true,
		},
		{
			name: "invalid API URL",
			flags: &cmdFlags{
				apiURL:         "//:foo",
				appID:          11,
				installationID: 22,
				secretName:     "foo",
				privateKeyPath: "some-path",
				secretType:     gas.SecretGit,
			},
			wantErr: true,
		},
		{
			name: "all valid",
			flags: &cmdFlags{
				apiURL:          "https://example.com",
				appID:           11,
				installationID:  22,
				secretName:      "foo",
				privateKeyPath:  "some-path",
				secretType:      gas.SecretGit,
				secretNamespace: "default",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateInput(tc.flags)
			if (err != nil) != tc.wantErr {
				t.Errorf("expected error %t, actual: %v", tc.wantErr, err)
			}
		})
	}
}
