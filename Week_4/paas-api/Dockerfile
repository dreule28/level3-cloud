# first stage
FROM golang:1.25-alpine AS build
WORKDIR /src

# enable tls verification https traffic (go get || go mod download)
RUN apk add --no-cache ca-certificates git

# cache dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# copy everything from local to workdir
COPY paas-api/ .

# build the binary with caching for faster rebuilds
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/paas-api ./cmd/api


# second stage: minimal distroless image for production
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /out/paas-api /paas-api
EXPOSE 8080
ENTRYPOINT ["/paas-api"]