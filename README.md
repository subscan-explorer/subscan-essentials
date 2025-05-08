![grants_badge](./grants_badge.png)

# Subscan Essentials

![License: GPL](https://img.shields.io/badge/license-GPL-blue.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/itering/subscan)](https://goreportcard.com/report/github.com/itering/subscan)
![subscan](https://github.com/subscan-explorer/subscan-essentials/workflows/subscan/badge.svg)

Subscan Essentials is a high-precision blockchain explorer scaffold project. 
It supports substrate-based blockchain networks with developer-friendly interface, standard or custom module parsing capabilities. 
It's developed by the Subscan team and used in [subscan.io](https://www.subscan.io/). 
Developers are free to use the codebase to extend functionalities and develop unique user experiences for their audiences.

## Contents

- [Feature](#Feature)
- [QuickStart](#QuickStart)
  - [Requirement](#Requirement)
  - [Structure](docs/tree.md)
  - [Installation](#Install)
  - [UI](#UI)
  - [Config](#Config)
  - [Usage](#Usage)
  - [Docker](#Docker)
  - [Test](#Test)
- [Contributions](#Contributions)
- [LICENSE](#LICENSE)
- [Resource](#Resource)

## Feature

1. Support Substrate network [custom](/custom_type.md) type registration 
2. Support index Block, Extrinsic, Event, log, EVM data(block, transaction...)
3. More data can be indexed by custom [plugins](/plugins)
4. [Gen](https://github.com/itering/subscan-plugin/tree/master/tools) tool can automatically generate plugin templates
5. Built-in default HTTP API [DOC](/docs/index.md)


## QuickStart

### Requirement

* Linux / Mac OSX
* Git
* Golang 1.23.0+
* Redis 3.0.4+
* MySQL 8.0+

### Install

```bash
./build.sh build

```

### Config

#### Init config file 

```bash
cp configs/config.yaml.example configs/config.yaml
```

## Environment Variables

### Common

| Name                   | Default Value | Describe               |
|------------------------|---------------|------------------------|
| CONF_DIR               | ../configs    | configs path           |
| VERIFY_SERVER          | NULL          | solidity verify server |
| SUBSTRATE_ADDRESS_TYPE | 0             | ss58 address type      |
| SUBSTRATE_ACCURACY     | 10            | native token accuracy  |



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

| Name                  | Default Value | Describe |
|-----------------------|---------------|----------|
| REDIS_HOST            | 127.0.0.1     |          |
| REDIS_PORT            | 6379          |          |
| REDIS_DATABASE        | 0             |          |
| REDIS_PASSWORD        |               |          |


### Usage

- Start DB

**Make sure you have started redis and mysql**

- Subscribe
```bash
cd cmd
./subscan start subscribe
```

- Worker
```bash
cd cmd
./subscan start worker
```

- Api Server
```bash
cd cmd
./subscan
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

### Docker

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

### Test


**default test mysql database is subscan_test. Please CREATE it or change configs/config.yaml**

```bash
go test ./...
```


## Contributions

We welcome contributions of any kind. Issues labeled can be good (first) contributions.

## LICENSE

GPL-3.0


## Resource
 
- [SUBSCAN] https://github.com/subscan-explorer
- [Darwinia] https://github.com/darwinia-network
