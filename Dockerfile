FROM golang:1.20.5-bullseye as builder
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get upgrade -y

WORKDIR /subscan

COPY . ./
RUN ./build.sh build

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
EXPOSE 80
