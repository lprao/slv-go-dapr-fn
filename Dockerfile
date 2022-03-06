FROM golang:1.17 as builder

WORKDIR /src/
COPY . /src/

ARG APP_VERSION=v0.0.1

ENV APP_VERSION=$APP_VERSION
ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -tags netgo -ldflags \
    "-w -extldflags '-static' -X main.AppVersion=${APP_VERSION}" \
    -o ./service .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /src/service .
COPY --from=builder /src/components ./components/

ENTRYPOINT ["./service"]
