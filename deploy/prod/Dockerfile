ARG GO_IMAGE=docker.m.daocloud.io/library/golang:1.26-alpine
ARG ALPINE_IMAGE=docker.m.daocloud.io/library/alpine:3.21
FROM ${GO_IMAGE} AS builder

WORKDIR /src

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/gateway ./cmd/server

FROM ${ALPINE_IMAGE}

WORKDIR /app

COPY --from=builder /out/gateway ./gateway
COPY conf/config.prod.yaml ./conf/config.yaml

RUN mkdir -p conf logs static/upload

EXPOSE 8079

CMD ["./gateway"]
