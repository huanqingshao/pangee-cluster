FROM ubuntu:24.04

RUN sed -i 's/#force_color_prompt=yes/force_color_prompt=yes/g' /root/.bashrc \
  && sed -i 's/ports\.ubuntu\.com/mirrors\.tuna\.tsinghua\.edu\.cn/g' /etc/apt/sources.list.d/ubuntu.sources\
  && apt update \
  && DEBIAN_FRONTEND=noninteractive apt install -y --no-install-recommends tzdata \
  && apt upgrade -y \
  && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && apt install curl git vim wget -y 

RUN install -m 0755 -d /etc/apt/keyrings \
  && curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc \
  && chmod a+r /etc/apt/keyrings/docker.asc  \
  && echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu noble stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null \
  && apt update \
  && apt install docker-ce-cli -y

RUN curl -fsSL https://deb.nodesource.com/setup_22.x -o nodesource_setup.sh \
  && bash nodesource_setup.sh \
  && apt-get install -y nodejs \
  && useradd -D

RUN wget https://go.dev/dl/go1.24.1.linux-$(dpkg --print-architecture).tar.gz \
  && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.1.linux-$(dpkg --print-architecture).tar.gz

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

RUN apt install python3.12-venv -y

ADD requirements.txt /root/requirements.txt

SHELL ["/bin/bash", "-c"] 

RUN python3 -m venv /root/.cmcc-venv \
  && source /root/.cmcc-venv/bin/activate \
  && pip install -U -r /root/requirements.txt

