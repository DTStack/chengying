## messages define

<a name="Register"></a>
### Register

```
// Register

message RegisterRequest {
    SidecarRequestHeader H = 1 [(gogoproto.nullable) = false, (gogoproto.embed) = true];

    string osType = 2;
    string osPlatform = 3;
    string osVersion = 4;
    string cpuSerial = 5;
    uint32 cpuCores = 6;
    uint64 memSize = 7;
    uint64 swapSize = 8;
    string host = 9;
    string localIp = 10;
    repeated string callBack = 11;
}

message RegisterResponse {

}
```

<a name="Controlling"></a>
### Controlling

```
// Controlling
message ControlRequest {
    SidecarRequestHeader H = 1 [(gogoproto.nullable) = false, (gogoproto.embed) = true];
}

message ControlResponse {

    enum ControlCmd {
        INSTALL_AGENT = 0;
        UNINSTALL_AGENT = 1;
        UPDATE_AGENT = 2;
        START_AGENT = 3;
        STOP_AGENT = 4;
        UPDATE_AGENT_CONFIG = 5;
        UPDATE_SIDECAR = 6;
        UPDATE_SIDECAR_CONFIG = 7;
        EXEC_SCRIPT = 8;
        CANCEL_AGENTS = 9;
        EXEC_REST = 10;
    }

    ControlCmd cmd = 1;

    uint32 seqno = 2;

    oneof options {
        InstallAgentOptions installAgentOptions = 100;
        UninstallAgentOptions uninstallAgentOptions = 101;
        UpdateAgentOptions updateAgentOptions = 102;
        StartAgentOptions startAgentOptions = 103;
        StopAgentOptions stopAgentOptions = 104;
        UpdateAgentConfigOptions updateAgentConfigOptions = 105;
        UpdateSidecarOptions updateSidecarOptions = 106;
        UpdateSidecarConfigOptions updateSidecarConfigOptions = 107;
        ExecScriptOptions execScriptOptions = 108;
        CancelOptions cancelOptions = 109;
        ExecRestOptions execRestOptions = 110;
    }

    message InstallAgentOptions {
        bytes agentId = 1 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
        string configurationPath = 2;
        string binaryPath = 3;
        repeated string parameter = 4;
        string name = 5;

        repeated string installParameter = 6;
        string installScript = 7;
        string workdir = 8;

        string healthShell = 9;
        google.protobuf.Duration healthPeriod = 10 [(gogoproto.stdduration) = true, (gogoproto.nullable) = false];

        google.protobuf.Duration timeout = 11 [(gogoproto.stdduration) = true, (gogoproto.nullable) = false];
        google.protobuf.Duration healthStartPeriod = 12 [(gogoproto.stdduration) = true, (gogoproto.nullable) = false];
        google.protobuf.Duration healthTimeout = 13 [(gogoproto.stdduration) = true, (gogoproto.nullable) = false];
        uint64 healthRetries = 14;
        string runUser = 15;
    }

    message UninstallAgentOptions {
        bytes agentId = 1 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
        repeated string parameter = 2;
        string uninstallScript = 3;

        google.protobuf.Duration timeout = 4 [(gogoproto.stdduration) = true, (gogoproto.nullable) = false];
    }

    message UpdateAgentOptions {
        bytes agentId = 1 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
        repeated string parameter = 2;
        string updateScript = 3;

        google.protobuf.Duration timeout = 4 [(gogoproto.stdduration) = true, (gogoproto.nullable) = false];
    }

    message StartAgentOptions {
        bytes agentId = 1 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
        float cpuLimit = 2;
        uint64 memLimit = 3;
        uint64 netLimit = 4;
        map<string, string> environment = 5;
    }

    message StopAgentOptions {
        bytes agentId = 1 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];

        enum StopAgentOptionsType {
            STOP_UNRECOVER = 0;
            STOP_RECOVER = 1;
        }

        StopAgentOptionsType stopAgentOptionsType = 2;
    }

    message UpdateAgentConfigOptions {
        bytes agentId = 1 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
        string configContent = 2;
        string configPath = 3;
    }

    message UpdateSidecarOptions {
    }

    message UpdateSidecarConfigOptions {
    }

    message ExecScriptOptions {
        repeated string parameter = 1;
        string execScript = 2;
        google.protobuf.Duration timeout = 3 [(gogoproto.stdduration) = true, (gogoproto.nullable) = false];
        bytes agentId = 4 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
    }

    message CancelOptions {
        bytes agentId = 1 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
    }

    message ExecRestOptions {
        string method = 1;
        string path = 2;
        string query = 3;
        bytes body = 4;
        google.protobuf.Duration timeout = 5 [(gogoproto.stdduration) = true, (gogoproto.nullable) = false];
        bytes agentId = 6 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
    }
}

```

<a name="Event"></a>
### Event

```
message Event {
    SidecarRequestHeader H = 1 [(gogoproto.nullable) = false, (gogoproto.embed) = true];

    enum EventType {
        EVT_OP_PROGRESS = 0;
        EVT_AGENT_ERR = 1;
        EVT_OS_RESOURCE_USAGES = 2;
        EVT_PROC_RESOURCE_USAGES = 3;
        EVT_EXEC_SCRIPT = 4;
        EVT_AGENT_HEALTH_CHECK = 5;
        EVT_REST_SCRIPT = 6;
    }

    EventType eventType = 2;

    oneof details {
        OperationProgress opProgress = 100;
        AgentError agentError = 101;
        OsResourceUsages osResourceUsages = 102;
        ProcessResourceUsages procResourceUsages = 103;
        ExecScriptResponse execScriptResponse = 104;
        AgentHealthCheck agentHealthCheck = 105;
        ExecRestResponse execRestResponse = 106;
    }

    message OperationProgress {
        uint32 seqno = 1;
        bytes agentId = 2 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
        float progress = 3;
        bool failed = 4;
        string message = 5;
    }

    message AgentError {
        bytes agentId = 1 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
        string errstr = 2;
    }

    message OsResourceUsages {
        float cpuUsage = 1;
        uint64 memUsage = 2;
        uint64 swapUsage = 3;
        float load1 = 4;
        double uptime = 5;
        repeated DiskUsage diskUsage = 6 [(gogoproto.nullable) = false];
        repeated NetStat netStats = 7 [(gogoproto.nullable) = false];
    }

    message DiskUsage {
        string mountPoint = 1;
        uint64 usedSpace = 2;
        uint64 totalSpace = 3;
    }

    message NetStat {
        string ifName = 1;
        repeated string ifIp = 2;
        uint64 bytesSent = 3;
        uint64 bytesRecv = 4;
    }

    message ProcessResourceUsages {
        bytes agentId = 1 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
        float cpuUsage = 2;
        uint64 memory = 3;
        string cmd = 4;
        uint64 bytesSent = 5;
        uint64 bytesRecv = 6;
    }

    message ExecScriptResponse {
        uint32 seqno = 1;
        bool failed = 2;
        string response = 3;
        bytes agentId = 4 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
    }

    message AgentHealthCheck {
        bytes agentId = 1 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
        bool failed = 2;
    }

    message ExecRestResponse {
            uint32 seqno = 1;
            bool failed = 2;
            bytes response = 3;
            bytes agentId = 4 [(gogoproto.customtype) = "Uuid", (gogoproto.nullable) = false];
    }
}

```