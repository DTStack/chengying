
## 用法

```text
go build

离线化rke k8s部署镜像
./easyrke --config example-config.yml --k8s v1.16.3-rancher1-1

package em2.0 k8s 产品包
cd package
sh  -x package.sh  v1.12.9-rancher1-1
package的版本号要跟离线镜像版本号一致

Product package create success: "DTK8S_DTK8S-v1.12.9-rancher1-1.tar"
```