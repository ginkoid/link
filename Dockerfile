FROM golang:alpine as builder
COPY . $GOPATH/src/github.com/ginkoid/link/
WORKDIR $GOPATH/src/github.com/ginkoid/link/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/app
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/app /go/bin/app
EXPOSE 80
ENTRYPOINT ["/go/bin/app"]
