package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	flag "github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	gas "github.com/fluxcd-community/github-app-secret"
)

var logger = log.Log

type cmdFlags struct {
	logLevel        int32
	apiURL          string
	appID           int64
	installationID  int64
	privateKeyPath  string
	secretType      string
	timeout         time.Duration
	secretName      string
	secretNamespace string
}

func main() {
	cf := &cmdFlags{}

	flag.Int32VarP(&cf.logLevel, "logLevel", "v", 0, "Log verbosity level")
	flag.StringVar(&cf.apiURL, "apiURL", "", "Github API URL (default \"https://api.github.com\")")
	flag.Int64Var(&cf.appID, "appID", 0, "Github App ID")
	flag.Int64Var(&cf.installationID, "installationID", 0, "Github App Installation ID")
	flag.StringVar(&cf.privateKeyPath, "privateKeyPath", "", "Path to the Github App private key")
	flag.StringVar(&cf.secretType, "secretType", gas.SecretGit, "Type of secret to create [git, plain]")
	flag.DurationVar(&cf.timeout, "timeout", 15*time.Second, "Timeout for token generation and secret creation")
	flag.StringVar(&cf.secretName, "secretName", "", "Name of the Secret to store the token in")
	flag.StringVar(&cf.secretNamespace, "secretNamespace", "", "Namespace of the Secret to store the token in")

	flag.Parse()

	log.SetLogger(zap.New(zap.Level(zapcore.Level(-cf.logLevel))))

	if err := validateInput(cf); err != nil {
		logger.Error(err, "validation failed")
		os.Exit(1)
	}

	cfg := config.GetConfigOrDie()
	kclient, err := client.New(cfg, client.Options{})
	if err != nil {
		logger.Error(err, "failed to configure kubernetes client")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cf.timeout)
	defer cancel()

	// Configure the secret namespace. Namespace is injected via downward API as
	// an environment variable. If no namespace is provided, use the "default"
	// namespace. Prioritize namespace set via flag.
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}
	if cf.secretNamespace != "" {
		namespace = cf.secretNamespace
	}

	secretKey := client.ObjectKey{
		Name:      cf.secretName,
		Namespace: namespace,
	}
	err = gas.NewAppSecret(kclient, logger, cf.apiURL, cf.privateKeyPath, cf.appID, cf.installationID).
		GenerateAndCreate(ctx, secretKey, cf.secretType)
	if err != nil {
		logger.Error(err, "failed to generate token and create secret")
		os.Exit(1)
	}
	logger.V(2).Info(fmt.Sprintf("token generated and created/updated Secret (%s)", secretKey))
}

func validateInput(cf *cmdFlags) error {
	if cf.apiURL != "" {
		_, err := url.Parse(cf.apiURL)
		if err != nil {
			return fmt.Errorf("invalid API URL: %w", err)
		}
	}

	if cf.appID == 0 {
		return fmt.Errorf("invalid Github App ID: %d", cf.appID)
	}
	if cf.installationID == 0 {
		return fmt.Errorf("invalid Github App Installation ID: %d", cf.installationID)
	}
	if cf.privateKeyPath == "" {
		return errors.New("must provide a privage key path with --privateKeyPath")
	}
	if cf.secretName == "" {
		return errors.New("must provide a Secret name with --secretName")
	}

	switch cf.secretType {
	case gas.SecretPlain, gas.SecretGit:
	default:
		return fmt.Errorf("invalid secret type %q", cf.secretType)
	}
	return nil
}
