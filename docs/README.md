---
sidebarDepth: 0
---

# PangeeCluster

中国移动开源的图形化的 K8S 集群离线安装、维护工具。

## 快速安装

找一台不低于 1 核 2G，不少于 10G 剩余磁盘空间，已经安装好 docker 的服务器，执行如下指令，即可完成 PangeeCluster 的安装：

[【腾讯云】云产品限时秒杀，爆款 1 核 2G 云服务器，首年 99 元](https://cloud.tencent.com/act/cps/redirect?redirect=1062&cps_key=2ee6baa049659f4713ddc55a51314372&from=console)

```sh
docker run -d \
  --privileged \
  --restart=unless-stopped \
  --name=pangee-cluster \
  -p 80:80/tcp \
  -e TZ=Asia/Shanghai \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ~/pangee-cluster-data:/data \
  opencmit/pangee-cluster:latest-amd64
  # 如果您是 arm64 环境，请将标签里的 amd64 修改为 arm64，例如 opencmit/pangee-cluster:latest-arm64
```

在浏览器地址栏中输入 `http://这台机器的IP地址`，输入用户名 `admin`，默认密码 `PangeeCluster123`，即可登录 pangee-cluster 界面，剩下的事情，在界面上操作一下，您就会啦。如果有困难，试试这篇文档 [使用 PangeeCluster 安装 Kubernetes 集群](./guide/install-k8s.md)

::: tip 常见问题

- PangeeCluster 的信息保存在容器的 `/data` 路径，请将其映射到一个您认为安全的地方，上面的命令中，将其映射到了 `~/pangee-cluster-data` 路径；
- 只要此路径的内容不受损坏，重启、升级、重新安装 Pangee-Cluster，或者将数据及 Pangee-Cluster 迁移到另外一台机器上，您都可以找回到原来的信息；
- 加个 [GITHUB Star](https://github.com/opencmit/pangee-cluster)，避免迷路。

:::

## 自制资源包

Pangee-Cluster 将定期提供最新版本的资源包，可以在 kubord-spray 资源包管理界面中查到，如果您是离线环境，也可以在 [https://pangee-cluster.cn/support/](./support/)找到最新的资源包。您也可以自己制作资源包，资源包的项目地址在 [pangee-cluster-resource](https://github.com/opencmit/pangee-cluster-resource)。

## 社区

