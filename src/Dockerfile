## Build
FROM golang:1.19.3-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /http-dump

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /http-dump /htt-dump

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/htt-dump-server"]
