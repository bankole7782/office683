#! /bin/bash
echo "Installing Dependencies"
sudo apt update
sudo apt install nano

echo "Fetching Assets"
rm -rf /opt/saenuma/office683
mkdir -p /opt/saenuma/office683
mkdir -p /var/lib/office683/
wget -q https://storage.googleapis.com/pandolee/office683/1/office683.tar.xz -O /opt/saenuma/office683.tar.xz
tar -xf /opt/saenuma/office683.tar.xz -C /opt/saenuma/office683

sudo chmod +x /opt/saenuma/office683/bin/o6sites
sudo chmod +x /opt/saenuma/office683/bin/o6ssl
sudo chmod +x /opt/saenuma/office683/bin/ssl-proxy-linux-amd64

sudo cp /opt/saenuma/office683/bin/o6sites.service /etc/systemd/system/o6sites.service
sudo cp /opt/saenuma/office683/bin/o6ssl.service /etc/systemd/system/o6ssl.service

sudo cp /opt/saenuma/office683/bin/o6sites /usr/local/bin/
sudo cp /opt/saenuma/office683/bin/o6ssl /usr/local/bin/
sudo cp /opt/saenuma/office683/bin/ssl-proxy-linux-amd64 /usr/local/bin/

echo "Saving Config"
cat <<EOT > /var/lib/office683/install.zconf
// name of the company that created this office tools information
company_name: Test1

// logo of the company on your website.
company_logo: https://sae.ng/static/logo.png

// admin_pass is the password used by all admins
admin_pass:

// admin_email is for contacting the admin to get access
admin_email: admin@admin.com

// domain must be set after launching your server
domain:

EOT

echo "Starting Services"
sudo systemctl daemon-reload
sudo systemctl start o6sites
