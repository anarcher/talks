docker service create --name nginx-1 --mount type=bind,source=/etc/hosts,target=/usr/share/nginx/html/hosts nginx
docker service update --publish-add 80:80 nginx

docker network create -d overlay my-net

docker service create --name shell --network my-net alpine sleep 100000
docker service update --constraint-add="node.role==manager" shell 

docker service create --name shell-2 --publish 8001:8001 --constraint="node.role==manager" alpine sleep 100000

docker service create --name nginx-2 --network my-net --endpoint-mode=dnsrr --mount type=bind,source=/etc/hosts,target=/usr/share/nginx/html/hosts nginx

docker service create --name nginx-3 --publish 8002:80 --network my-net --endpoint-mode=vip --mount type=bind,source=/etc/hosts,target=/usr/share/nginx/html/hosts nginx

docker service create --mode global --constraint "node.role==worker" --endpoint-mode=dnsrr --name cadvisor --mount type=bind,source=/,target=/rootfs:ro --mount type=bind,source=/var/run,target=/var/run --mount type=bind,source=/sys,target=/sys --mount type=bind,source=/var/lib/docker,target=/var/lib/docker google/cadvisor



