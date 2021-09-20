#slim alpine golang image
#build stage
ARG GO_VERSION=1.14
ARG ALPINE_VERSION=3.11

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build-stage


#COPY working directory
WORKDIR /app
COPY . /app


# Install git.
# Git is required for fetching the dependencies.
RUN apk --update add --no-cache ca-certificates\
    git \
    curl \
    tzdata \
    gettext \
    && apk add --update bash \
    && CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` GO111MODULE=`on` go build cmd/socketserver/main.go \
    && apk del wget curl git


#production stage
FROM scratch

# Add in certs
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Add the built binary
COPY --from=build-stage /app .


# Run the application
ENTRYPOINT ["./main"]