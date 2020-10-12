# API Document

1.  [metadata](#metadata)
1.  [blocks](#blocks)
1.  [block](#block)
1.  [extrinsics](#extrinsics)
1.  [extrinsic](#extrinsic)
1.  [events](#events)
1.  [check_hash](#check-hash)
1.  [runtime_list](#runtime-list)
1.  [runtime_info](#runtime-info)
1.  [plugin_list](#plugin-list)
1.  [plugin_ui](#plugin-ui)


## Description

### Global Header 

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| Content-Type | application/json | yes     |


#### List row&page

| Name          | Type   | Desc |
| ------------- | ------ | ------- |
|row| int | min=1,max=100 defines how many results should be presented per page    |
|page| int | min=0 defines which page should be shown   |


-----

## metadata

### URL Request

`POST /api/scan/metadata`

### Example Response

`200 OK` and

```json

{
    "code": 0,
    "message": "Success",
    "ttl": 1,
    "data": {
        "blockNum": "1825213",
        "count_event": "39706",
        "count_extrinsic": "39016",
        "count_signed_extrinsic": "94"
    }
}
```
-----

## blocks

### URL Request

`POST /api/scan/blocks`

### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| row | int | yes     |
| page| int | yes     |

### Example Response

`200 OK` and

```json

{
    "code": 0,
    "message": "Success",
    "ttl": 1,
    "data": {
        "blocks": [
            {
                "id": 21286,
                "block_num": 1825212,
                "created_at": "2019-06-18T18:19:22+08:00",
                "hash": "0x0b2f52c4b744df62ea93724e15bcdce85bbcf8687c1ca00121194895f8afce8d",
                "parent_hash": "0xc06ccdc2ae060b686620f95649ad7fa8025ae2d9e77eda389b8937373d12ef98",
                "state_root": "0xfa05f4fd97e256596605fabb609b399ce994371b17a6aa6d8989d4c3045047c7",
                "extrinsics_root": "0x3a4e596ab65e724f1a38d6eaecd0404c2d18ed8ee3d273dcc1ff949f3bbf1a13",
                "logs": "[\"0x046175726121017174810f00000000d393eca3bf848b7f459abe58ea9e8c555eafbfe3f46c65d3a02b103b157cfa2a8305dd0c274c1a492718ad28cd6f5a86d196ca3916980bbda71012047fb95904\"]",
                "extrinsics": "[\"0x01000003a6ba085d\",\"0x010d0000\"]",
                "decode_extrinsics": "[{\"valueRaw\": \"010000\", \"extrinsic_length\": null, \"version_info\": \"01\", \"call_code\": \"0000\", \"call_module_function\": \"set\", \"call_module\": \"timestamp\", \"params\": [{\"name\": \"now\", \"type\": \"Compact<Moment>\", \"value\": \"2019-06-18T10:19:18\", \"valueRaw\": \"03a6ba085d\"}]}, {\"valueRaw\": \"010d00\", \"extrinsic_length\": null, \"version_info\": \"01\", \"call_code\": \"0d00\", \"call_module_function\": \"set_heads\", \"call_module\": \"parachains\", \"params\": [{\"name\": \"heads\", \"type\": \"Vec<AttestedCandidate>\", \"value\": [], \"valueRaw\": \"\"}]}]",
                "event": "0x080000000000000000010000000000",
                "decode_event": "[{\"phase\": 0, \"extrinsic_idx\": 0, \"type\": \"0000\", \"module_id\": \"system\", \"event_id\": \"ExtrinsicSuccess\", \"params\": [], \"event_idx\": 0}, {\"phase\": 0, \"extrinsic_idx\": 1, \"type\": \"0000\", \"module_id\": \"system\", \"event_id\": \"ExtrinsicSuccess\", \"params\": [], \"event_idx\": 1}]"
            }
        ],
        "count": 21215
    }
}
```
-----

## block

### URL Request

`POST /api/scan/block`

### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| block_num | int | no     |
| block_hash| string | no     |


### Example Response

`200 OK` and

```json

    {
        "code": 0,
        "message": "Success",
        "ttl": 1,
        "data": {
            "block_num": 1763772,
            "created_at": "2019-06-14T11:40:35+08:00",
            "hash": "0x953e762c5bda331f4887958a38936cef5e51d48849238aafc975c2948f56917d",
            "parent_hash": "0x738f3334448f53cd6818ad4c2b61a8be9920438dcf008c6db91a2b65cebea86c",
            "state_root": "0xea541acde43379121972712be76e9e43f328882b54c3814ec0b4a94895a87ade",
            "extrinsics_root": "0x9aa30a3b0b2c237fa515f0f30b2d53bf6388b3b116bf0f3ecf98e9450f0d2c4e",
            "extrinsics": [
                {
                    "extrinsic_index": "",
                    "value_raw": "",
                    "extrinsic_length": "",
                    "version_info": "01",
                    "call_code": "0000",
                    "call_module_function": "set",
                    "call_module": "set",
                    "params": "[{\"name\":\"now\",\"type\":\"Compact\\u003cMoment\\u003e\",\"value\":\"2019-06-14T03:40:30\",\"valueRaw\":\"032e17035d\"}]",
                    "account_length": "",
                    "account_id": "",
                    "account_index": "",
                    "signature": "",
                    "nonce": 0,
                    "era": "",
                    "extrinsic_hash": "",
                    "success": false
                }
            ],
            "events": [
                {
                    "event_index": "",
                    "phase": 0,
                    "extrinsic_idx": 0,
                    "type": "0000",
                    "module_id": "system",
                    "event_id": "ExtrinsicSuccess",
                    "params": "[]",
                    "event_idx": 0
                }
            ]
        }
    }
```


-----

## extrinsics

### URL Request

`POST /api/scan/extrinsics`


### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| row | int | yes     |
| page| int | yes     |
| signed| string | no     |
| address| string | no     |
| module| string | no     |
| call| string | no     |

### Example Response

`200 OK` and

```json

{
    "code": 0,
    "message": "Success",
    "ttl": 1,
    "data": {
        "count": 613,
        "extrinsics": [
            {
                "block_timestamp": 0,
                "block_num": 71009,
                "extrinsic_index": "71009-2",
                "value_raw": "",
                "extrinsic_length": "",
                "version_info": "81",
                "call_code": "0a00",
                "call_module_function": "bond",
                "call_module": "staking",
                "params": "[{\"name\":\"controller\",\"type\":\"Address\",\"value\":\"ceeb8cf9ce7c6c23fc2a70d861eac4e6e94214f6db1063cb1279f5262b99b41f\",\"valueRaw\":\"ffceeb8cf9ce7c6c23fc2a70d861eac4e6e94214f6db1063cb1279f5262b99b41f\"},{\"name\":\"value\",\"type\":\"Compact\\u003cBalanceOf\\u003e\",\"value\":10,\"valueRaw\":\"28\"},{\"name\":\"payee\",\"type\":\"RewardDestination\",\"value\":\"Staked\",\"valueRaw\":\"00\"}]",
                "account_length": "ff",
                "account_id": "6caf778f29c84796e5dcba9122fcd4ca838d9caffd496e140f831f9faa80695c",
                "account_index": "",
                "signature": "9689cf2a20f670d1b3b100d3a8cc7f8f01fa7177835d4985179272908508f46c489341a1cb71c9e03402a2f6866cb22aed40fef87832500c1e1cb2b885fccc08",
                "nonce": 8,
                "era": "00",
                "extrinsic_hash": "0xb20a86ebeb3f59c27bec60d4bcce4cc2a0ca029e5ccff64f25779ad8562a8ea3",
                "success": true
            }
        ]
    }
}   
```

-----

## extrinsic

### URL Request

`POST /api/scan/extrinsic`

### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| extrinsic_index | string | no     |
| hash| string | no     |

### Example Response

`200 OK` and

```json

{
    "code": 0,
    "message": "Success",
    "ttl": 1,
    "data": {
        "block_timestamp": 0,
        "block_num": 820,
        "extrinsic_index": "820-2",
        "call_module_function": "transfer",
        "call_module": "balances",
        "account_id": "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
        "signature": "522d0e61a3a3bfa1f0584d999af4cb16d0ad9b5390324a5ecbb80cceb11093694d8629276dc8c41de032ada9212cb499f1f8f4ba4a5784cad8c3f58db65ada06",
        "nonce": 0,
        "extrinsic_hash": "0x8d2ad11ee3cbd0f58286a2d67739db4b7d10210d394bb8ca7a1168006fcb41ca",
        "success": false,
        "params": [
            {
                "name": "dest",
                "type": "Address",
                "value": "1486518478c79befe09ffe69dc1eb8cb862e29ee013097b021fafcb74642127b",
                "valueRaw": "ff1486518478c79befe09ffe69dc1eb8cb862e29ee013097b021fafcb74642127b"
            },
            {
                "name": "value",
                "type": "Compact<Balance>",
                "value": 100000,
                "valueRaw": "821a0600"
            }
        ],
        "transfer": {
            "from": "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
            "to": "5CXciU5Pamk9z67KFiqG1r8HZKqGLYKJLCdXLQjVfcPU42KM",
            "module": "balances",
            "amount": "0.0001",
            "hash": "0x8d2ad11ee3cbd0f58286a2d67739db4b7d10210d394bb8ca7a1168006fcb41ca"
        },
        "event": {
            "event_index": "820-2",
            "block_num": 820,
            "phase": 0,
            "extrinsic_idx": 0,
            "type": "0802",
            "module_id": "staking",
            "event_id": "OfflineSlash",
            "params": "[{\"type\":\"AccountId\",\"value\":\"0xfe65717dad0447d715f660a0a58411de509b42e6efb8375f562f58a554d5860e\",\"valueRaw\":\"fe65717dad0447d715f660a0a58411de509b42e6efb8375f562f58a554d5860e\"},{\"type\":\"Balance\",\"value\":0,\"valueRaw\":\"00000000000000000000000000000000\"}]",
            "event_idx": 0
        }
    }
}


``` 
-----

## events

### URL Request

`POST /api/scan/events`

### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| row | int | yes     |
| page| int | yes     |

### Example Response

`200 OK` and

```json

{
    "code": 0,
    "message": "Success",
    "ttl": 1,
    "data": {
        "count": 40084,
        "events": [
            {
                "event_index": "1825126-2",
                "phase": 1,
                "extrinsic_idx": 0,
                "type": "0300",
                "module_id": "session",
                "event_id": "NewSession",
                "params": "[{\"type\":\"BlockNumber\",\"value\":45356,\"valueRaw\":\"2cb1000000000000\"}]",
                "event_idx": 2
            }
        ]
    }
}
    
```
-----


## check-hash

### URL Request

`POST /api/scan/check_hash`

### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| hash | string | yes     |


### Example Response

`200 OK` and
```json

{
    "code": 0,
    "message": "Success",
    "ttl": 1,
    "data": {
        "hash_type": "block/extrinsic"
    }
}
    
```

-----

## runtime-list

### URL Request

`POST /api/scan/runtime/list`


### Example Response

`200 OK` and
```json
{
    "code": 0,
    "message": "Success",
    "ttl": 1,
    "data": {
        "list": [
            {
                "spec_version": 83,
                "modules": "RandomnessCollectiveFlip|System|Babe|Balances|Indices|Kton|Timestamp|Balances|Authorship|GrandpaFinality|ImOnline||Offences|Session|Staking|Contract|Sudo||EthRelay|EthBacking"
            }
        ]
    }
}
```
-----

## runtime-info

### URL Request

`POST /api/scan/runtime/metadata`

### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| spec| int | yes   |
| module| string | yes   |


### Example Response

`200 OK` and
```json
{
    "code": 0,
    "message": "Success",
    "ttl": 1,
    "data": {
        "info": {
            "name": "RandomnessCollectiveFlip",
            "prefix": "RandomnessCollectiveFlip",
            "storage": [
                {
                    "name": "RandomMaterial",
                    "modifier": "Default",
                    "type": {
                        "origin": "PlainType",
                        "plain_type": "Vec<Hash>"
                    },
                    "fallback": "0x00",
                    "docs": [
                        " Series of block headers from the last 81 blocks that acts as random seed material. This",
                        " is arranged as a ring buffer with `block_number % 81` being the index into the `Vec` of",
                        " the oldest hash."
                    ]
                }
            ],
            "calls": null,
            "events": null,
            "errors": null
        }
    }
}
```
-----

## plugin-list

### URL Request

`POST /api/scan/plugins`


### Example Response

`200 OK` and
```json
{
    "code": 0,
    "message": "Success",
    "ttl": 1,
    "data": {
        "list": [
            {
                "name": "balance",
                "version": "0.1"
            },
            {
                "name": "system",
                "version": "0.1"
            }
        ]
    }
}
```

-----
## plugin-ui

### URL Request

`POST /api/scan/plugins/ui`


### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| name| string | yes(plugin name)   |

### Example Response

`200 OK` and
```json
{
    "code": 0,
    "message": "Success",
    "ttl": 1,
    "data": {
        "type": "page",
        "body": {
            "type": "crud",
            "api": {
                "method": "POST",
                "url": "api/plugin/balance/accounts",
                "requestAdaptor": "return {...api, data: {...api.data, page: api.data.page - 1, row: api.data.perPage,} }",
                "adaptor": "return {...payload, status: payload.code, data: {items: payload.data.list, count: payload.data.count}, msg: payload.message };"
            },
            "syncLocation": false,
            "headerToolbar": [],
            "columns": [
                {
                    "name": "address",
                    "label": "address"
                },
                {
                    "name": "nonce",
                    "label": "nonce"
                },
                {
                    "name": "balance",
                    "label": "balance"
                },
                {
                    "name": "lock",
                    "label": "lock"
                }
            ]
        }
    }
}
```

-----

### Possible Responses Code

| HTTP Code        | Condition              |
| ---------------- | ---------------------- |
| 200 OK           | successful             |
| 500 Server Error | API server error       |
| 400 Bad request  | Invalid params or else |
