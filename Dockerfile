FROM python:3.7-alpine3.13

# Setting User and Home
USER root
WORKDIR /root
ENV HOSTNAME kube-debug

# Installing debug tools


COPY bin/ctop /usr/local/bin/ctop
COPY bin/etcdctl /usr/local/bin/etcdctl
COPY bin/termshark /usr/local/bin/termshark
COPY bin/kube-debug-ttyd /usr/local/bin/kube-debug-ttyd

RUN set -ex \
    && echo "http://nl.alpinelinux.org/alpine/edge/main" > /etc/apk/repositories \
    && echo "http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories \
    && echo "http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
    && apk update \
    && apk upgrade \
    && apk add --no-cache \
    apache2-utils \
    bash \
    bash-completion \
    bind-tools \
    bird \
    bridge-utils \
    busybox-extras \
    conntrack-tools \
    curl \
    dhcping \
    drill \
    ebtables \
    ethtool \
    file\
    fping \
    graphviz \
    httpie \
    iftop \
    iperf \
    iproute2 \
    ipset \
    iptables \
    iptraf-ng \
    iputils \
    ipvsadm \
    jq \
    libc6-compat \
    liboping \
    mtr \
    net-snmp-tools \
    netcat-openbsd \
    nftables \
    ngrep \
    nmap \
    nmap-nping \
    openssl \
    scapy \
    socat \
    sox \
    strace \
    tcpdump \
    tcptraceroute \
    tshark \
    util-linux \
    vim \
    git \
    websocat \
    wget

EXPOSE 3080/tcp 22/tcp

# bash Themes
COPY config/bashrc /root/.bashrc
COPY config/motd /etc/motd

# Running
CMD ["kube-debug-ttyd","-p","3080","bash"]

