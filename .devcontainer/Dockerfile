FROM ubuntu:20.04

RUN sed -i 's/#force_color_prompt=yes/force_color_prompt=yes/g' /root/.bashrc \
  && sed -i 's/ports\.ubuntu\.com/mirrors\.tuna\.tsinghua\.edu\.cn/g' /etc/apt/sources.list \
  && sed -i 's/archive\.ubuntu\.com/mirrors\.tuna\.tsinghua\.edu\.cn/g' /etc/apt/sources.list \
  && apt update \
  && DEBIAN_FRONTEND=noninteractive apt install -y --no-install-recommends tzdata \
  && apt upgrade -y \
  && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && apt install curl git vim wget -y 

RUN wget https://download.docker.com/linux/ubuntu/dists/focal/pool/stable/$(dpkg --print-architecture)/docker-ce-cli_20.10.24~3-0~ubuntu-focal_$(dpkg --print-architecture).deb \
  && dpkg -i docker-ce-cli_20.10.24~3-0~ubuntu-focal_$(dpkg --print-architecture).deb \
  && rm docker-ce-cli_20.10.24~3-0~ubuntu-focal_$(dpkg --print-architecture).deb

RUN curl -fsSL https://deb.nodesource.com/setup_22.x -o nodesource_setup.sh \
  && bash nodesource_setup.sh \
  && apt-get install -y nodejs \
  && useradd -D

RUN wget https://go.dev/dl/go1.24.2.linux-$(dpkg --print-architecture).tar.gz \
  && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.2.linux-$(dpkg --print-architecture).tar.gz

RUN apt install -y python3-pip \
  && apt install -y sshpass

RUN apt install sudo -y \
  && echo "ubuntu ALL=(ALL:ALL) NOPASSWD: ALL" >> /etc/sudoers

RUN wget -qO- https://get.pnpm.io/install.sh | ENV="/root/.bashrc" SHELL="/usr/bin/bash" bash -

RUN echo "export PATH=$PATH:/usr/local/go/bin:/root/.local/bin" >> /root/.profile \
  && echo "export GO111MODULE=on" >> /root/.profile \
  && echo "export GOPROXY=https://goproxy.cn" >> /root/.profile \
  && echo "source /root/.cmcc-venv/bin/activate" >> /root/.profile \
  && echo "export LC_ALL=C.UTF-8" >> /root/.profile

RUN apt install python3-venv -y

ADD docker/requirements.txt /root/requirements.txt

SHELL ["/bin/bash", "-c"] 

RUN python3 -m venv /root/.cmcc-venv \
  && source /root/.cmcc-venv/bin/activate \
  && pip install -U -r /root/requirements.txt

COPY docker/ansible-patch/config/base.yml /root/.cmcc-venv/lib/python3.8/site-packages/ansible/config/base.yml
COPY docker/ansible-patch/plugins_callback/default.py /root/.cmcc-venv/lib/python3.8/site-packages/ansible/plugins/callback/default.py
COPY docker/ansible-patch/plugins_callback/__init__.py /root/.cmcc-venv/lib/python3.8/site-packages/ansible/plugins/callback/__init__.py
COPY docker/ansible-patch/plugins_action/raw.py /root/.cmcc-venv/lib/python3.8/site-packages/ansible/plugins/action/raw.py

RUN apt install -y rsync netcat

RUN wget https://github.com/mikefarah/yq/releases/download/v4.47.1/yq_linux_amd64 -O /usr/local/bin/yq \
    && chmod +x /usr/local/bin/yq

RUN apt install -y unzip

RUN source /root/.profile && go version && apt install -y libgpgme-dev libassuan-dev libbtrfs-dev pkg-config \
    && git clone https://github.com/containers/skopeo.git \
    && cd skopeo \
    && git checkout v1.21.0 \
    && make bin/skopeo \
    && mv bin/skopeo /usr/local/bin \
    && cd ../ \
    && rm -rf skopeo

# 创建与宿主机相同 UID/GID 的用户
ARG USER_UID=1000
ARG USER_GID=1000
RUN if [ "$USER_GID" != "1000" ]; then groupmod -g $USER_GID vscode; fi && \
    if [ "$USER_UID" != "1000" ]; then usermod -u $USER_UID vscode; fi

ENTRYPOINT [ "tail", "-f", "/dev/null" ]