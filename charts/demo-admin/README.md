# Helm Chart for NUTS demo-admin
This chart allows the ease of running NUTS demo-admin on a Kubernetes cluster. 
All the NUTS demo-admin information is persisted on [Persisted Volumes](https://kubernetes.
io/docs/concepts/storage/persistent-volumes/).

## Configure your NUTS node
All the configurable properties can be found at [./values.yaml](./values.yaml).

The configuration contains default Helm properties. In addition to these values,
there are `demo-admin` config properties. This contains 3 sections:

| Section     | Description                                                                                                                                                                                                                                            |
|-------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `admin.config` | Represents the `/app/server.config.yaml` file. All configurable properties can be found in the <br/>main [README](../../README.md#Configuration). The properties are loaded into a `ConfigMap` and mounted as `/app/server.config.yaml` inside the Pod(s). |
| `admin.data`   | Contains configurable properties for the `PersistedVolume` that will be created. This will be used to write all NUTS demo-admin data to.                                                                                                                          |

### Special properties
NUTS demo-adminn allows binding to specific interfaces on the host machines. In the case of Kubernetes, this is already 
taken care 
of. However, we do need to expose the `http` and `gRPC` ports. This is extracted from the following properties:

| Property                                                     | Value (default) |
|--------------------------------------------------------------|-----------------|
| `admin.config.port` (must align with `service.internalPort`) | :1303 |


### Overriding values
#### From Source
The properties can be manually changed in the [./values.yaml](./values.yaml), or they can be overwritten whilst running
`helm install` via the `--set x=y` parameter.

#### From the NUTS Helm Repo
 
The default values can be viewed with the following command: 
```shell
helm show values nuts-repo/nuts-admin-demo-chart
```

You can then override any of these settings in a YAML formatted file, and then pass that file during [installation](#from-the-nuts-helm-repo-1).

## Installing NUTS admin-demo
### From Source

Execute the following command from the root of the chart folder. Replace `<NAME>` with the name you 
wish to give this Helm installation.
```
helm install <NAME> .
```
### From the NUTS Helm Repo

Add the NUTS helm Repo with the following command:
```shell 
helm repo add nuts-repo https://nuts-foundation.github.io/nuts-registry-admin-demo/
```
This should list available releases with the following command:
```shell
helm search repo nuts-registry-admin-demo
```

After this, the desired version can be installed with the following command:
```shell
helm repo update              # Make sure we get the latest list of charts
helm install -f values.yaml <NAME> nuts-repo/nuts-registry-admin-demo-chart
```

Note that the `values.yaml` in the above command is the result from the [configuration step](#from-the-nuts-helm-repo).

## Uninstalling NUTS
As the `PersistedVolume` can contain crucial data (like the private keys), by default, the uninstall command will not remove it and its 
`PersistedVolumeClaim`. If you're sure it can be deleted, this can be done with the following command:
```shell
kubectl delete pvc nuts- nuts-registry-admin-demo-data-pvc
kubectl delete pv nuts- nuts-registry-admin-demo-data-pv
```
