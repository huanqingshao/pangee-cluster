FROM ubuntu:22.04

ADD .docker/pip.conf /root/.pip/pip.conf

ENV LANG=C.UTF-8
ENV TZ=Asia/Shanghai

RUN apt-get update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    libssl-dev sshpass apt-transport-https moreutils git wget skopeo \
    ca-certificates curl gnupg2 python3-pip unzip rsync tzdata \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://github.com/mikefarah/yq/releases/download/v4.47.1/yq_linux_amd64 -O /usr/local/bin/yq \
    && chmod +x /usr/local/bin/yq

WORKDIR /pangee-cluster
COPY ./.docker/requirements.txt ./requirements.txt
RUN /usr/bin/python3 -m pip install --no-cache-dir pip -U \
    && python3 -m pip install --no-cache-dir -r requirements.txt \
    && update-alternatives --install /usr/bin/python python /usr/bin/python3 1

COPY .docker/ansible-patch/config/base.yml /usr/local/lib/python3.10/dist-packages/ansible/config/base.yml
COPY .docker/ansible-patch/plugins_callback/default.py /usr/local/lib/python3.10/dist-packages/ansible/plugins/callback/default.py
COPY .docker/ansible-patch/plugins_callback/__init__.py /usr/local/lib/python3.10/dist-packages/ansible/plugins/callback/__init__.py
COPY .docker/ansible-patch/plugins_action/raw.py /usr/local/lib/python3.10/dist-packages/ansible/plugins/action/raw.py

ENV PANGEE_CLUSTER_WEB_DIR="/pangee-cluster/ui"
ENV PANGEE_CLUSTER_PORT="8080"
ENV GIN_MODE=release
ENV PANGEE_CLUSTER_LOGRUS_LEVEL="info"
ENV PANGEE_CLUSTER_ADMIN_LOGRUS_LEVEL="info"
ENV TERM=xterm

COPY ./admin/pangee-cluster-admin pangee-cluster-admin
COPY ./server/ansible-script ansible-script
COPY ./server/ansible-rpc ansible-rpc
COPY ./server/pangee-cluster pangee-cluster
COPY ./web/dist /pangee-cluster/ui
COPY ./server/pull-resource-package.sh pull-resource-package.sh
RUN chmod +x pangee-cluster pangee-cluster-admin pull-resource-package.sh

ENTRYPOINT ["./pangee-cluster"]
