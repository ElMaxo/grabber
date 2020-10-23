FROM golang:1.13 as build-env

# init project dir and download dependencies
RUN mkdir /grabber
WORKDIR /grabber

# download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

RUN GO111MODULE=on go get -u github.com/go-swagger/go-swagger/cmd/swagger@v0.25.0

# copy sources and generate code
COPY . .
RUN swagger generate client --target=internal/rest -f api/grabber.swagger.yml && \
    swagger generate server --target=internal/rest -f api/grabber.swagger.yml

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /go/bin/grabber /grabber/cmd/grabber

FROM alpine:3.5
RUN apk --no-cache add ca-certificates
COPY --from=build-env /go/bin/grabber /go/bin/grabber

CMD ["/go/bin/grabber"]
