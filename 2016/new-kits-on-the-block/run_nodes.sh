SWARM_TOKEN=$(docker swarm join-token -q worker)
SWARM_MASTER=$(docker info | grep -w 'Node Address' | awk '{print $3}')
NUM_WORKERS=3 
for i in $(seq "${NUM_WORKERS}"); do 
	docker run -d --privileged --name worker-${i} \
            --hostname=worker-${i} -p ${i}2375:2375 docker:1.12.3-dind
    docker --host localhost:${i}2375 swarm join --token ${SWARM_TOKEN} ${SWARM_MASTER}:2377
done
