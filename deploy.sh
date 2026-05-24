#!/bin/bash
set -e

APP_HOME="/opt/mywebapp"

sudo apt update
sudo apt install -y wget mariadb-server mariadb-client nginx apt-transport-https ca-certificates curl software-properties-common

sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(# shellcheck disable=SC1091
          . /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/docker.asc
EOF

sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io
sudo usermod -aG docker "$USER"

sudo systemctl enable --now mariadb
sleep 5

if [ -f ".env" ]; then
    # shellcheck disable=SC1091
    source ".env"
fi

sudo mariadb -e "CREATE DATABASE IF NOT EXISTS \"${MARIADB_DATABASE}\";"
sudo mariadb -e "CREATE USER IF NOT EXISTS '${MARIADB_USER}'@'localhost' IDENTIFIED BY '${MARIADB_PASSWORD}';"
sudo mariadb -e "GRANT ALL PRIVILEGES ON \"${MARIADB_DATABASE}\".* TO '${MARIADB_USER}'@'localhost';"
sudo mariadb -e "FLUSH PRIVILEGES;"

sudo cp ./nginx/nginx.conf /etc/nginx/
sudo cp -r ./nginx/sites-available /etc/nginx/
sudo ln -sf /etc/nginx/sites-available/mywebapp.conf /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo systemctl enable --now nginx

sudo mkdir -p "$APP_HOME"
sudo cp .env "$APP_HOME/"
sudo chown -R root:root "$APP_HOME"
sudo chmod 600 "$APP_HOME/.env"

sudo useradd --system --shell /usr/sbin/nologin app || true

sudo useradd student -G sudo -m -s /bin/bash || true
echo "student:12345678" | sudo chpasswd
sudo chage -d 0 student

sudo useradd teacher -G sudo -m -s /bin/bash || true
echo "teacher:12345678" | sudo chpasswd
sudo chage -d 0 teacher

echo "Creating gradebook for student"
echo "24" | sudo tee /home/student/gradebook
sudo chown student:student /home/student/gradebook
sudo chmod 644 /home/student/gradebook

sudo useradd -m -U -s /bin/bash operator || true
echo "operator:12345678" | sudo chpasswd
sudo chage -d 0 operator 

echo "operator ALL=(ALL) NOPASSWD: /usr/bin/systemctl start mywebapp.service, /usr/bin/systemctl stop mywebapp.service, /usr/bin/systemctl restart mywebapp.service, /usr/bin/systemctl status mywebapp.service, /usr/bin/systemctl reload nginx" | sudo tee /etc/sudoers.d/operator
sudo chmod 440 /etc/sudoers.d/operator

if [ -n "$SUDO_USER" ]; then
  sudo usermod -L "$SUDO_USER"
fi
