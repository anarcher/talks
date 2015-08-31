docker-machine create -d virtualbox --virtualbox-boot2docker-url=https://github.com/anarcher/boot2docker-experimental/releases/download/1.9/boot2docker-1.9.iso infra

docker $(docker-machine config infra) run -d \
    -p "8500:8500" \
    -h "consul" \
    progrium/consul -server -bootstrap

docker-machine create -d virtualbox --virtualbox-boot2docker-url=https://github.com/anarcher/boot2docker-experimental/releases/download/1.9/boot2docker-1.9.iso --engine-opt="kv-store=consul:$(docker-machine ip infra):8500" --engine-label="com.docker.network.driver.overlay.bind_interface=eth0" demo0

docker-machine create -d virtualbox --virtualbox-boot2docker-url=https://github.com/anarcher/boot2docker-experimental/releases/download/1.9/boot2docker-1.9.iso --engine-opt="kv-store=consul:$(docker-machine ip infra):8500" --engine-label="com.docker.network.driver.overlay.bind_interface=eth0" demo1

docker $(docker-machine config demo0) run --rm -it --name first busybox

docker $(docker-machine config demo1) run --rm -it --name second busybox

ping second

ping first

