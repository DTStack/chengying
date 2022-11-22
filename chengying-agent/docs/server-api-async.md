# [server api doc]

|Version|Update Time|Author|
|:--------:|:--------:|:--------:|
|1.0|2021-06-01|huanxiong@dtstack.com|

## Agent Install

### Request

    api/v1/agent/installSync

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

    //#ok response
    {
    "msg": "ok",
    "code": 0,
    "data": {
        "agent_id": "****************",
        "operation_seq": 12,
        "result": {}
        }
    }

    //#wrong response
    {
    "msg": "DB Model error",
    "code": 102,
    "data": {

            }
    }

> description

|value|type|desc|
|--------|--------|--------|
|data.agent_id|string|agent uuid|
|data.operation_seq|int|operation seq|　　

## Agent Uninstall

### Request

     api/v1/agent/{agent_id}/uninstallSync

### Method

    POST


### Params

| Param | Type | Description |
| ---- | ----| ---- |
|parameter|string|parameters of the shellScript ，separated by ＇,＇|
|shellScript|string|run shellScript to clear sth if needed|

### Response

       //#ok response
        {
        "msg": "ok",
        "code": 0,
        "data": {
            "agent_id": "****************",
            "opration_seq": 12,
            "result": {}
            }
        }
    
        //#wrong response
        {
        "msg": "DB Model error",
        "code": 102,
        "data": {}
        }

> Description

|Value|Description|
|:--------:|:--------:|
|data.agent_id||
|data.opration_seq||


## agent start|stop|restart


### Request

    api/v1/agent/{agent_id}/startSync
    api/v1/agent/{agent_id}/stopSync
    api/v1/agent/{agent_id}/restartSync

### Method

    GET


### Params

| Param | Type | Description |
| ---- | ----| ---- |
|agent_id|string|agent uuid|


### Response
    //ok
    {
    "msg": "ok",
    "code": 0,
    "data": {
        "agent_id": "****************",
        "opration_seq": 12,
        "result": {}
        }
    }

    //wrong response

    {
    "msg": "DB Model error",
    "code": 102,
    "data": {}
    }
> Description

|Value|Description|
|:--------:|:--------:|
|data.agent_id|agent uuid|
|data.opration_seq|operation seq|


## Agent config get/update

### Request

    api/v1/agent/{agent_id}/configSync

### Method

    GET


### Params

| Param | Type | Description |
| ---- | ----| ---- |
|agent_id|string|agent uuid|

### Response

    //ok response
    {
    "msg": "ok",
    "code": 0,
    "data": {
        "agent_id": "****************",
        "opration_seq": 12
        }
    }

    //wrong response

    {
    "msg": "DB Model error",
    "code": 102,
    "data": "sql: no rows in result set"
    }


> Description

|Value|Type|Description|
|--------|--------|--------|
|data.agent_id|string|agent uuid|
|data.config_content|string|agent configuration content|　　


##  agent configuration update

### Request

    api/v1/agent/{agent_id}/configSync

### Method

    POST


### Params

| Param | Type | Description |
| ---- | ----| ---- |
|agent_id|string|agent uuid|
|config_content|string|agent configuration context|

### Response

    //#ok response
    {
    "msg": "ok",
    "code": 0,
    "data": {
        "agent_id": "****************",
        "opration_seq": 12
        }
    }

    //#wrong response

    {
    "msg": "DB Model error",
    "code": 102,
    "data": "sql: no rows in result set"
    }

> Description

|Value|Type|Description|
|--------|--------|--------|
|data.agent_id|string|agent uuid|
|data.opration_seq|operation seq|


## Exec script on the target sidecar(host)


### Request

	api/v1/sidecar/{sidecar_id:uuid}/execscriptSync

### Method

	POST

### Params


| Param | Type | Description |
| ---- | ---- | ---- |
|execScript|string|the content of the script|
|parameter|string|Params，separated by ＇,＇|
|timeout|string|A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".|


### Response

    //ok response
    {
        "msg": "ok",
        "code": 0,
        "data": {
            "operation_seq": 51662,
            "result": {
                "seqno": 51662,
                "response": "app\nbin\nboot\ndata\ndev\netc\nhome\nlib\nlib64\nlost+found\nmedia\nmnt\nogg\nopt\noracle\nproc\nroot\nrun\nsbin\nselinux\nsrv\nsys\ntmp\nusr\nvar\n"
            },
            "sidecar_id": "4a4b3343-0839-4e8d-9b9e-2ae0c7384644"
        }
    }

    //wrong response

    {
        "msg": "ok",
        "code": 0,
        "data": {
            "operation_seq": 51672,
            "result": {
                "seqno": 51672,
                "failed": true,
                "response": "fork/exec /tmp/exec-script-077126330: no such file or directory"#base64 encoding when SyncBase64
            },
            "sidecar_id": "4a4b3343-0839-4e8d-9b9e-2ae0c7384644"
        }
    }

> Description

|Value|Type|Description|
|--------|--------|--------|
|data.agent_id|string|agent uui|
|data.operation_seq|int||　


