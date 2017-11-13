# go-mongo

`Starting MongoDB server`

```sh
$ docker run --name mongo -d -p 27017:27017 mongo
```

`Install dep`

```sh
$ go get -v github.com/golang/dep/cmd/dep
```

`Install the project's dependencies`

```sh
$ dep ensure
```

`Build docker API`

```sh
$ CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/api
```

`Start API`

```sh
$ go run main.go
```

```sh
$ ./build/api
```