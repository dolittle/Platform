FROM golang:latest as build-cli
WORKDIR /build
COPY Source/Backups/Cli/ ./
ENV CGO_ENABLED=0
RUN go build .

FROM scratch
COPY --from=build-cli /build/backups-cli /bin/backups-cli
ENTRYPOINT [ "/bin/backups-cli" ]
