#!/bin/bash

# $1: 文件夹路径
# $2: 资源包名称
# $3: 下载源 (github.com / gitee.com)
# $4: 是否使用 proxy (true / false)
# $5: proxy 地址

echo "------------------------------------"
echo "$1"
echo "------------------------------------"
echo "$2"
echo "------------------------------------"
echo "$3"
echo "------------------------------------"
echo "$4"
echo "------------------------------------"
echo "$5"
echo "------------------------------------"

echo "TASK [创建文件夹] ****"

mkdir -p $1/content

echo ""

echo "TASK [下载资源包] ****"

if [ "$3" != "null" ]; then
    git clone -b $2 https://$3/Horan-Z/test.git $1/content
else
    echo "skip"
fi

echo ""

echo "TASK [配置代理] ****"

ENABLE_PROXY=$4         /root/go/bin/yq -i '.all.hosts.localhost.enable_proxy = env(ENABLE_PROXY)' $1/content/inventory.yaml
if [ "$4" = "true" ]; then
    echo "HTTP_PROXY=$5"
    HTTP_PROXY=$5       /root/go/bin/yq -i '.all.hosts.localhost.http_proxy = env(HTTP_PROXY)' $1/content/inventory.yaml
    HTTPS_PROXY=$5      /root/go/bin/yq -i '.all.hosts.localhost.https_proxy = env(HTTPS_PROXY)' $1/content/inventory.yaml
else
    echo "No proxy"
fi

echo ""

echo "TASK [执行资源下载脚本] ****"

bash $1/content/3.sh

echo ""

echo "PLAY RECAP *********************************************************************"
echo "pangeecluster : ok=4    unreachable=0    failed=0"