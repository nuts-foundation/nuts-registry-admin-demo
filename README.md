# nuts-registry-admin-demo
Demo application which shows how to integrate with the Nuts registry.

To build for production:

```
$ npm install
$ npm run build
$ go run main.go
```

This will will serve the front end from the embedded filesystem.

During development you might want to use the real filesystem and webpack in watch mode:

```
$ npm install
$ npm run watch
$ go run main.go live
```
