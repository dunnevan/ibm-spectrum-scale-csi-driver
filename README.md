
   * [Welcome to the public Beta of IBM Spectrum Scale Container Storage Interface (CSI) Driver](#welcome-to-the-public-beta-of-ibm-spectrum-scale-container-storage-interface-csi-driver)
      * [IBM Spectrum Scale Introduction](#ibm-spectrum-scale-introduction)
      * [IBM Spectrum Scale Container Storage Interface (CSI) driver](#ibm-spectrum-scale-container-storage-interface-csi-driver)
         * [Supported Features of the CSI driver](#supported-features-of-the-csi-driver)
         * [Limitations of the CSI driver](#limitations-of-the-csi-driver)
         * [Pre-requisites for installing and running the CSI driver](#pre-requisites-for-installing-and-running-the-csi-driver)
      * [Building the docker image](#building-the-docker-image)
      * [Install and Deploy the Spectrum Scale CSI driver](#install-and-deploy-the-spectrum-scale-csi-driver)
      * [Static Provisioning](#static-provisioning)
      * [Dynamic Provisioning](#dynamic-provisioning)
         * [Storageclass](#storageclass)
      * [Advanced Configuration](#advanced-configuration)
         * [Remote mount support](#remote-mount-support)
         * [Node Selector](#node-selector)
         * [Kubernetes node to Spectrum Scale node mapping](#kubernetes-node-to-spectrum-scale-node-mapping)
      * [Cleanup](#cleanup)
      * [Troubleshooting](#troubleshooting)
      * [Environments in Test](TESTCONFIG.md#environments-in-test)
      * [Example Hardware Configs](TESTCONFIG.md#example-hardware-configs)
      * [Example of using the Install Toolkit to build a Spectrum Scale cluster for testing the CSI driver](TESTCONFIG.md#example-of-using-the-install-toolkit-to-build-a-spectrum-scale-cluster-for-testing-the-csi-driver)
      * [Links](#links)

  

# Welcome to the public Beta of IBM Spectrum Scale Container Storage Interface (CSI) Driver

DISCLAIMER: This Beta driver is provided as is, without warranty. Any issue will be handled on a best-effort basis. See the Spectrum Scale Users Group links at the very bottom for a community to share and discuss test efforts.

  
## IBM Spectrum Scale Introduction

IBM Spectrum Scale is a clustered file system providing concurrent access to a single file system or set of file systems from multiple nodes. The nodes can be SAN attached, network attached, a mixture of SAN attached and network attached, or in a shared nothing cluster configuration. This enables high performance access to this common set of data to support a scale-out solution or to provide a high availability platform.

IBM Spectrum Scale has many features beyond common data access including data replication, policy based storage management, and multi-site operations. You can create a cluster of AIX® nodes, Linux nodes, Windows server nodes, or a mix of all three. IBM Spectrum Scale can run on virtualized instances providing common data access in environments, leverage logical partitioning, or other hypervisors. Multiple IBM Spectrum Scale clusters can share data within a location or across wide area network (WAN) connections. For more information on IBM Spectrum Scale features, see the Product overview section in the IBM Spectrum Scale: Concepts, Planning, and Installation Guide.

Please refer to the [IBM Spectrum Scale Knowledge Center](https://www.ibm.com/support/knowledgecenter/en/STXKQY/ibmspectrumscale_welcome.html) for more information.
  

## IBM Spectrum Scale Container Storage Interface (CSI) driver

The IBM Spectrum Scale Container Storage Interface (CSI) driver allows IBM Spectrum Scale to be used as persistent storage for stateful application running in Kubernetes clusters. Through this CSI Driver, Kubernetes persistent volumes (PVs) can be provisioned from IBM Spectrum Scale. Thus, containers can be used with stateful microservices, such as database applications (MongoDB, PostgreSQL etc), web servers (nginx, apache), or any number of other containerized applications needing provisioned storage.

### Supported Features of the CSI driver

IBM Spectrum Scale Container Storage Interface (CSI) driver supports the following features:

- **Static provisioning:** Ability to use existing directories as persistent volumes
- **Lightweight dynamic provisioning:** Ability to create directory-based volumes dynamically
- **Fileset-based dynamic provisioning:** Ability to create fileset-based volumes dynamically
- **Multiple file systems support:** Volumes can be created across multiple file systems
- **Remote mount support:** Volumes can be created on a remotely mounted file system
  
### Limitations of the CSI driver

The IBM Spectrum Scale Container Storage Interface (CSI) driver has the following limitations:

- The size specified in PersistentVolumeClaim for lightweight volume and dependent fileset volume, is not honored.
- Volumes cannot be mounted in read-only mode.
- Maximum number of supported volumes that can be created using independent fileset storage class is 998 (excluding the root fileset and primary fileset reserved for CSI driver). This is based upon the [fileset maximums for IBM Spectrum Scale](https://www.ibm.com/support/knowledgecenter/STXKQY/gpfsclustersfaq.html#filesets)
- The IBM Spectrum Scale GUI server is relied upon for performing file system and cluster operations. If the GUI password or CA certificate expires, manual intervention is needed by the admin to reset the GUI password or generate a new certificate and update the configuration of the CSI driver. In this case, a restart of the CSI driver will be necessary.
- Rest API status, used by the CSI driver, may lag from actual state, causing PVC mount or unmount failures.
- Although multiple instances of the Spectrum Scale GUI are allowed, the CSI driver is currently limited to point to a single GUI node.
- External attacher and external provisioner run as statefulsets, which by design do not failover to different node in case docker/kubelet is brought down. (They however failover to another node when the node itself is deleted explicitly). It is recommended to run the attacher and provisioner on two separate infrastructure nodes which can be done by using [Node Selector](#node-selector) or by configuring it during [operator deployment](https://github.com/IBM/ibm-spectrum-scale-csi-operator)

### Pre-requisites for installing and running the CSI driver

- IBM Spectrum Scale / Kubernetes overlap should be as follows

  | Node Type | Spectrum Scale | Kubernetes |
  |--|--|--|
  | **Master node(s)** | do not install | required |
  | **Worker node(s)** | required | required |
  | **GUI node** | required | do not install |
  | **NSD node** | required | optional |

- Red Hat 7.6 (**kernel 3.10.0-957 or higher**) on Spectrum Scale nodes

- IBM Spectrum Scale version 5.0.4.1 is installed.

- An IBM Spectrum Scale GUI is up and running on a Spectrum Scale node and a user is created and part of the `CsiAdmin` group

  ```
  /usr/lpp/mmfs/gui/cli/mkuser <__username__> -p <__password__> -g CsiAdmin
  ```

- Kubernetes ver 1.13+ cluster is created
 
- All Kubernetes worker nodes must also be Spectrum Scale client nodes. Install the Spectrum Scale client on all Kubernetes worker nodes and ensure they are added to the Spectrum Scale cluster. (To install Spectrum Scale and CSI driver only on selected nodes, perform the steps from [Node Selector](#node-selector)

- The Filesystem to be used for persistent storage must be mounted on the Spectrum Scale GUI node as well as all Kubernetes worker nodes. (*If multiple filesystems are to be used as persistent storage for containers, then all need to be mounted*)

- Quota must be enabled on the filesystem (*required for fileset based dynamic provisioning*)
  ```
  mmchfs <__filesystem_name__> -Q yes
  ```

## Building the docker image

**Using multi-stage build**

Pre-requisite: Docker 17.05 or higher is installed on local build machine.


1. Clone the code

   ```
   git clone https://github.com/IBM/ibm-spectrum-scale-csi-driver.git
   cd ibm-spectrum-scale-csi-driver
   ```

2. Invoke multi-stage build

   ```
   docker build -t ibm-spectrum-scale-csi:v0.9.2 -f Dockerfile.msb .
   ```

   *On podman setup, use this command instead:*

   ```
   podman build -t ibm-spectrum-scale-csi:v0.9.2 -f Dockerfile.msb .
   ```

3. save the docker image

   ```
   docker save ibm-spectrum-scale-csi:v0.9.2 -o ibm-spectrum-scale-csi_v0.9.2.tar
   ```

   *On podman setup, use this command instead:*

   ```
   podman save ibm-spectrum-scale-csi:v0.9.2 -o ibm-spectrum-scale-csi_v0.9.2.tar
   ```

   A tar file of docker image will be stored under the _output directory.

## Install and Deploy the Spectrum Scale CSI driver

1. Copy and load the docker image on all Kubernetes worker nodes

   ```
   docker image load -i ibm-spectrum-scale-csi_v0.9.2.tar
   ```

   *On OpenShift setup, use this command instead:*

   ```
   podman image load -i ibm-spectrum-scale-csi_v0.9.2.tar
   ```

2. Deploy CSI driver

   **Method 1: Operator (Recommended)**

   Follow the instructions from [ibm-spectrum-scale-csi-operator](https://github.com/IBM/ibm-spectrum-scale-csi-operator) for deployment of CSI driver

   **Method 2: Install script**

   a. Update `deploy/spectrum-scale-driver.conf` with your cluster and environment details.

   Note that on OpenShift setup, the image is listed as `localhost/ibm-spectrum-scale-csi:v0.9.2`. Change the value of "spectrumscaleplugin" parameter in images section accordingly. 

   b. Set the environment variable CSI_SCALE_PATH to ibm-spectrum-scale-csi-driver directory

   ```
   export CSI_SCALE_PATH=$GOPATH/src/github.com/IBM/ibm-spectrum-scale-csi-driver
   ```

   c. Run the install helper script:

   ```
   tools/spectrum-scale-driver.py $CSI_SCALE_PATH/deploy/spectrum-scale-driver.conf
   ```

   Review the generated configuration files in deploy.

   d. Run the `deploy/create.sh` script to deploy the plugin

3. Check that the csi pods are up and running

   ```
   % kubectl get pod
   NAME READY STATUS RESTARTS AGE
   ibm-spectrum-scale-csi-7d8jg 2/2 Running 0 7s
   ibm-spectrum-scale-csi-attacher-0 1/1 Running 0 8s
   ibm-spectrum-scale-csi-provisioner-0 1/1 Running 0 8s
   ```


## Static Provisioning

In static provisioning, the backend storage volumes and PVs are created by the administrator. Static provisioning can be used to provision a directory or fileset with existing data.

For static provisioning of existing directories perform the following steps:

- Generate static pv yaml file using helper script

   ```
   tools/generate_pv_yaml.sh --filesystem gpfs0 --size 10 \
   --linkpath /ibm/gpfs0/pvfileset/static-pv --pvname static-pv
   ```

- Use sample static_pvc and pod files for sanity test under `examples/static`

   ```
   kubectl apply -f examples/static/static_pv.yaml
   kubectl apply -f examples/static/static_pvc.yaml
   kubectl apply -f examples/static/static_pod.yaml
   ```
  

## Dynamic Provisioning

Dynamic provisioning is used to dynamically provision the storage backend volume based on the storageclass.

### Storageclass
Storageclass defines what type of backend volume should be created by dynamic provisioning. IBM Spectrum Scale CSI driver supports creation of directory based (also known as lightweight volumes) and fileset based (independent as well as dependent) volumes. Following parameters are supported by BM Spectrum Scale CSI driver storageclass:

 - **volBackendFs**: Filesystem on which the volume should be created. This is a mandatory parameter.
 - **clusterId**: Cluster ID on which the volume should be created. 
 - **volDirBasePath**: Base directory path relative to the filesystem mount point under which directory based volumes should be created. If specified, the storageclass is used for directory based (lightweight) volume creation. If not specified, storageclass creates fileset based volumes.
 - **uid**: UID with which the volume should be created. Optional
 - **gid**: UID with which the volume should be created. Optional
 - **filesetType**: Type of fileset. Valid values are "independent" or "dependent". Default is "independent". 
 - **parentFileset**: Specifies the parent fileset under which dependent fileset should be created. Mandatory if "fileset-type" is specified.
 - **inodeLimit**: Inode limit for fileset based volumes. If not specified, default Spectrum Scale inode limit of 1million is used.
 
For dynamic provisioning, use sample storageclass, pvc and pod files for sanity test under examples/dynamic

Example:

   ```
   kubectl apply -f examples/dynamic/fileset/storageclassfileset.yaml
   kubectl apply -f examples/dynamic/fileset/pvcfset.yaml
   kubectl apply -f examples/dynamic/fileset/podfset.yaml
   ```

## Advanced Configuration

Following is advanced configuration of IBM Spectrum Scale CSI driver and is not supported through the installer "spectrum-scale-driver.py". Perform the below steps after running the installer "spectrum-scale-driver.py".

Note: This advanced configuration is supported through operator and manual steps given below are not needed.

### Remote mount support

IBM Spectrum Scale provides a feature to mount a Spectrum Scale file system that belongs to another IBM Spectrum Scale cluster. Consider the case where Kubernetes worker nodes are part of a "primary" Spectrum Scale cluster. This primary cluster has filesystems mounted from a "remote" Spectrum Scale cluster. 

In order to deploy CSI driver on such a configuration, following steps should be performed after running the installer "spectrum-scale-driver.py":

- Update `deploy/spectrum-scale-config.json` file with remote cluster and filesystem name information under "primary" section by adding the two parameters as:

   * "**remoteCluster**":"<remote cluster ID>",  
   * "**remoteFs**":"remote filesystem name" (Required only if remote filesystem name is different than the locally mounted filesystem name)

- Make another entry for this cluster under the "clusters" section of `deploy/spectrum-scale-config.json` as:

   ```
   {"id":"2954738785946888888",
    "secrets":"secret2",
     "restApi": [
        {"guiHost":"172.16.1.33"
        }
     ]
   }
   ```

- Ensure that a new entry is created in secrets list in `deploy/spectrum-scale-secret.json` for "secret2" as:

   ```
   {
      "kind": "Secret",
      "apiVersion": "v1",
      "metadata": {
         "name": "secret2"
      },
      "data": {
         "username": "YWRtaW4=",
         "password": "MWYyZDFlMmU2N2Rm"
      }
   }
   ```
   **Note:** username and passoword are base64 encoded.

- Add an entry for secret2 in deploy/csi-plugin.yaml file under "volumes":

   ```
   - name: secret2
     secret:
       secretName: secret2
   ```

   Add corresponding entry under "containers -> ibm-spectrum-scale-csi -> volumeMounts" section:

   ```
   - name: secret2
     mountPath: /var/lib/ibm/secret2
     readOnly: true
   ```

- Deploy the driver by running `deploy/create.sh`

- For lightweight dynamic provisioning, no change in storageclass is needed.

- For fileset based dynamic provisioning, use the storageclass parameters as below:

   * **volBackendFs**: Filesystem on which the volume should be created. Use the remote cluster filesystem name here.
   * **clusterId**: Remote Cluster ID on which the volume (fileset) should be created. 
   * **localFs**: Name of the locally mounted filesystem. This is required only if the local name and remote filesystem names are different.
	Rest of the storageclass parameteres remain valid.

### Node Selector

Node selector is used to control on which Kubernetes worker nodes the IBM Spectrum Scale CSI driver should be running. Node selector also helps in cases where new worker nodes are added to Kubernetes cluster but does not have IBM Spectrum Scale installed, in this case we would not want the CSI driver to be deployed on these nodes. If node selector is not used, CSI driver gets deployed on all worker nodes.

To use this feature, perform the following steps after running the installer "spectrum-scale-driver.py":

- Label the Kubernetes worker nodes where IBM Spectrum Scale is running. Example:
	```
	kubectl label node node7 spectrumscalenode=yes --overwrite=true
	```
- Uncomment the lines from the following files: 
    * `deploy/csi-plugin-attacher.yaml`
    * `deploy/csi-plugin-provisioner.yaml`
    * `deploy/csi-plugin.yaml`

    ```
    #      nodeSelector:  
    #        spectrumscalenode: "yes"
    ```

- Deploy the driver by running `deploy/create.sh`

**Note:** If you choose to run csi plugin on selective nodes using the node selector then make sure pod using scale csi pvc are getting scheduled on nodes where csi driver is running.

### Kubernetes node to Spectrum Scale node mapping

In an environment where Kubernetes node names are different than the Spectrum Scale node names, this mapping feature must be used for application pods with Spectrum Scale as persistent storage to be successfully mounted.

To use this feature, perform the following steps after running the installer "spectrum-scale-driver.py":

- Add new environment variable in `deploy/csi-plugin.yaml` under container "*- name: ibm-spectrum-scale-csi*", where name of the environment variable is Kubernetes node name and value is the Spectrum Scale node name. 

  ```
  env:
     - name: k8snodename1  
       value: "scalenodename1"
     - name: k8snodename2  
       value: "scalenodename2"
  ```

  **Note:** Only add those nodes whose name is different in Kubernetes (`kubectl get nodes`) and Spectrum Scale (`mmlscluster/mmlsnode`)
	
- Deploy the driver by running `deploy/create.sh`

## Cleanup

1. Delete the resources that were created (pod, pvc, pv, storageclass)

2. Run deploy/destroy.sh script to cleanup the plugin resources

3. Find the CSI driver docker image and remove from all Kubernetes worker nodes

   ```
   % docker images -a | grep ibm-spectrum-scale-csi
   ibm-spectrum-scale-csi v0.9.2 465ca978127a 18 minutes ago 109MB

   % docker rmi 465ca978127a
   ```

## Troubleshooting

The tool `spectrum-scale-driver-snap.sh` collects the CSI driver debug data and stores in the given output directory.

Usage:

   ```
   spectrum-scale-driver-snap.sh [-n namespace] [-o output-dir] [-h]
   -n: Debug data for CSI resources under this namespace will be collected.
       If not specified, default namespace is used. The tool returns error
       if CSI is not running under the given namespace.
   -o: Output directory where debug data will be stored. If not specified,
       the debug data is stored in current directory.
   -h: Prints the usage
   ```

## Links

[IBM Spectrum Scale Knowledge Center Welcome Page](https://www.ibm.com/support/knowledgecenter/en/STXKQY/ibmspectrumscale_welcome.html)
The Knowledge Center contains all official Spectrum Scale information and guidance.

[IBM Spectrum Scale FAQ](https://www.ibm.com/support/knowledgecenter/en/STXKQY/gpfsclustersfaq.html)
Main starting page for all Spectrum Scale compatibility information.

[IBM Spectrum Scale Protocols Quick Overview](https://www.ibm.com/developerworks/community/wikis/home?lang=en#!/wiki/fa32927c-e904-49cc-a4cc-870bcc8e307c/page/Protocols%20Quick%20Overview%20for%20IBM%20Spectrum%20Scale)
Guide showing how to quickly install a Spectrum Scale cluster. Information similar to the above Install Toolkit example.

[IBM Block CSI driver](https://github.com/IBM/ibm-block-csi-driver)
CSI driver supporting multiple IBM storage systems.

[Installing kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)
Main Kubernetes site detailing how to install kubeadm and create a cluster.

[OpenShift Container Platform 4.x Tested Integrations](https://access.redhat.com/articles/4128421)
Red Hat's test matrix for OpenShift 4.x.

[IBM Storage Enabler for Containers Welcome Page](https://www.ibm.com/support/knowledgecenter/en/SSCKLT/landing/IBM_Storage_Enabler_for_Containers_welcome_page.html)
Flex Volume driver released in late 2018 with a HELM update in early 2019, providing compatibility with IBM Spectrum Scale for file storage and multiple IBM storage systems for block storage. Future development efforts have shifted to CSI.

[Spectrum Scale Users Group](http://www.gpfsug.org/)
A group of both IBM and non-IBM users, interested in Spectrum Scale

[Spectrum Scale Users Group Mailing List and Slack Channel](https://www.spectrumscaleug.org/join/)
Join everyone and let the team know about your experience with the CSI driver
