FROM golang:1.20.4-bullseye as builder
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get upgrade -y

WORKDIR /subscan

COPY go.mod go.sum ./
RUN go mod download
COPY . /subscan
WORKDIR /subscan/cmd
RUN go build -o subscan

FROM buildpack-deps:bullseye-scm
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get upgrade -y

WORKDIR subscan
COPY configs configs
COPY configs/config.yaml.example configs/config.yaml

COPY --from=builder /subscan/cmd/subscan cmd/subscan
WORKDIR cmd
RUN mkdir -p /subscan/log


ENTRYPOINT ["/subscan/cmd/subscan"]
EXPOSE 4399