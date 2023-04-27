#! /bin/bash
apt-get update
apt-get install -y wget
wget https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
EOF