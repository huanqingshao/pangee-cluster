FROM ubuntu:22.04

ADD .devcontainer/docker/pip.conf /root/.pip/pip.conf

ENV LANG=C.UTF-8
ENV TZ=Asia/Shanghai

RUN apt-get update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    libssl-dev sshpass apt-transport-https moreutils git wget skopeo \
    ca-certificates curl gnupg2 python3-pip unzip rsync tzdata \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://github.com/mikefarah/yq/releases/download/v4.47.1/yq_linux_amd64 -O /usr/local/bin/yq \
    && chmod +x /usr/local/bin/yq

RUN update-alternatives --install /usr/bin/python python /usr/bin/python3 1
RUN mkdir /pangee-cluster 

WORKDIR /pangee-cluster
COPY ./.devcontainer/docker/requirements.txt ./requirements.txt
RUN python3 --version
RUN /usr/bin/python3 -m pip install --no-cache-dir pip -U \
    && python3 -m pip install --no-cache-dir -r requirements.txt

COPY .devcontainer/docker/ansible-patch/config/base.yml /usr/local/lib/python3.10/dist-packages/ansible/config/base.yml
COPY .devcontainer/docker/ansible-patch/plugins_callback/default.py /usr/local/lib/python3.10/dist-packages/ansible/plugins/callback/default.py
COPY .devcontainer/docker/ansible-patch/plugins_callback/__init__.py /usr/local/lib/python3.10/dist-packages/ansible/plugins/callback/__init__.py
COPY .devcontainer/docker/ansible-patch/plugins_action/raw.py /usr/local/lib/python3.10/dist-packages/ansible/plugins/action/raw.py

ENV PANGEE_CLUSTER_WEB_DIR="/pangee-cluster/ui"
ENV PANGEE_CLUSTER_PORT="10080"
ENV GIN_MODE=release
ENV PANGEE_CLUSTER_LOGRUS_LEVEL="info"
ENV PANGEE_CLUSTER_ADMIN_LOGRUS_LEVEL="info"
ENV PANGEE_CLUSTER_DEFAULT_ADMIN_PASSWORD="2432612431302464617343536e4858706f567278764f532f5274752f755a526c72694f30306d73456657674c7a39644b712f43436c6d64544e506375"
ENV TERM=xterm

COPY --chmod=+x ./admin/pangee-cluster-admin pangee-cluster-admin
COPY ./server/ansible-script ansible-script
COPY ./server/ansible-rpc ansible-rpc
COPY --chmod=+x ./server/pangee-cluster pangee-cluster
COPY ./web/dist /pangee-cluster/ui
COPY --chmod=+x ./server/pull-resource-package.sh pull-resource-package.sh

ENTRYPOINT ["./pangee-cluster"]
