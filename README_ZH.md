# SUBSCAN

![License: GPL](https://img.shields.io/badge/license-GPL-blue.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/itering/subscan)](https://goreportcard.com/report/github.com/itering/subscan)
![subscan](https://github.com/itering/subscan/workflows/subscan/badge.svg)


SUBSCAN是一个高精度的区块链浏览器，聚合Substrate生态网络并为您提供完美的浏览体验。

此项目是[subscan](https://subscan.io)的开源版本


## API doc

默认的API可以在这边找到 [DOC](/docs/index.md)


### 功能

1. API Server 与后台监听程序分离
2. 支持substrate 网络自定义type 注册
3. 支持索引block, Extrinsic, Event, log
4. 可自定义插件索引更多的数据[Plugins](/plugins)
5. [gen](/tools/gen-plugin)工具可自动生成插件模版
6. 内置默认的HTTP API [DOC](/docs/index.md)


### 安装

```bash
make &&  ./cmd/subscan --conf configs install
```

### 运行

> API 

```bash

./cmd/subscan --conf configs

```

> Daemon

```bash
./cmd/subscan --conf configs start substrate
./cmd/subscan --conf configs stop substrate
```


### docker

```bash

docker-composer build

docker-composer up -d

```

## LICENSE

GPL-3.0


## resource
 
[ITERING] https://github.com/itering

[SUBSCAN] https://subscan.io/

[Darwinia] https://github.com/darwinia-network/darwinia

[freehere107] https://github.com/freehere107