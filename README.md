![grants_badge](./grants_badge.png)

# Subscan Essentials

[![License: GPL](https://img.shields.io/badge/license-GPL-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/itering/subscan)](https://goreportcard.com/report/github.com/itering/subscan)
![CI/CD](https://github.com/subscan-explorer/subscan-essentials/workflows/subscan/badge.svg)

Subscan Essentials is a high-precision blockchain explorer scaffold supporting Substrate-based networks. Developed by
the Subscan team and powering [subscan.io](https://www.subscan.io/), it provides:

- Developer-friendly interface
- Standard/custom module parsing
- Extensible plugin system
- Multi-chain compatibility

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Configuration](#configuration)
    - [Running Services](#running-services)
- [Docker Deployment](#docker-deployment)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)
- [Resources](#resources)

## Features

- **Substrate Integration**
    - Custom type registration ([guide](/custom_type.md))
    - Indexes blocks, extrinsics, events, logs, and EVM data
- **Extensibility**
    - Custom plugins framework ([docs](/plugins))
    - Auto-generate plugin templates via [gen tool](https://github.com/itering/subscan-plugin/tree/master/tools)
- **APIs**
    - Built-in HTTP API documentation ([docs](/docs))

---

## Quick Start

### Prerequisites

- Linux/macOS
- Git
- **Go 1.23+**
- **Redis 3.0.4+**
- **MySQL 8.0+** or **PostgreSQL 16+**

### Installation

```bash
./build.sh build

```

### configuration

#### Init config file

```bash
cp configs/config.yaml.example configs/config.yaml
```

```yaml
server:
  http:
    addr: 0.0.0.0:4399 # http api port
    timeout: 30s       # http timeout 
database:
  mysql:
    api: "mysql://root:helloload@127.0.0.1:3306?writeTimeout=3s&parseTime=true&loc=Local&charset=utf8mb4,utf8" # mysql default dsn
  postgres:
    api: "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable" # postgres default dsn
redis:
  proto: tcp
  addr: 127.0.0.1:6379 # redis host
  password: "" # redis password
  read_timeout: 1s 
  write_timeout: 1s
  idle: 10
  active: 100
UI:
  enable_substrate: true # if true, ui will show substrate data
  enable_evm: true       # if true, ui will show evm data
```



## Available Environment Variables




### Common

| Name                   | Default Value | Describe               |
|------------------------|---------------|------------------------|
| CONF_DIR               | ../configs    | configs path           |
| VERIFY_SERVER          | NULL          | solidity verify server |
| SUBSTRATE_ADDRESS_TYPE | 0             | ss58 address type      |
| SUBSTRATE_ACCURACY     | 10            | native token accuracy  |
| CHAIN_WS_ENDPOINT      |               | websocket endpoint url |
| NETWORK_NODE           | moonbeam      | network node name      |
| WORKER_GOROUTINE_COUNT | 10            | worker goroutine count |
| ETH_RPC                |               | Evm rpc endpoint       |

### Database

| Name              | Default Value      | Describe               |
|-------------------|--------------------|------------------------|
| DB_DRIVER         | mysql              | support mysql/postgres |
| MYSQL_HOST        | 127.0.0.1          | mysql host             |
| MYSQL_USER        | root               | mysql user             |
| MYSQL_PASS        |                    | mysql password         |
| MYSQL_DB          | subscan-essentials | mysql db name          |
| MYSQL_PORT        | 3306               | mysql port             |
| POSTGRES_HOST     | 127.0.0.1          | postgres port          |
| POSTGRES_USER     | gorm               | postgres user          |
| POSTGRES_PASS     | gorm               | postgres password      |
| POSTGRES_DB       | subscan-essentials | postgres db name       |
| POSTGRES_PORT     | 9920               | postgres port          |
| POSTGRES_SSL_MODE | disable            | postgres ssl mode      |
| MAX_DB_CONN_COUNT | 200                | gorm max db conn count |

### Redis

| Name           | Default Value | Describe                   |
|----------------|---------------|----------------------------|
| REDIS_HOST     | 127.0.0.1     | redis host                 |
| REDIS_PORT     | 6379          | redis host port            |
| REDIS_DATABASE | 0             | redis db                   |
| REDIS_PASSWORD |               | redis password default nil |

### running-services

- Start DB

**Make sure you have started redis and mysql/postgres**

- Subscribe

```bash
cd cmd && ./subscan start subscribe
```

- Worker

```bash
cd cmd && ./subscan start worker
```

- Api Server

```bash
cd cmd && ./subscan
```

- Help

```
NAME:
   SUBSCAN - SUBSCAN Backend Service, use -h get help

USAGE:
   cmd [global options] command [command options] [arguments...]

VERSION:
   2.0

DESCRIPTION:
   SubScan Backend Service, substrate blockchain explorer

COMMANDS:
   start              Start one worker, E.g. subscribe
   install            Install default database and create default conf file
   CheckCompleteness  Create blocks completeness
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --conf value   (default: "../configs")
   --help, -h     show help
   --version, -v  print the version

```

### docker-deployment

Use [docker-compose](https://docs.docker.com/compose/) can start projects quickly

Create local network

```
docker network create app_net
```

Run mysql and redis container

```bash
docker-compose -f docker-compose.db.yml up  -d
```

Run subscan service

```bash
docker-compose build
docker-compose up -d
```

### testing

1. Create test database (if using MySQL):

```sql
CREATE DATABASE subscan-essentials;
```

2. Run tests:

```bash
go test -v ./...
```

## contributing

We welcome contributions! Please see CONTRIBUTING.md for guidelines. Good first issues are labeled with **good first
issue**.

## LICENSE

GPL-3.0

## resources

- [SUBSCAN] https://github.com/subscan-explorer
- [scale.go] https://github.com/subscan-explorer/scale.go SCALE codec implementation
- [Darwinia] https://github.com/darwinia-network
