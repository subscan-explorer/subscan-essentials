FROM golang:1.23.9 as builder

WORKDIR /subscan

COPY go.mod go.sum ./
RUN go mod download
COPY . /subscan
WORKDIR /subscan/cmd
RUN go build -o subscan

FROM alpine:3

WORKDIR subscan
COPY configs configs
COPY configs/config.yaml.example configs/config.yaml

COPY --from=builder /subscan/cmd/subscan cmd/subscan
WORKDIR cmd
RUN apk update && apk add gcompat
ENTRYPOINT ["/subscan/cmd/subscan"]
EXPOSE 4399