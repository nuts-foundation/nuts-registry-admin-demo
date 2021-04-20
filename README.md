# nuts-registry-admin-demo
Demo application which shows how to integrate with the Nuts registry.

## Building and running
### Production
To build for production:

```shell
$ npm install
$ npm run build
$ go run main.go
```

This will serve the front end from the embedded filesystem.
### Development

During front-end development, you probably want to use the real filesystem and webpack in watch mode:

```shell
$ npm install
$ npm run watch
$ go run . live
```

The API is generated from the `api/api.yaml`.
```shell
$ oapi-codegen -generate types,server -package api api/api.yaml > api/generated.go
```

### Docker
```shell
$ docker run -p 1303:1303 nutsfoundation/nuts-registry-admin-demo
```

## Configuration
You can configure the application by changing the values in `server.config.yaml`.

The Default http port is `1303`.

Credentials to get a session: user:`demo@nuts.nl` password:`demo`.
