# Easyaddons

![](addon.png)

```text
亚当是《圣经·旧约》创世纪篇章中的人物。
根据记载，他是世上的第一个人类与第一个男人，是耶和华按照自己的形象用尘土而造的活人。
亚当的毕生事迹记载于创世记2~5章。
耶和华神用五天时间创造了天地万物，一切准备好了之后，第六天，造了亚当。
在创世记2章7节，这里说到：耶和华神用地上的尘土造人，将生气吹在他鼻孔里，他就成了有灵的活人，名叫亚当。
```

## 是什么？

亚当可以是所有通过EM安装下发的插件。
插件和产品包的区别，插件是为了满足EM的某种私欲而产生的，是为了更好的服务于通过EM部署的产品。
需提供如下运行规范：
```json
{
    "name":"easyfiler",
    "version": "v1.0.0",
    "type":"rpc",
    "binaryPath":"easyfiler/sbin/easyfiler",
    "configuraionPaths": ["easyfiler/conf/config.yml"],
    "configTpl": "",
    "logs": ["easyfiler/logs","/tmp/logs"],
    "parameter":"--config.file,easyfiler/conf/config.yml",
    "extendInfo": {"port": 8787}
}
```

|  字段   | 选项  |
|  :----  | ----  |
| name  | 必须 |
| version  | 必须 |
| type  | 必须 |
| binaryPath  | 必须 |
| configuraionPaths  | 可选 |
| configTemplate  | 可选 |
| logs  | 可选 |
| parameter  | 可选 |
| extendInfo  | 可选 |

## 能做什么？

```text
以不变应万变
EM能力增强，使EM对主机端的控制具备开放扩展能力；
```