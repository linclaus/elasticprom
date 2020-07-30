FROM golang:1.13.11-alpine3.11 as builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

#RUN CGO_ENABLED=0 \
#    GOOS=linux \
#    GOARCH=amd64 \
#    go build -mod=readonly -ldflags="-w -s" -o mantisd cmd/main.go
#
#FROM scratch

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -mod=readonly --gcflags="-N -l" -o elasticpromd cmd/main.go

FROM alpine:3.11

WORKDIR /app

COPY --from=builder /build/elasticpromd .

CMD ["/app/elasticpromd"]