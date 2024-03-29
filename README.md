# nuts-registry-admin-demo
Demo application which shows how to integrate with the Nuts registry.

## Building and running
### Production
To build for production:

```shell
$ npm install
$ npm run build
$ go run .
```

This will serve the front end from the embedded filesystem.
### Development

During front-end development, you probably want to use the real filesystem and webpack in watch mode:

```shell
$ npm install
$ npm run watch
$ go run . live
```

The API and domain types are generated from the `api/api.yaml`.
```shell
$ oapi-codegen -generate server -package api api/api.yaml > api/generated.go
$ oapi-codegen -generate types -package domain -o domain/generated_types.go api/api.yaml

```

### Docker
```shell
$ docker run -p 1303:1303 nutsfoundation/nuts-registry-admin-demo
```

## Configuration
When running in Docker without a config file mounted at `/app/server.config.yaml` it will use the default configuration.
In this case the default username will be `demo@nuts.nl`. The password is generated and printed in the log on startup.

The `nutsnodeapikeyfile` config parameter should point to a PEM encoded private key file. The corresponding public key should be configured on the Nuts node in SSH authorized keys format.
`nutsnodeapiuser` Is required when using Nuts node API token security. It must match the user in the SSH authorized keys file.

## Technology Stack

Frontend framework is vue.js 3.x

Icons are from https://heroicons.com

CSS framework is https://tailwindcss.com



