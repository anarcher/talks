# Infrakit

## What is Infrakit 

> A toolkit for building __declarative__ , __self-healing__ infrastructure

- [https://github.com/docker/infrakit](https://github.com/docker/infrakit)


## Declarative

- ansible
  - ``` apt: name=docker state=absent ```
  - ``` lineinfile: dest=/etc/sudoers regexp='^%wheel' state=present```
 
----

- JSON configuration for desired infrastructure state
   - Spec of instances : vm image, instance type, etc
   - Group properties : size, logical identifiers, etc

- Config is input to all operations

## Self-healing

- active components/processes
   - monitor infrastructure state
   - detects state divergence
   - takes action

- like Cron
- No downtime : __rolling update__

## Toolkit

- A collection of components : __Plugins__
- __Abstractions__ & API
- For managing __Groups__ of resources
    - __Scaling groups__ for any environment

### Plugin Types

- __instance__ : a resource (e.g. EC2 Instance)
- __group__ : a collection of 'like' resources (e.g. ASG)
- __flavor__ : modification to an instance (e.g. swarm join)
- manager : leadership detection , state storage for groups

### Plugin responsibilities

- __group__ : ```commit```,```free``` ,```destroy``` 
- __instance__ : ```provision```,```destory```
- __flavor__ : ```prepare```,```health```,```drain```

### Abstraction & Configuration

A common pattern for a JSON object

```
{ 
   "SomeForKey" : "SomeForValue",
   "Properties" : {..}
}
```

#### group.Spec

```
{
	"ID" : "swarm-dind":
	"Properties" : {
		"Allocation" : {
			"worker-1",
			"worker-2",
			"worker-3"
		},
		"Instance": {
      		"Plugin": "instance-dind",
      		"Properties": {
		 		"Name": "worker",
		 		"HostName" : "worker"
      		}
    	},
    	"Flavor": {
      		"Plugin": "flavor-vanilla",
      		"Properties": {
          	"Init": [
              	"docker swarm join --token SWMTKN-1-3cfnsnumc4ptz1ame7rac2dgq4atklr9nza6amux438jkd02g9-csb58zcltq94m5uf7q2im4dvi 192.168.65.2:2377"

          	]
      		}
    	}	
	}
}
```
----

```
/ Spec is the specification for a Group.  The full schema for a Group is defined by the plugin.
// In general, a Spec of an entity is set as the raw JSON value of another object's Properties.
type Spec struct {
	// ID is the unique identifier for the group.
	ID ID

	// Properties is the configuration for the group.
	// The schema for the raw JSON can be found as the *.Spec of the plugin used.
	// For instance, if the default group plugin is used, the value here will be
	// a JSON representation of github.com/docker/infrakit/plugin/group/types.Spec
	Properties *json.RawMessage
}
```

```
// Spec is the configuration schema for the plugin, provided in group.Spec.Properties
type Spec struct {
	Instance   InstancePlugin
	Flavor     FlavorPlugin
	Allocation AllocationMethod
}

// AllocationMethod defines the type of allocation and supervision needed by a flavor's Group.
type AllocationMethod struct {
	Size       uint
	LogicalIDs []instance.LogicalID
}

// InstancePlugin is the structure that describes an instance plugin.
type InstancePlugin struct {
	Plugin     string
	Properties *json.RawMessage // this will be the Spec of the plugin
}

// FlavorPlugin describes the flavor configuration
type FlavorPlugin struct {
	Plugin     string
	Properties *json.RawMessage // this will be the Spec of the plugin
}
```

## Hello infrakit

- [https://github.com/anarcher/infrakit-dind](https://github.com/anarcher/infrakit-dind)

```
{
	"ID" : "swarm-dind":
	"Properties" : {
		"Allocation" : {
			"worker-1",
			"worker-2",
			"worker-3"
		},
		"Instance": {
      		"Plugin": "instance-dind",
      		"Properties": {
		 		"Name": "worker",
		 		"HostName" : "worker"
      		}
    	},
    	"Flavor": {
      		"Plugin": "flavor-vanilla",
      		"Properties": {
          	"Init": [
              	"docker swarm join --token SWMTKN-1-3cfnsnumc4ptz1ame7rac2dgq4atklr9nza6amux438jkd02g9-csb58zcltq94m5uf7q2im4dvi 192.168.65.2:2377"

          	]
      		}
    	}	
	}
}
```

```
$./infrakit-group-default --log=5
DEBU[0000] Opening: /Users/anarch/.infrakit/plugins     
DEBU[0000] Discovered plugin at /Users/anarch/.infrakit/plugins/flavor-vanilla 
DEBU[0000] Discovered plugin at /Users/anarch/.infrakit/plugins/instance-gcp 
INFO[0000] Listening at: /Users/anarch/.infrakit/plugins/group

$./infrakit-flavor-vanilla 
INFO[0000] Listening at: /Users/anarch/.infrakit/plugins/flavor-vanilla 

$./infrakit-dind 
INFO[0000] Listening at: /Users/anarch/.infrakit/plugins/instance-dind 

$./infrakit plugin ls
NAME                	LISTEN
flavor-vanilla      	/Users/anarch/.infrakit/plugins/flavor-vanilla
group               	/Users/anarch/.infrakit/plugins/group
instance-dind       	/Users/anarch/.infrakit/plugins/instance-dind

#./infrakit group -h 
Access group plugin

Usage:
  ./infrakit group [command]

Available Commands:
  commit      commit a group configuration
  describe    describe the live instances that make up a group
  destroy     destroy a group
  free        free a group from active monitoring, nondestructive
  inspect     return the raw configuration associated with a group
  ls          list groups
  

$./infrakit group commit ./hello-dind.json
```

## RPC Processing flow

- Group : ```CommitGroup(grp Spec, pretend bool) (string, error)```
- Instance: ```DescribeInstances(tags map[string]string) ([]Description, error)```
- Flavor: ```Prepare(flavorProperties json.RawMessage, spec instance.Spec, allocation types.AllocationMethod) (instance.Spec, error)```
- Instance: ```Provision(spec Spec) (*ID, error)```

---
- Instance: ```DescribeInstances(tags map[string]string) ([]Description, error)```
- Flavor: ```Healthy(flavorProperties json.RawMessage, inst instance.Description) (Health, error)```
- Instance: ```DescribeInstances(tags map[string]string) ([]Description, error)```
- Flavor: ```Healthy(flavorProperties json.RawMessage, inst instance.Description) (Health, error)```
- ...

---

## infrakit plugin list


| plugin                                               | type     | description                             |
|:-----------------------------------------------------|:---------|:----------------------------------------|
| [swarm](pkg/example/flavor/swarm)                    | flavor   | runs Docker in Swarm mode               |
| [vanilla](pkg/example/flavor/vanilla)                | flavor   | manual specification of instance fields |
| [zookeeper](pkg/example/flavor/zookeeper)            | flavor   | run an Apache ZooKeeper ensemble        |
| [infrakit/file](pkg/example/instance/file)           | instance | useful for development and testing      |
| [infrakit/terraform](pkg/example/instance/terraform) | instance | creates instances using Terraform       |
| [infrakit/vagrant](pkg/example/instance/vagrant)     | instance | creates Vagrant VMs                     |
| [infrakit/group](cmd/group)                                   | group    | supports Instance and Flavor plugins, rolling updates |
| [docker/infrakit.aws](https://github.com/docker/infrakit.aws) | instance | creates Amazon EC2 instances                          |



- https://github.com/anarcher/infrakit.gcp

## Design Goals

![](https://github.com/docker/infrakit/raw/master/docs/images/arch.png)



