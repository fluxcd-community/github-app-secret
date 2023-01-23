package githubappsecret

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestPopulateSecret(t *testing.T) {
	testToken := "a-token"

	cases := []struct {
		name       string
		secret     *corev1.Secret
		secretType string
		wantData   map[string]string
	}{
		{
			name:       "empty string data, git secret",
			secret:     &corev1.Secret{},
			secretType: SecretGit,
			wantData: map[string]string{
				"username": AccessTokenUsername,
				"password": testToken,
			},
		},
		{
			name:       "plain secret",
			secret:     &corev1.Secret{},
			secretType: SecretPlain,
			wantData: map[string]string{
				"token": testToken,
			},
		},
		{
			name: "overwrite existing data",
			secret: &corev1.Secret{
				StringData: map[string]string{
					"username": "foo",
					"password": "bar",
				},
			},
			secretType: SecretGit,
			wantData: map[string]string{
				"username": AccessTokenUsername,
				"password": testToken,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			populateSecret(tc.secret, tc.secretType, testToken)
			if !reflect.DeepEqual(tc.secret.StringData, tc.wantData) {
				t.Errorf("unexpected secret data:\n\tWNT: %v\n\tGOT: %v", tc.wantData, tc.secret.StringData)
			}
		})
	}
}
