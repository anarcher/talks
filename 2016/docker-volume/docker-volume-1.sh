$docker volume create --name=scribe_conf
scribe_conf

$docker volume ls 
DRIVER              VOLUME NAME
local               scribe_conf

$docker volume inspect scribe_conf
[
    {
        "Name": "scribe_conf",
        "Driver": "local",
        "Mountpoint": "/var/lib/docker/volumes/scribe_conf/_data"
    }
]

$docker volume rm scribe_conf
scribe_conf

$docker volume ls
DRIVER              VOLUME NAME

