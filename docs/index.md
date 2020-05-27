# API Document

1.  [socket](#socket)
1.  [subscribe](#subscribe)
1.  [now](#now)
1.  [metadata](#metadata)
1.  [blocks](#blocks)
1.  [extrinsics](#extrinsics)
1.  [extrinsic](#extrinsic)
1.  [events](#events)
1.  [event](#event)
1.  [search](#search)
1.  [dailyStat](#daily)
1.  [transfers](#transfers)
1.  [check_hash](#check-hash)


## Description

### Global Header 

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| Content-Type | application/json | yes     |


#### List row&page

| Name          | Type   | Desc |
| ------------- | ------ | ------- |
|row| int | min=1,max=100     |
|page| int | min=0     |



## socket

### URL Request

`ws /socket`

### Example Response

topic
> metadata_update: metadata更新
 
> block_new:新的block
  
    {"content":{"count_event":115},"time":1563775286,"topic":"metadata_update"}
 
    {"content":{"block_num":181272,"created_at":"2019-07-22T14:01:19.418808+08:00","decode_event":"[{\"phase\": 0, \"extrinsic_idx\": 0, \"type\": \"0000\", \"module_id\": \"system\", \"event_id\": \"ExtrinsicSuccess\", \"params\": [], \"event_idx\": 0}]","decode_extrinsics":"[{\"valueRaw\": \"010100\", \"extrinsic_length\": 8, \"version_info\": \"01\", \"call_code\": \"0100\", \"call_module_function\": \"set\", \"call_module\": \"timestamp\", \"params\": [{\"name\": \"now\", \"type\": \"Compact\u003cMoment\u003e\", \"value\": \"2019-07-22T06:01:18\", \"valueRaw\": \"032e51355d\"}]}]","decode_logs":"[{\"index\": \"PreRuntime\", \"type\": \"(u32, Bytes)\", \"value\": {\"engine\": 1634891105, \"data\": \"dde2880f00000000\"}}, {\"index\": \"Seal\", \"type\": \"(u32, Bytes)\", \"value\": {\"engine\": 1634891105, \"data\": \"4a434cfd9da1e22d212c81831f2ac90db639a51d86cf510ce569c8e1bcb5925fda63537688d7065263b175cc2e2e9f1d63b812668cd2540e38d30deb46600204\"}}]","event":"0x040000000000000000","extrinsics":"[\"0x20010100032e51355d\"]","extrinsics_root":"0x83730e3c02f630a7a4ef81a2dba09d339954f261d94c70b64df9d4c89b992a96","hash":"0xfc8b26a440993737da41ddb673fc065482966790d060a24fda48f3e97ba4aa3a","id":2791,"logs":"[\"0x066175726120dde2880f00000000\",\"0x056175726101014a434cfd9da1e22d212c81831f2ac90db639a51d86cf510ce569c8e1bcb5925fda63537688d7065263b175cc2e2e9f1d63b812668cd2540e38d30deb46600204\"]","parent_hash":"0xdd929487fa95da203823d7652644c925ac2c8d42c6e4f4e0f1b5ec27e885911b","spec_version":78,"state_root":"0x43e5376c3eb78e5cd01158e351c8fcddf2293c6f8143b2a6208dc0e863388c63"},"time":1563775279,"topic":"block_new"}



## subscribe
     
### URL Request

`POST /api/subscribe`

## FormData :

| Name          | Type   | Desc |
| ------------- | ------ | ------- |
|email| string |      |

### Example Response

`200 OK` and

 {
     "code": 0,
     "message": "success",
     "ttl": 1
 }

-----


## now

### URL Request

`POST /api/now`

### Example Response

`200 OK` and

    {
        "code": 0,
        "message": "success",
        "ttl": 1,
        "data": 1559545576
    }

-----

## metadata

### URL Request

`POST /api/scan/metadata`

### Example Response

`200 OK` and

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


>signed: signed|all

### Example Response

`200 OK` and

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
-----


## event

### URL Request

`POST /api/scan/event`

### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| event_index | string | no     |

### Example Response

`200 OK` and
    
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
-----

## search

### URL Request

`POST /api/scan/search`

### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| key | string | yes     |
| row | int | yes     |
| page | int | yes     |


-----

## daily

### URL Request

`POST /api/scan/daily`

### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| start | Date(2019-07-04) | yes     |
| end | Date(2019-07-04) | yes     |


### Example Response

`200 OK` and
    
    {
        "code": 0,
        "message": "Success",
        "ttl": 1,
        "data": {
            "list": [
                {
                    "ID": 1,
                    "time_utc": "2019-07-02T00:00:00+08:00",
                    "transfer_count": 109
                },
                {
                    "ID": 42,
                    "time_utc": "2019-07-03T00:00:00+08:00",
                    "transfer_count": 69
                },
                {
                    "ID": 68,
                    "time_utc": "2019-07-04T00:00:00+08:00",
                    "transfer_count": 44
                }
            ]
        }
    }
-----


## transfers

### URL Request

`POST /api/scan/transfers`

### payload

| Name          | Type   | Require |
| ------------- | ------ | ------- |
| row | int | yes     |
| page| int | yes     |
| address| string | no     |

### Example Response

`200 OK` and

    {
        "code": 0,
        "message": "Success",
        "ttl": 1,
        "data": {
            "count": 0,
            "extrinsics": [
                {
                    "from": "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
                    "to": "5CXciU5Pamk9z67KFiqG1r8HZKqGLYKJLCdXLQjVfcPU42KM",
                    "module": "balances",
                    "amount": "0.0001",
                    "hash": "0x8d2ad11ee3cbd0f58286a2d67739db4b7d10210d394bb8ca7a1168006fcb41ca"
                },
            ]
        }
    }

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

    {
        "code": 0,
        "message": "Success",
        "ttl": 1,
        "data": {
            "hash_type": "block/extrinsic"
        }
    }
-----



### Possible Responses Code

| HTTP Code        | Condition              |
| ---------------- | ---------------------- |
| 200 OK           | successful             |
| 500 Server Error | API server error       |
| 400 Bad request  | Invalid params or else |
