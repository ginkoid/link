FROM golang:alpine as builder
COPY . $GOPATH/src/app/
WORKDIR $GOPATH/src/app/
RUN adduser -D -g '' app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /app
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app /app
USER app
EXPOSE 8000
ENTRYPOINT ["/app"]
