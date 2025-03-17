---
description: KuboardSpray 使用图形界面离线安装 Kubernetes 高可用集群
sidebarDepth: 1
meta:
  - name: keywords
    content: 配置 KuboardSpray 开发测试环境
---

# 配置开发测试环境

本文描述了如何配置 KuboardSpray 的开发测试环境，并编译和构建 KuboardSpray。本文的目标读者为想要尝试对 KuboardSpray 进行二次开发的开发者。

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
  git clone https://github.com/eip-work/kuboard-spray.git
  cd kuboard-spray
  git submodule update --init
  ```

- 使用 VS code 在容器中打开工作区

- 在 vscode 中导航到 `kuboard-spray/requirements`，点击右键，并点击 `在集成终端中打开`，如下图所示

  ```sh
  # 当前路径为 kuboard-spray
  pip install -r requirements.txt --break-system-packages
  ```

- 替换对 `ansible` 定制的文件：（待完善）

  ```sh
  cp .docker/ansible-patch/config/base.yml /usr/local/lib/python3.8/dist-packages/ansible/config/base.yml
  cp .docker/ansible-patch/plugins_callback/default.py /usr/local/lib/python3.8/dist-packages/ansible/plugins/callback/default.py
  cp .docker/ansible-patch/plugins_callback/__init__.py /usr/local/lib/python3.8/dist-packages/ansible/plugins/callback/__init__.py
  cp .docker/ansible-patch/plugins_action/raw.py /usr/local/lib/python3.8/dist-packages/ansible/plugins/action/raw.py
  ```

## 运行开发环境

需按如下步骤逐个运行各组件

- 运行 kuboard-spray server
- 运行 kuboard-spray web

### 运行 kuboard-spray server

- 在 vscode 中导航到 `kuboard-spray/web`，点击右键，并点击 `在集成终端中打开`，如下图所示

  ![web](./dev.assets/iShot_2022-08-06_20.21.03.png)

- 在集成终端中（kuboard/web 路径下）执行命令：

  ```sh
  pnpm install
  pnpm build
  ```

- 在 vscode 中导航到 `kuboard-spray`，点击右键，并点击 `在集成终端中打开`，省略截图；

- 在集成终端中（kuboard-spray 路径下）执行命令：

  ```sh
  mkdir -p data/user
  ```

- 在 vscode 中导航到 `kuboard-spray/server`，点击右键，并点击 `在集成终端中打开`，省略截图；

- 在集成终端中（kuboard-spray/server 路径下）执行命令：

  ```sh
  go run kuboard-spray.go
  ```

  如果执行成功，最后的输出日志如下所示：

  ```sh
  # 如果是 amd64 环境
  [GIN-debug] Listening and serving HTTP on :8006
  # 如果是 arm64 环境
  [GIN-debug] Listening and serving HTTP on :8007
  ```

### 运行 kuboard-spray web

- 在 vscode 中导航到 `kuboard-spray/web`，点击右键，并点击 `在集成终端中打开`，省略截图；

- 编辑 `/etc/hosts` 文件，添加如下信息：

  ```sh
  192.168.64.68  kb kuboard-spray-arm
  # 假设 192.168.64.68 为您开发环境的 IP 地址
  ```

- 在集成终端中（kuboard-spray/web/public 路径下）执行命令：

  ```sh
  ln -s version-$(dpkg --print-architecture).json version.json
  # 仅需在首次运行时执行一次即可
  ```

- 在集成终端中（kuboard-spray/web 路径下）执行命令：

  ```sh
  pnpm serve
  ```

- 在浏览器打开如下路径

  `http://localhost:25702`（如果是 amd64 CPU）
  或 `http://localhost:25703` （如果是 arm64 CPU）。

- 在登录界面中输入默认用户名 `admin`，默认密码 `Kuboard123`，可登录到 KuboardSpray 的界面。

## 构建容器镜像

构建容器镜像时，只需要在开发环境的 `kuboard-spray` 目录中执行如下命令

```sh
./build.sh v1.2.0 no_push
# 其中 v1.2.0 为本次构建时的版本号，no_push 参数代表只生成镜像，不推送到镜像仓库
```
