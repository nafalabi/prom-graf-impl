#!/bin/bash

BASE_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
IS_DOWNLOADED="$BASE_DIR/resources/.downloaded"

if [[ -e $IS_DOWNLOADED ]]; then
  echo "Docker is already downloaded"
  echo "Skipping...."
  exit
fi

echo "Downloading docker..."

mkdir /vagrant/resources/docker
curl -o /vagrant/resources/docker.tgz https://download.docker.com/linux/static/stable/x86_64/docker-27.1.2.tgz
tar -xzvf /vagrant/resources/docker.tgz -C /vagrant/resources

mkdir /vagrant/resources/docker-compose
curl -SL https://github.com/docker/compose/releases/download/v2.29.1/docker-compose-linux-x86_64 -o /vagrant/resources/docker-compose/docker-compose
chmod +x /vagrant/resources/docker-compose/docker-compose

echo "true" > $IS_DOWNLOADED

echo "Download completed!"
