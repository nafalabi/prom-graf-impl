Vagrant.configure("2") do |config|
  # config.vagrant.plugins = "vagrant-docker-compose"

  config.vm.define "master" do |master|
    # master.vm.provision :docker
    # master.vm.provision :docker_compose
    master.vm.box = "hashicorp/bionic64"
    master.vm.box_version = "1.0.282"
    master.vm.network "private_network", ip: "192.168.56.2"
    master.vm.provision :shell, inline: "sh /vagrant/retrieve-docker.sh"
    master.vm.provision :shell, inline: "sh /vagrant/install-docker.sh"
  end
  
  config.vm.define "client1" do |client1|
    # client1.vm.provision :docker
    # client1.vm.provision :docker_compose
    client1.vm.box = "hashicorp/bionic64"
    client1.vm.box_version = "1.0.282"
    client1.vm.network "private_network", ip: "192.168.56.3"
    client1.vm.provision :shell, inline: "sh /vagrant/install-docker.sh"
  end
  
  config.vm.define "client2" do |client2|
    # client2.vm.provision :docker
    # client2.vm.provision :docker_compose
    client2.vm.box = "hashicorp/bionic64"
    client2.vm.box_version = "1.0.282"
    client2.vm.network "private_network", ip: "192.168.56.4"
    client2.vm.provision :shell, inline: "sh /vagrant/install-docker.sh"
  end

end
