FROM ubuntu:20.04

ARG HTTP_PROXY http://10.2.193.51:10811

RUN sed -i 's/#force_color_prompt=yes/force_color_prompt=yes/g' /root/.bashrc \
  && sed -i 's/ports\.ubuntu\.com/mirrors\.tuna\.tsinghua\.edu\.cn/g' /etc/apt/sources.list \
  && sed -i 's/archive\.ubuntu\.com/mirrors\.tuna\.tsinghua\.edu\.cn/g' /etc/apt/sources.list \
  && apt update \
  && DEBIAN_FRONTEND=noninteractive apt install -y --no-install-recommends tzdata \
  && apt upgrade -y \
  && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && apt install curl git vim wget -y 

RUN export https_proxy=${HTTP_PROXY} \
  && wget https://download.docker.com/linux/ubuntu/dists/focal/pool/stable/$(dpkg --print-architecture)/docker-ce-cli_20.10.24~3-0~ubuntu-focal_$(dpkg --print-architecture).deb \
  && dpkg -i docker-ce-cli_20.10.24~3-0~ubuntu-focal_$(dpkg --print-architecture).deb \
  && rm docker-ce-cli_20.10.24~3-0~ubuntu-focal_$(dpkg --print-architecture).deb

RUN curl -fsSL https://deb.nodesource.com/setup_22.x -o nodesource_setup.sh \
  && bash nodesource_setup.sh \
  && apt-get install -y nodejs \
  && useradd -D

RUN export https_proxy=${HTTP_PROXY} && wget https://go.dev/dl/go1.24.1.linux-$(dpkg --print-architecture).tar.gz \
  && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.1.linux-$(dpkg --print-architecture).tar.gz

RUN apt install -y python3-pip \
  && apt install -y sshpass

RUN apt install sudo -y \
  && echo "ubuntu ALL=(ALL:ALL) NOPASSWD: ALL" >> /etc/sudoers

RUN export https_proxy=${HTTP_PROXY} && wget -qO- https://get.pnpm.io/install.sh | ENV="/root/.bashrc" SHELL="/usr/bin/bash" bash -

RUN echo "export PATH=$PATH:/usr/local/go/bin:/root/.local/bin" >> /root/.profile \
  && echo "export GO111MODULE=on" >> /root/.profile \
  && echo "export GOPROXY=https://goproxy.cn" >> /root/.profile \
  && echo "source /root/.cmcc-venv/bin/activate" >> /root/.profile \
  && echo "export LC_ALL=C.UTF-8" >> /root/.profile

RUN apt install python3-venv -y

ADD requirements.txt /root/requirements.txt

SHELL ["/bin/bash", "-c"] 

RUN python3 -m venv /root/.cmcc-venv \
  && source /root/.cmcc-venv/bin/activate \
  && pip install -U -r /root/requirements.txt

RUN apt install -y rsync netcat

RUN  export https_proxy=${HTTP_PROXY} \
  && wget https://github.com/containerd/containerd/releases/download/v2.0.5/containerd-2.0.5-linux-$(dpkg --print-architecture).tar.gz \
  && tar Cxzvf /usr/local containerd-2.0.5-linux-$(dpkg --print-architecture).tar.gz \
  && rm containerd-2.0.5-linux-$(dpkg --print-architecture).tar.gz \
  && wget https://github.com/opencontainers/runc/releases/download/v1.3.0/runc.$(dpkg --print-architecture) \
  && install -m 755 runc.$(dpkg --print-architecture) /usr/local/sbin/runc \
  && rm runc.$(dpkg --print-architecture) \
  && wget https://github.com/containerd/nerdctl/releases/download/v2.0.4/nerdctl-2.0.4-linux-$(dpkg --print-architecture).tar.gz \
  && tar Cxzvf /usr/local/bin nerdctl-2.0.4-linux-$(dpkg --print-architecture).tar.gz \
  && rm nerdctl-2.0.4-linux-$(dpkg --print-architecture).tar.gz

RUN /usr/local/go/bin/go install github.com/mikefarah/yq/v4@latest \
  && apt install -y unzip

ENTRYPOINT [ "containerd" ]