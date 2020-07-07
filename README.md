# Subscan Essentials

![License: GPL](https://img.shields.io/badge/license-GPL-blue.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/itering/subscan)](https://goreportcard.com/report/github.com/itering/subscan)
![subscan](https://github.com/itering/subscan/workflows/subscan/badge.svg)

Subscan Essentials is a high-precision blockchain explorer scaffold project, supports substrate-based blockchain networks with developer-friendly interface and custom module parsing capabilities.  It's developed by the Subscan team and used in subscan.io.  Developers are free to use the codebase to extend functionalities and develop unique user experiences for their audiences.

## API doc

The default API can be found here [DOC](/docs/index.md)


### Feature

1. Separation of API Server and background listener
2. Support Substrate network custom type registration [Type](/custom_type.md)
3. Support index block, Extrinsic, Event, log
4. More data can be indexed by custom plugins [Plugins](/plugins)
5. [gen](/tools/gen-plugin) tool can automatically generate plug-in templates
6. Built-in default HTTP API [DOC](/docs/index.md)

### Install

```bash
make &&  ./cmd/subscan --conf configs install
```

### RUN

> API 

```bash

./cmd/subscan --conf configs

```

> Daemon

```bash
./cmd/subscan --conf configs start substrate
./cmd/subscan --conf configs stop substrate
```


### Docker

```bash

docker-composer build

docker-composer up -d

```

## LICENSE

GPL-3.0


## Resource
 
[ITERING] https://github.com/itering

[SUBSCAN] https://subscan.io/

[Darwinia] https://github.com/darwinia-network/darwinia

[freehere107] https://github.com/freehere107
