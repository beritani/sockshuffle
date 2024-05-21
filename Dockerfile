ARG GO_VERSION=1.22

# Build
FROM docker.io/golang:${GO_VERSION} AS build

WORKDIR /src
COPY src/ .
RUN CGO_ENABLED=0 go build -o /sockshuffle

# Final
FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=build --chown=nonroot:nonroot /sockshuffle /sockshuffle
USER nonroot:nonroot
ENTRYPOINT [ "/sockshuffle" ]

LABEL org.opencontainers.image.title="sockshuffle"
LABEL org.opencontainers.image.description="A lightweight SOCKS5 proxy load balancer"
LABEL org.opencontainers.image.licenses="MIT"