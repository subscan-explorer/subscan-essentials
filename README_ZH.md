# Subscan Essentials

![License: GPL](https://img.shields.io/badge/license-GPL-blue.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/itering/subscan)](https://goreportcard.com/report/github.com/itering/subscan)
![subscan](https://github.com/itering/subscan/workflows/subscan/badge.svg)

Subscan Essentials是一个高精度的区块链浏览器脚手架项目，它具有开发人员友好的界面和自定义模块解析功能，支持基于substrate的区块链网络。 它由Subscan团队开发，并在subscan.io中使用。
开发人员可以自由使用代码库来扩展功能并为其受众开发独特的用户体验。


## API doc

默认的API文档可以在这边找到 [DOC](/docs/index.md)


### 功能

1. API Server 与后台监听程序分离
2. 支持substrate 网络自定义type 注册 [Custom](/custom_type.md)
3. 支持索引block, Extrinsic, Event, log
4. 可自定义插件索引更多的数据[Plugins](/plugins)
5. [Gen](https://github.com/itering/subscan-plugin/tree/master/tool)工具可自动生成插件模版
6. 内置默认的HTTP API [DOC](/docs/index.md)


### 安装

```bash
./build.sh build &&  ./cmd/subscan --conf configs install
```

### 使用

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