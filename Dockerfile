FROM golang:1.12.4 as builder

WORKDIR /subscan-end

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . /subscan-end

WORKDIR /subscan-end/cmd

RUN go build -o subscan

FROM buildpack-deps:buster-scm

RUN mkdir subscan

WORKDIR subscan

RUN mkdir log

COPY configs configs

COPY --from=builder /subscan-end/cmd/subscan cmd/subscan

COPY cmd/run.py cmd/run.py

WORKDIR cmd

EXPOSE 4399

CMD ["/subscan/cmd/subscan","-conf", "../configs"]