<a name="29c80db5"></a>
## development

- Golang: >=1.12
- OS：linux\unix\windows
- Depends: [protoc](https://github.com/protocolbuffers/protobuf/releases/tag/v3.17.1), [go-bindata](https://github.com/go-bindata/go-bindata)


<a name="dd96ac2a"></a>
## build

just for linux\windows\darwin
```
make all 
```

## clear

just for linux\windows\darwin
```
make clean
```

## run server
- prepare the config file 
- prepare mysql and run the init sql;
- run the easyagent server with the cmd below
```
./easy-agent-server -c example-config.yml --debug
```

## install agent
- prepare the install curl cmd
- exec the curl cmd on the target host
- check the sidecar process
```
#change 127.0.0.1 to server ip
curl 'http://127.0.0.1:8889/api/v1/deploy/sidecar/install/shell?TargetPath=/opt/dtstack/easymanager/easyagent&CallBack=aHR0cDovLzE3Mi4xNi4xMC4zNzo4ODY0L2FwaS92Mi9hZ2VudC9pbnN0YWxsL2NhbGxiYWNrP2FpZD0tMQ==' | sh
#check sidecar process
ps -ef|grep sidecar

```