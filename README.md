

- anarcher/talks:2015-gokit

```
docker run --net=host anarcher/talks:2015-gokit  

#boot2docker (mac)
docker run --net=host anarcher/talks:2015-gokit -http=:3999 -orighost=$(docker-machine ip dev)
```
