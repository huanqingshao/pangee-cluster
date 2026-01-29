---
description: PangeeCluster 使用图形界面离线安装 Kubernetes 高可用集群
sidebarDepth: 1
meta:
  - name: keywords
    content: 配置 PangeeCluster 开发测试环境
---

# 配置开发测试环境

本文描述了如何配置 PangeeCluster 的开发测试环境，并编译和构建 PangeeCluster。本文的目标读者为想要尝试对 PangeeCluster 进行二次开发的开发者。

## 开发环境要求

开发环境的配置要求不应低于：

- CPU 4 核
- 内存 8G
- 磁盘空间 60G，因为经常构建镜像，推荐磁盘空间不小于 120G

开发环境所需要的软件列表如下：

- Docker 20.10.17 或更高版本
  - https://docs.docker.com/engine/install/ubuntu/#install-using-the-convenience-script
    ```sh
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh ./get-docker.sh
    # 检查安装是否成功
    docker version
    ```

## 进入开发环境

### 导入代码到开发环境

- 执行如下命令，将代码 clone 到开发环境。

  ```sh
  git clone https://github.com/opencmit/pangee-cluster.git
  cd pangee-cluster
  git submodule update --init
  ```

- 使用 VS code 在容器中打开工作区

- 替换对 `ansible` 定制的文件：（待完善）

  ```sh
  cp .docker/ansible-patch/config/base.yml /usr/local/lib/python3.8/dist-packages/ansible/config/base.yml
  cp .docker/ansible-patch/plugins_callback/default.py /usr/local/lib/python3.8/dist-packages/ansible/plugins/callback/default.py
  cp .docker/ansible-patch/plugins_callback/__init__.py /usr/local/lib/python3.8/dist-packages/ansible/plugins/callback/__init__.py
  cp .docker/ansible-patch/plugins_action/raw.py /usr/local/lib/python3.8/dist-packages/ansible/plugins/action/raw.py
  ```

## 运行开发环境

需按如下步骤逐个运行各组件

- 运行 pangee-cluster/server
- 运行 pangee-cluster/web

### 运行 pangee-cluster server

- 在 vscode 中导航到 `pangee-cluster/web`，点击右键，并点击 `在集成终端中打开`，如下图所示

  ![web](./dev.assets/screenshot_2025-10-04_15.57.13.png)

- 在集成终端中（pangee-cluster/web 路径下）执行命令：

  ```sh
  pnpm install
  pnpm build
  ```

- 在 vscode 中导航到 `pangee-cluster`，点击右键，并点击 `在集成终端中打开`，省略截图；

- 在集成终端中（pangee-cluster 路径下）执行命令：

  ```sh
  mkdir -p data/user
  ```

- 在 vscode 中导航到 `pangee-cluster/server`，点击右键，并点击 `在集成终端中打开`，省略截图；

- 在集成终端中（pangee-cluster/server 路径下）执行命令：

  ```sh
  go run pangee-cluster.go
  ```

  如果执行成功，最后的输出日志如下所示：

  ```sh
  # 如果是 amd64 环境
  [GIN-debug] Listening and serving HTTP on :9080
  # 如果是 arm64 环境
  [GIN-debug] Listening and serving HTTP on :9081
  ```

### 运行 pangee-cluster web

- 在 vscode 中导航到 `pangee-cluster/web`，点击右键，并点击 `在集成终端中打开`，省略截图；

- 编辑 `/etc/hosts` 文件，添加如下信息：

  ```sh
  192.168.64.68  kb pangee-cluster-arm
  # 假设 192.168.64.68 为您开发环境的 IP 地址
  ```

- 在集成终端中（pangee-cluster/web/public 路径下）执行命令：

  ```sh
  ln -s version-$(dpkg --print-architecture).json version.json
  # 仅需在首次运行时执行一次即可
  ```

- 在集成终端中（pangee-cluster/web 路径下）执行命令：

  ```sh
  pnpm dev
  ```

- 在浏览器打开如下路径

  `http://localhost:8848`

- 在登录界面中输入默认用户名 `admin`，默认密码 `PangeeCluster123`，可登录到 PangeeCluster 的界面。

## 构建容器镜像

构建容器镜像时，只需要在开发环境的 `pangee-cluster` 目录中执行如下命令

```sh
./build.sh v2.0.0
# 其中 v2.0.0 为本次构建时的版本号
```

## 构建二进制可执行文件

在开发容器所在宿主机安装 dockerc，并执行 dockerc 构建命令

```shell
wget https://github.com/NilsIrl/dockerc/releases/download/v0.3.2/dockerc_x86-64
mv dockerc_x86-64 dockerc
chmod +x dockerc
./dockerc --image docker-daemon:pangee-cluster:v2.0.0-amd64 --output pangee-cluster-bin
```

## 发布资源包

将打包好的资源包 zip 上传至 github release 并选定 tag 即可，余下步骤由 github action 自动完成

**请确保压缩包文件名(去除.zip)与 package.yaml 中的version名一致**

**如不一致，用户手动下载压缩包然后加载本地资源包时会出现问题，用户直接在线下载资源包不会出现问题**