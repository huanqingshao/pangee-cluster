#!/bin/bash

datetime=`date "+%Y-%m-%d %H:%M:%S"`

echo $datetime

arch="amd64"
if [ $(uname -m) != "x86_64" ]; then
  arch="arm64"
fi

echo "【构建】 ${1}-${arch}"

tag=eipwork/pangee-cluster
tag_backup=swr.cn-east-2.myhuaweicloud.com/kuboard/pangee-cluster

echo
echo "【构建 web】"

cd web
yarn install
yarn build

echo \{\"version\":\"${1}-${arch}\",\"buildDate\":\"$datetime\"\} > ./dist/version.json
cd ..

echo
echo "【构建 server】"
rm -f ./server/pangee-cluster
docker run --rm -v ${PWD}:/usr/src/pangee-cluster \
  -v ~/temp/build-temp/pkg:/go/pkg \
  -w /usr/src/pangee-cluster/server golang:1.18.0-buster \
  sh -c "export GOPROXY=https://goproxy.io,direct && go build pangee-cluster.go"

ls -lh ./server/pangee-cluster

echo
echo "【构建 admin】"
rm -f ./admin/pangee-cluster-admin
docker run --rm -v ${PWD}:/usr/src/pangee-cluster \
  -v ~/temp/build-temp/pkg:/go/pkg \
  -w /usr/src/pangee-cluster/admin golang:1.18.0-buster \
  sh -c "export GOPROXY=https://goproxy.io,direct && go build pangee-cluster-admin.go"

ls -lh ./admin/pangee-cluster-admin

echo
echo "【构建 镜像】"

docker build -f Dockerfile -t $tag:$1-${arch} --build-arg arch=${arch} .

echo "【构建 成功】"
echo $tag:$1-${arch}

docker tag $tag:$1-${arch} $tag:latest-${arch}
docker tag $tag:$1-${arch} $tag_backup:$1-${arch}
docker tag $tag:$1-${arch} $tag_backup:latest-${arch}

if [ "$2" == "" ]
then

  docker push $tag:$1-${arch}
  docker push $tag:latest-${arch}
  docker push $tag_backup:$1-${arch}
  docker push $tag_backup:latest-${arch}

else

  echo
  echo "【稍后推送镜像】"

  echo "#!/bin/bash" > push.sh
  echo "echo $datetime" >> push.sh
  echo "docker push $tag:$1-${arch}" >> push.sh
  echo "docker push $tag:latest-${arch}" >> push.sh
  echo "docker push $tag_backup:$1-${arch}" >> push.sh
  echo "docker push $tag_backup:latest-${arch}" >> push.sh
  echo "echo $datetime" >> push.sh

  echo "【已生成 push.sh】"

fi

echo $datetime