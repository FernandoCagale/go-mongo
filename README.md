# go-mongo

`Starting MongoDB server`

```sh
$ docker run --name mongo -d -p 27017:27017 mongo
```

```sh
$ export GOPATH=`pwd`
```

```sh
$ go get gopkg.in/mgo.v2
$ go get github.com/labstack/echo
$ go get github.com/stretchr/testify/assert
```

```sh
$ go install main
```

```sh
$ src/main go test *.go
```

```sh
$ bin/main
```