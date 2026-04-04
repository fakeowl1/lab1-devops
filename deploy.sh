#!/bin/bash

set -e

APP_HOME="/opt/mywebapp"

echo "Installing necessary packages"
sudo apt update
sudo apt install wget mariadb-server mariadb-client nginx -y

sudo systemctl enable --now mariadb
sleep 5

echo "Configuring MariaDB"

[ -f .env ] && source .env

sudo mariadb -e "CREATE DATABASE ${MARIADB_DATABASE};"
sudo mariadb -e "CREATE USER '${MARIADB_USER}'@'localhost' IDENTIFIED BY '${MARIADB_PASSWORD}';"
sudo mariadb -e "GRANT ALL PRIVILEGES ON ${MARIADB_DATABASE}.* TO '${MARIADB_USER}'@'localhost';"
sudo mariadb -e "FLUSH PRIVILEGES;"

echo "Configuring nginx"
sudo cp ./nginx/nginx.conf /etc/nginx/
sudo cp -r ./nginx/sites-available /etc/nginx/
sudo ln -s /etc/nginx/sites-available/mywebapp.conf /etc/nginx/sites-enabled/
sudo systemctl enable --now nginx

echo "Installing golang"
wget https://go.dev/dl/go1.26.1.linux-amd64.tar.gz

sudo rm -rf /usr/local/go 

sudo tar -C /usr/local -xzf go1.26.1.linux-amd64.tar.gz 

rm go1.26.1.linux-amd64.tar.gz

echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/golang.sh
source /etc/profile.d/golang.sh

go mod download

echo "Building webapp"
CGO_ENABLED=0 GOOS=linux go build -o mywebapp ./cmd/mywebapp/main.go

echo "Coping project files"
sudo mkdir -p $APP_HOME
sudo cp .env $APP_HOME
sudo cp ./mywebapp $APP_HOME
sudo cp -r ./templates $APP_HOME

sudo useradd --system --shell /usr/sbin/nologin app

sudo chown -R app:app $APP_HOME
sudo chmod -R 755 $APP_HOME
sudo chmod +x "$APP_HOME/mywebapp"
sudo chmod 600 "$APP_HOME/.env"

sudo useradd student -G sudo -m -s /bin/bash
echo "student:12345678" | sudo chpasswd
sudo chage -d 0 student

sudo useradd teacher -G sudo -m -s /bin/bash
echo "teacher:12345678" | sudo chpasswd
sudo chage -d 0 teacher

echo "Creating gradebook for student"
echo "24" | sudo tee /home/student/gradebook
sudo chown student:student /home/student/gradebook
sudo chmod 644 /home/student/gradebook

sudo useradd -g operator -m -s /bin/bash operator
echo "operator:12345678" | sudo chpasswd
sudo chage -d 0 operator 

echo "operator ALL=(ALL) NOPASSWD: /bin/systemctl start mywebapp.service, /bin/systemctl stop mywebapp.service, /bin/systemctl restart mywebapp.service, /bin/systemctl status mywebapp.service, /bin/systemctl reload nginx" | sudo tee /etc/sudoers.d/operator
sudo chmod 440 /etc/sudoers.d/operator

if [ -n "$SUDO_USER" ]; then
  sudo usermod -L "$SUDO_USER"
fi

sudo cp systemd/mywebapp.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now mywebapp.service
sudo systemctl restart nginx

echo "Deployment environment ready!"
