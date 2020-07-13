FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags '-w -s' -o phrase-server server.go
RUN rm -rf .dockerignore server.go

FROM scratch
WORKDIR /src
COPY --from=builder /src .
ENTRYPOINT ["/src/phrase-server"]
