docker-machine create -d virtualbox --virtualbox-boot2docker-url=https://github.com/anarcher/boot2docker-experimental/releases/download/1.9/boot2docker-1.9.iso infra

docker $(docker-machine config infra) run -d \
    -p "8500:8500" \
    -h "consul" \
    progrium/consul -server -bootstrap


docker-machine --debug create -d virtualbox --virtualbox-boot2docker-url=https://github.com/anarcher/boot2docker-experimental/releases/download/1.9/boot2docker-1.9.iso --engine-opt="kv-store=consul:$(docker-machine ip infra):8500" --engine-label="com.docker.network.driver.overlay.bind_interface=eth1" demo0

docker-machine --debug create -d virtualbox --virtualbox-boot2docker-url=https://github.com/anarcher/boot2docker-experimental/releases/download/1.9/boot2docker-1.9.iso --engine-opt="kv-store=consul:$(docker-machine ip infra):8500" --engine-label="com.docker.network.driver.overlay.bind_interface=eth1" --engine-label="com.docker.network.driver.overlay.neighbor_ip=$(docker-machine ip demo0)" demo1


docker network create -d overlay dev 
docker network ls
docker network info dev

docker $(docker-machine config demo0) run -i -t --rm --publish-service=demo0.dev.overlay -h demo0 --name demo0 debian
docker service publish demo0.dev
docker service attach $cid demo0.dev

docker $(docker-machine config demo1) run -i -t --rm --publish-service=demo1.dev.overlay -h demo1 --name demo1 debian
docker service publish demo1.dev
docker service attach $cid demo1.dev


ping demo0

ping demo1.dev

