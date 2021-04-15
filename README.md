# nuts-registry-admin-demo
Demo application which shows how to integrate with the Nuts registry.

## Building and running
### Local
To build for production:

```
$ npm install
$ npm run build
$ go run main.go
```

This will serve the front end from the embedded filesystem.

During development you might want to use the real filesystem and webpack in watch mode:

```
$ npm install
$ npm run watch
$ go run main.go live
```

### Docker
```
$ docker run -p 1303:1303 nutsfoundation/nuts-registry-admin-demo
```

## Configuration
You can configure the application by changing the values in `server.config.yaml`.

The Default http port is `1303`.

Credentials to get a session: user:`demo@nuts.nl` password:`demo`.
