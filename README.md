## Prometheus + Grafana Implementation Using Vagrant

an overview of using prometheus & grafana as a server monitoring using vagrant vm manager

### Prerequisite

- vagrant v2.4.1 or equal
- virtualbox v7.x or equal

### To Start

1. make sure vagrant & virtualbox is installed in your environtment
2. create and configure `aletmanager.yml` file by using the template located at `./master/alertmanager_conf/alertmanager.yml_sample`
3. run `vagrant up`

### To test sending alert

1. make sure you have started vagrant
2. ssh into master node by running `vagrant ssh master`
3. cd into `/vagrant/master`
4. run `sh trigger-alert.sh`
