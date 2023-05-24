# Subscan Essentials

![License: GPL](https://img.shields.io/badge/license-GPL-blue.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/itering/subscan)](https://goreportcard.com/report/github.com/itering/subscan)
![subscan](https://github.com/itering/subscan/workflows/subscan/badge.svg)

Subscan Essentials is a high-precision blockchain explorer scaffold project.
It supports substrate-based blockchain networks with developer-friendly interface, standard or custom module parsing capabilities.
It's developed by the Subscan team and used in subscan.io.
Subscan Essentials is a high-precision blockchain explorer scaffold project.
It supports substrate-based blockchain networks with developer-friendly interface, standard or custom module parsing capabilities.
It's developed by the Subscan team and used in subscan.io.
Developers are free to use the codebase to extend functionalities and develop unique user experiences for their audiences.

## Contents

- [Feature](#feature)
- [QuickStart](#quickstart)
  - [Requirement](#requirement)
  - [Structure](docs/tree.md)
  - [Installation](#install)
  - [Config](#config)
  - [Usage](#usage)
  - [Docker](#docker)
  - [Test](#test)
- [Contributions](#contributions)
- [LICENSE](#license)
- [Resource](#resource)

## Feature

1. Support Substrate network [custom](/custom_type.md) type registration
1. Support Substrate network [custom](/custom_type.md) type registration
2. Support index Block, Extrinsic, Event, log
3. More data can be indexed by custom [plugins](/plugins)
4. [Gen](https://github.com/itering/subscan-plugin/tree/master/tools) tool can automatically generate plugin templates
5. Built-in default HTTP API [DOC](/docs/index.md)

## QuickStart

### Requirement

- Linux / Mac OSX
- Git
- Golang 1.20+
- Redis 3.0.4+
- MySQL 5.6+
- Node 8.9.0+

### Install

```bash
./build.sh build
```

### Config

#### Init config file
#### Init config file

```bash
cp configs/config.yaml.example configs/config.yaml
```

#### Set

1. Redis  configs/redis.toml

   > addr： redis host and port (default: 127.0.0.1:6379)

2. Mysql  configs/mysql.toml

   > host: mysql host (default: 127.0.0.1)
   > user: mysql user (default: root)
   > pass: mysql user passwd (default: "")
   > db:   mysql db name (default: "subscan")

3. Http   configs/http.toml

   > addr: local http server port (default: 0.0.0.0:4399)

### Usage

- Start DB
   **Make sure you have started redis and mysql**
- Substrate Daemon

   ```bash
   cd cmd
   ./subscan start substrate
   ```

- Api Server

   ```bash
   cd cmd
   ./subscan
   ```

- Help
- Help

   ```text
   NAME:
      SubScan - SubScan Backend Service, use -h get help

   USAGE:
      main [global options] command [command options] [arguments...]

   VERSION:
      1.0

   DESCRIPTION:
      SubScan Backend Service, substrate blockchain explorer

   COMMANDS:
      start    Start one worker, E.g substrate
      install  Create database and create default conf file
      help, h  Shows a list of commands or help for one command

   GLOBAL OPTIONS:
      --conf value   (default: "../configs")
      --help, -h     show help
      --version, -v  print the version


   ```

### Docker

If you prefer running containers, you can use [docker-compose](https://docs.docker.com/compose/).

Create local network

```bash
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

### Test

Note: The default test mysql database is subscan_test. Please CREATE it or change configs/mysql.toml.

```bash
go test ./...
```

## Contributions

We welcome contributions of any kind. Issues labeled can be good (first) contributions.

## LICENSE

GPL-3.0

## Resource

- [ITERING](https://github.com/itering)
- [Darwinia](https://github.com/darwinia-network/darwinia)
