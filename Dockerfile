# syntax=docker/dockerfile:1

# TARGETOS/TARGETARCH are populated by Docker Buildx for cross-platform builds
# and left empty (native build) by classic builders such as Google Cloud Build.
# CGO_ENABLED=0 produces a static binary for the target platform.
FROM golang:1.27.1-alpine AS build

RUN apk add --no-cache ca-certificates
WORKDIR /src

COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY web/ ./web/

ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0 GOTOOLCHAIN=local
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -buildvcs=false -ldflags="-s -w" \
    -o /jevmoji ./cmd/jevmoji

FROM scratch

# Go's HTTPS client needs a trusted CA bundle to reach the TypeSafe API.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /jevmoji /jevmoji
COPY LICENSE /LICENSE

LABEL org.opencontainers.image.title="jevmoji" \
      org.opencontainers.image.description="Emoji reactions powered by TypeSafe Jev" \
      org.opencontainers.image.licenses="MIT"

USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/jevmoji"]
CMD ["-addr", ":8080"]
