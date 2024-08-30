#!/bin/bash

cp /vagrant/resources/docker/* /usr/bin/

sudo groupadd docker
sudo usermod -aG docker vagrant

mkdir -p /usr/local/lib/docker/cli-plugins
cp /vagrant/resources/docker-compose/docker-compose /usr/local/lib/docker/cli-plugins
cp /vagrant/resources/docker-compose/docker-compose /usr/bin

cp /vagrant/resources/docker.socket /etc/systemd/system/
cp /vagrant/resources/docker.service /etc/systemd/system/

systemctl start docker.service
systemctl enable docker.service
