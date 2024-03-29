## linodectl create lkecluster

Create an LKE Cluster

```
linodectl create lkecluster NAME [args...] [flags]
```

### Examples

```
  # Create a 1.21 Cluster with a 5 node pool of type g6-standard-2
  linodectl create lkecluster -v1.21 --region us-east --pool g6-standard-2:5

  # Create a Cluster on the latest support version of Kubernetes
  linodectl create lkecluster --region us-central --pool g6-standard-1:3

  # Create a Cluster with 2 node pools
  linodectl create lkecluster --region eu-west -v1.20 --pool g6-standard-1:3 --pool g6-standard-3:2

  # Create a Cluster with a highly available control plane
  linodectl create lkecluster --region us-east -v1.20 --ha --pool g6-standard-1:3
```

### Options

```
      --ha               If true, this cluster will be deployed with a highly available control plane
  -h, --help             help for lkecluster
      --pool strings     Node pool configuration in format <instance-type>:<count> (i.e. g6-standard-2:3)
  -p, --profile string   The profile to use for communicating with the Linode API
      --region string    Region to deploy in
      --tags strings     Tags to add to this cluster
  -v, --version string   Version of Kubernetes to deploy
```

### SEE ALSO

* [linodectl create](linodectl_create.md)	 - Create a Linode resource

###### Auto generated by spf13/cobra on 2-Dec-2021
