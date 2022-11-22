# [server api doc]

|Version|Update Time|Author|
|:--------:|:--------:|:--------:|
|1.0|2021-06-01|huanxiong@dtstack.com|

## Agent Install

### Request

    api/v1/agent/install

### Method

    POST

### Params

| param | type | desc | required |
| ---- | ----| ---- | ---- |
|collectorId|string|sidecar UUID|yes|
|configurationPath|string|agent configuration path||
|binaryPath|string|agent binary file|yes|
|name|string|agent name|yes|
|parameter|string|agent parameters，more than one parameter separated by ＇,＇||
|installScript|string|install shell|yes|
|timeout|string|timeout of install shell exec, default 15m eg. ||
|runUser|string|exec user of the agent, default as the user of the sideccar ||
|installParameter|string|install shell parameter||
|healthShell|string|health check shell||
|healthPeriod|string|health check period, 1s/1m/1h...||
|healthTimeout|string|timeout of health check shell, 20s eg. ||

### Response

    {
    "msg": "ok",
    "code": 0,
    "data": {
        "agent_id": "****************",
        "operation_seq": 12
        }
    }

> description

|value|type|desc|
|--------|--------|--------|
|data.agent_id|string|agent uuid|
|data.operation_seq|int|operation seq|　　

## Agent Uninstall


### Request

     api/v1/agent/{agent_id}/uninstall

### Method

    POST


### Params

| Param | Type | Description |
| ---- | ----| ---- |
|parameter|string|parameters of the shellScript ，separated by ＇,＇|
|shellScript|string|run shellScript to clear sth if needed|

### Response

    {
    "msg": "ok",
    "code": 0,
    "data": {
        "agent_id": "****************",
        "opration_seq": 12
        }
    }

> Description

|Value|Description|
|:--------:|:--------:|
|data.agent_id||
|data.opration_seq||

## agent start|stop|restart


### Request

    api/v1/agent/{agent_id}/start
    api/v1/agent/{agent_id}/stop
    api/v1/agent/{agent_id}/restart

### Method

    GET


### Params

| Param | Type | Description |
| ---- | ----| ---- |
|agent_id|string|agent uuid|


### Response

    {
    "msg": "ok",
    "code": 0,
    "data": {
        "agent_id": "****************",
        "opration_seq": 12
        }
    }

> Description

|Value|Description|
|:--------:|:--------:|
|data.agent_id|agent uuid|
|data.opration_seq|operation seq|


## Agent config get/update

### Request

    api/v1/agent/{agent_id}/config

### Method

    GET


### Params

| Param | Type | Description |
| ---- | ----| ---- |
|agent_id|string|agent uuid|

### Response

    {
    "msg": "ok",
    "code": 0,
    "data": {
        "agent_id": "****************",
        "config_content": "this is the content of configuration!"
        }
    }

> Description

|Value|Type|Description|
|--------|--------|--------|
|data.agent_id|string|agent uuid|
|data.config_content|string|agent configuration content|　　



##  agent configuration update

### Request

    api/v1/agent/{agent_id}/config

### Method

    POST


### Params

| Param | Type | Description |
| ---- | ----| ---- |
|agent_id|string|agent uuid|
|config_content|string|agent configuration context|

### Response

    {
    "msg": "ok",
    "code": 0,
    "data": {
        "agent_id": "****************",
        "opration_seq": 12
        }
    }

> Description

|Value|Type|Description|
|--------|--------|--------|
|data.agent_id|string|agent uuid|
|data.opration_seq|operation seq|


## Exec script on the target sidecar(host)


### Request

	api/v1/sidecar/{sidecar_id:uuid}/execscript

### Method

	POST

### Params


| Param | Type | Description |
| ---- | ---- | ---- |
|execScript|string|the content of the script|
|parameter|string|Params，separated by ＇,＇|
|timeout|string|A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".|


### Response

    {
    "msg": "ok",
    "code": 0,
    "data": {
        "sidecar_id": "****************",
        "operation_seq": 12
        }
    }

> Description

|Value|Type|Description|
|--------|--------|--------|
|data.agent_id|string|agent uui|
|data.operation_seq|int||　


## Get Progress of operation by seq and agentId

### Request

    api/v1/agent/{agent_id}/progress

### Method

    GET

### Params

| Param | Type | Description |
| ---- | ----| ---- |
|agent_id|string|agent uuid|
|operation_seq|string||


### Response

    {
    "msg": "ok",
    "code": 0,
    "data": {
        "operation_seq": 23,
        "op_time":2017-07-28T15:45:32,
        "op_name":"install"|"updateInstall",
        "progress": 0.24,
        "op_result":0,
        "ret_msg": "",
        "msg": "",
        "finish_time":2017-07-28T16:45:32,
        "ts":2017-07-28T16:45:32,
        "collector_id":"sdasdecaccrsdaaa",
        "agent_id": "cskdfmpoervmkef"
        }
    }

> description

|value|type|desc|
|--------|--------|--------|
|data.operation_seq|int|install operation uuid seq|
|data.op_time|date|start time of install operation|
|data.op_name|string|name of install operation|
|data.progress|float|the progress of the install operation|
|data.op_result|int|0 is success, 1 is failed|
|data.ret_msg|string|install shell return  message|
|data.msg|string|result message|
|data.finish_time|date|the finish time of the operation|
|data.ts|date|timestampe of the progress|
|data.collector_id|string|sidecar uuid|
|data.agent_id|string|agent uuid|

