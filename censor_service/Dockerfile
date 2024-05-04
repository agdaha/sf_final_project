FROM golang:1.21.9-alpine3.19 AS builder

WORKDIR /usr/local/go/src/

ADD . /usr/local/go/src/

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o app cmd/main/main.go

FROM alpine:3.19

COPY --from=builder /usr/local/go/src/app /
COPY --from=builder /usr/local/go/src/configs/config.yaml /
ENV CONFIG_PATH=/config.yaml
CMD ["/app"]