FROM golang:latest as build-cli
WORKDIR /build
COPY Source/Backups/Cli/ ./
ENV CGO_ENABLED=0
RUN go build .

FROM alpine:latest as ssl-certs

FROM scratch
COPY --from=ssl-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-cli /build/backups-cli /bin/backups-cli
ENTRYPOINT [ "/bin/backups-cli" ]
