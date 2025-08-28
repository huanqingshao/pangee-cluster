#!/bin/bash

# $1: 文件夹路径
# $2: tag/文件名
# $3: 下载源 (github.com / gitee.com)
# $4: 是否使用 proxy (true / false)
# $5: proxy 地址

echo "TASK [创建文件夹] ****"

mkdir -p $1/content

echo ""

echo "TASK [下载资源包] ****"

if [ "$3" != "null" ]; then
    # TODO: 修改地址
    wget -O $1/resource-pack.zip https://$3/Horan-Z/test/releases/download/$2
else
    echo "skip"
fi

echo ""

echo "TASK [解压资源包] ****"

unzip $1/resource-pack.zip -d $1/content
cp $1/content/package.yaml $1/package.yaml

echo ""

echo "TASK [配置代理] ****"

ENABLE_PROXY_ON_DOWNLOAD=$4         /usr/local/bin/yq -i '.all.hosts.localhost.enable_proxy_on_download = env(ENABLE_PROXY_ON_DOWNLOAD)' $1/content/inventory.yaml
if [ "$4" = "true" ]; then
    echo "HTTP_PROXY=$5"
    HTTP_PROXY=$5       /usr/local/bin/yq -i '.all.hosts.localhost.http_proxy = env(HTTP_PROXY)' $1/content/inventory.yaml
    HTTPS_PROXY=$5      /usr/local/bin/yq -i '.all.hosts.localhost.https_proxy = env(HTTPS_PROXY)' $1/content/inventory.yaml
else
    echo "No proxy"
fi

echo ""

echo "TASK [执行资源下载脚本] ****"

bash $1/content/download-dependency.sh

echo ""

echo "TASK [调整路径] ****"

PACKAGE_YAML="$1/content/package.yaml"
if [ ! -f "$PACKAGE_YAML" ]; then
    echo "Error: package.yaml not found in zip file"
    exit 1
fi

YAML_VERSION=$(/usr/local/bin/yq e '.metadata.version' "$PACKAGE_YAML")
DIR_NAME=$(basename "$1")
if [ "$YAML_VERSION" != "$DIR_NAME" ]; then
    echo "Version mismatch: YAML version=$YAML_VERSION, Directory version=$DIR_NAME"
    
    # 定义目标路径
    PARENT_DIR=$(dirname "$1")
    TARGET_DIR="$PARENT_DIR/$YAML_VERSION"
    
    # 删除目标文件夹(如果存在)
    if [ -d "$TARGET_DIR" ]; then
        echo "Removing existing target directory: $TARGET_DIR"
        rm -rf "$TARGET_DIR"
    fi
    
    # 重命名当前文件夹
    echo "Renaming $1 to $TARGET_DIR"
    mv "$1" "$TARGET_DIR"
else
    echo "Version matches: $YAML_VERSION, no action needed"
fi