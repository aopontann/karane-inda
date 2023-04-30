#! /bin/bash
apt-get update
apt-get install -y wget
wget -q https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go install github.com/aopontann/karane-inda@v0.1.1
echo -e 'OPENAI_API_KEY=<API_KEY> \nDISCORD_TOKEN=<TOKEN>' >> .env
go/bin/karane-inda