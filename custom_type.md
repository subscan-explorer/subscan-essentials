## Substrate Custom Type format


Subscan custom type most cases are consistent with polkadot.js type.json, due to language characteristics, 
there are some inconsistencies between struct and enum


### String

```json

{
  "typeName": "inheritTypeName"
}
```

Example

```json
{
  "address": "H256"
}

```

### Struct

```json

{
  "typeName": {
    "type": "struct",
    "type_mapping": [
        [
          "field1", 
          "inheritTypeName"
        ]
    ]
  }
}
```


Example
```json
{
    "BalanceLock<Balance>": {
      "type": "struct",
      "type_mapping": [
        [
          "id", 
          "LockIdentifier"
        ]
      ]
    }
}

```


### Enum


```json

{
  "typeName": {
    "type": "enum",
    "type_mapping": [
        [
          "field1", 
          "inheritTypeName"
        ]
    ]
  }
}
```


Example
```json
{
  "RedeemStrategy": {
    "type": "enum",
    "type_mapping": [
      [
        "Immediately",
        "Null"
      ]
    ]
  }
}

```


### Resource

[scale.go] https://github.com/itering/scale.go/tree/master/source