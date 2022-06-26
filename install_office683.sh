#! /bin/bash
echo "Installing Dependencies"
sudo apt update
sudo apt install nano

echo "Fetching Assets"
rm -rf /opt/saenuma/office683
mkdir -p /opt/saenuma/office683
wget -q https://storage.googleapis.com/pandolee/office683/3/office683.tar.xz -O /opt/saenuma/office683.tar.xz
tar -xf /opt/saenuma/office683.tar.xz -C /opt/saenuma/office683

sudo chmod +x /opt/saenuma/office683/bin/o6sites
sudo chmod +x /opt/saenuma/office683/bin/o6ssl
sudo chmod +x /opt/saenuma/office683/bin/ssl-proxy-linux-amd64

sudo cp /opt/saenuma/office683/bin/o6sites.service /etc/systemd/system/o6sites.service
sudo cp /opt/saenuma/office683/bin/o6ssl.service /etc/systemd/system/o6ssl.service

sudo cp /opt/saenuma/office683/bin/o6sites /usr/local/bin/
sudo cp /opt/saenuma/office683/bin/o6ssl /usr/local/bin/
sudo cp /opt/saenuma/office683/bin/ssl-proxy-linux-amd64 /usr/local/bin/

echo "Starting Services"
sudo systemctl daemon-reload
sudo systemctl start o6ssl
sudo systemctl start o6sites
