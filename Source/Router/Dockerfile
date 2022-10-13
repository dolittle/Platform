FROM golang:1.19.2 AS build

COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -o /build/platform-proxy


FROM gcr.io/distroless/base

COPY --from=build /build/platform-proxy /bin/platform-proxy
COPY config.yml /var/lib/platform-proxy/config.yml
WORKDIR /var/lib/platform-proxy
ENTRYPOINT [ "/bin/platform-proxy" ]
CMD [ "proxy", "--config=config.yml" ]
EXPOSE 8080
