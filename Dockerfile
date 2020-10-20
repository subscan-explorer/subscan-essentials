FROM golang:1.14.6 as builder

WORKDIR /subscan

COPY go.mod go.sum ./
RUN go mod download
COPY . /subscan
WORKDIR /subscan/cmd
RUN go build -o subscan

FROM buildpack-deps:buster-scm

WORKDIR subscan
COPY configs configs
COPY configs/redis.toml.example configs/redis.toml
COPY configs/mysql.toml.example configs/mysql.toml
COPY configs/http.toml.example configs/http.toml

COPY --from=builder /subscan/cmd/subscan cmd/subscan
COPY cmd/run.py cmd/run.py
WORKDIR cmd
RUN mkdir -p /subscan/log

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

CMD ["/subscan/cmd/subscan"]
EXPOSE 4399