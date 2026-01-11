---
---

# 制作资源包

## 资源包的代码也是开源的

[pangee-cluster-resource](https://github.com/opencmit/pangee-cluster-resource)

## 资源包的版本编码规则

资源包的版本号由如下几部分组成：

## 资源包制作方法

### 第一步 下载资源包模板

1. 克隆资源包 GitHub 仓库

```shell
# 只克隆保存模板的 main 分支
# TODO: 修改地址
git clone -b main --single-branch https://github.com/Horan-Z/test.git
cd pangee-cluster-resource-package
```

### 第二步 制作并测试资源包

1. 修改 package.yaml 内容：
   * 资源版本号（无需修改 checksum ）
   * metadata 内容
   * 修改并复制 metadata 中的资源包版本号，以下用版本号代替

2. 执行 create-package.sh 获取资源文件以及对应 checksum

```shell
./create-package.sh
```

3. 打包资源包为 .zip 文件
```shell
zip -r <版本号>.zip . -x "cache/**\*" "temp/**\*"
```

4. 在 Pangee Cluster 中通过以下路径加载离线资源包

   系统设置 -> 资源包管理 -> 加载离线资源包

   注意加载离线资源包时所填代理信息是 Pangee Cluster 后端可以使用的代理，而不是访问网页端的机器所使用的代理

### 第三步 测试资源包
1. 测试集群上检测资源包是否可用

2. 如资源包不可用（资源版本冲突、资源文件损坏等），重新执行第二步

### 第四步 添加资源包到官方仓库（可选）

1. 在官方仓库的 **Issues** 板块创建新issue，标题建议格式为 **申请添加资源包 - [版本号]**
2. 在issue内容中明确提供以下信息：
   - 资源包版本号
   - 资源包是否经过充分测试
3. 仓库管理员收到issue后，会为你创建一个与版本号同名的分支，并在issue中回复通知
4. 你需基于该分支提出 Pull Request，注意将base分支设置为管理员新创建的版本号分支
5. 管理员审核通过并merge后，即完成资源包添加