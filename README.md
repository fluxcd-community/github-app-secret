# github-app-secret

Generate Github app auth token and write it into a Kubernetes Secret and refresh
it periodically.

The application `./cmd/github-app-secret` takes Github app private key, app ID,
installation ID and a secret name, and generates an auth token and writes it to
a Kubernetes Secret with the given secret name. This can be used by any
application that needs Github app based authentication.

## Instructions

⚠️ __WARNING:__ Please make sure that the system time where this program runs is
up-to-date. The token generation requests contain expiry time. If the expiry
time used in the request is in the past, token generation would fail with **401
Unauthorized** error.

Create a new Github app with the appropriate permissions, generate a private key
for the app and install the app in the target repositories. [Refer the official
docs](https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#accessing-api-endpoints-as-a-github-app)
for detailed instructions.

The app ID can be obtained from the app settings page at
`https://github.com/settings/apps/<app-name>`.

The installation ID can be obtained from
`https://github.com/settings/installations` page. On clicking an installed app,
the URL will contain the installation ID
`https://github.com/settings/installations/<installation-id>`. For
organizations, the first part of the URL may be different, but it follows the
same pattern.

Put the private key in a Kubernetes Secret with

```shell
$ kubectl create secret generic github-app-private-key --from-file=privatekey.pem=/path-to-private-key.pem
```

This secret will be mounted as a volume and used by `github-app-secret`.

`github-app-secret` is run as a Kubernetes CronJob. Modify the manifests from
`./deploy` directory, adding the parameters collected above as argument to the
`github-app-secret` container. For example:

```yaml
    ...
    containers:
        - name: github-app-secret
          args:
            - "-v=3"
            - --privateKeyPath=/etc/secret-volume/privatekey.pem
            - --appID=<app-id>
            - --installationID=<installation-id>
            - --secretName=<secret-name>
    ...
```

Update the CronJob schedule depending on the needs, ensuring that the token gets
refreshed before expiry.

Make sure that the manifests in `./deploy/rbac.yaml`, which provide
`github-app-secret` the necessary permissions it needs to create and update the
Secret, are applied along with the CronJob manifest.

For cloning git repositories, the secret of type `git` can be used. This is the
default type of Secret. It creates secret data with `username` field
`x-access-token` as [required by Github for http based clone](https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#http-based-git-access-by-an-installation).

For just the auth token, the secret of type `plain` can be used. This can be
configured in `github-app-secret` by using `--secretType` flag.

For Github Enterprise, the Github API URL can be configured with `--apiURL`
flag.

## Build

Since this is a very basic golang application, [`ko`](https://ko.build/) can be
used to build a container image for it.

Install `ko` and run `make ko-build` to build a container image for it. This
will build the image and load it in the local container image store.

In order to build and publish to a remote repository, run
`KO_DOCKER_REPO=<container-repo-address> make ko-publish`. Refer
https://ko.build/get-started/#choose-destination for more examples.
