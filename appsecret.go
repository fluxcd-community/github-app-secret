package githubappsecret

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	SecretGit   string = "git"
	SecretPlain string = "plain"

	AccessTokenUsername = "x-access-token"
)

// AppSecret helps generates Github app auth token and save it in a Kubernetes
// Secret.
type AppSecret struct {
	client.Client
	log logr.Logger

	apiURL         string
	privateKey     string
	appID          int64
	installationID int64
}

// NewAppSecret constructs and returns a new AppSecret instance.
func NewAppSecret(kclient client.Client, log logr.Logger, apiURL, privateKey string, appID, installationID int64) *AppSecret {
	return &AppSecret{
		Client:         kclient,
		log:            log,
		apiURL:         apiURL,
		privateKey:     privateKey,
		appID:          appID,
		installationID: installationID,
	}
}

// GenerateAndCreate generates an auth token and creates a secret to store the
// token in Kubernetes based on the configured parameters.
func (as *AppSecret) GenerateAndCreate(ctx context.Context, namespacedName client.ObjectKey, secretType string) error {
	token, err := as.GenerateToken(ctx)
	if err != nil {
		return err
	}
	return as.CreateOrUpdateSecret(ctx, namespacedName, secretType, token)
}

// GenerateToken generates an auth token based on the configured parameters.
func (as *AppSecret) GenerateToken(ctx context.Context) (string, error) {
	tr := http.DefaultTransport

	itr, err := ghinstallation.NewKeyFromFile(tr, as.appID, as.installationID, as.privateKey)
	if err != nil {
		return "", err
	}
	if as.apiURL != "" {
		itr.BaseURL = as.apiURL
	}
	return itr.Token(ctx)
}

// CreateOrUpdateSecret creates a new secret or updates an existing secret with
// the new secret data.
func (as *AppSecret) CreateOrUpdateSecret(ctx context.Context, namespacedName client.ObjectKey, secretType, token string) error {
	secret := &corev1.Secret{}
	secret.Name = namespacedName.Name
	secret.Namespace = namespacedName.Namespace

	_, err := controllerutil.CreateOrPatch(ctx, as.Client, secret, func() error {
		populateSecret(secret, secretType, token)
		return nil
	})
	return err
}

func populateSecret(secret *corev1.Secret, secretType, token string) {
	if secret.StringData == nil {
		secret.StringData = map[string]string{}
	}

	switch secretType {
	case SecretGit:
		secret.StringData["username"] = AccessTokenUsername
		secret.StringData["password"] = token
	case SecretPlain:
		secret.StringData["token"] = token
	}
}
