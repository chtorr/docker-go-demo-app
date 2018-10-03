# build stage
FROM golang:alpine AS builder
ADD . /go/src/github.com/chtorr/docker-go-demo-app
RUN  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app /go/src/github.com/chtorr/docker-go-demo-app/src/*.go

# final stage
FROM scratch
COPY --from=builder /app /app
CMD ["/app"]