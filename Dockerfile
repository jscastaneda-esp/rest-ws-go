ARG GO_VERSION=1.19.5

FROM golang:${GO_VERSION}-alpine AS builder


RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /src

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
  -installsuffix 'static' \
  -o /rest-ws-go

FROM scratch AS runner

WORKDIR /src

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY .env ./
COPY --from=builder /rest-ws-go /rest-ws-go

EXPOSE 5050

ENTRYPOINT ["/rest-ws-go"]