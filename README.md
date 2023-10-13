# NCP IAM Authenticator for Kubernetes

NAVER Cloud Platform Kubernetes Service provides IAM authentication through `ncp-iam-authenticator`.  
To use the kubectl command through IAM authentication, you should install `ncp-iam-authenticator` and edit the kubectl configuration file to use it for authentication.  
The initial project was developed by NAVER Cloud Platform Kubernetes Service engineers, and now anyone can contribute to the project.
## Contents
- [Guide](#guide)
- [Installation](#installation)
- [Creating IAM authentication kubeconfig](#creating-iam-authentication-kubeconfig)
    - [Set ncp-iam-authenticator API authentication key value](#set-ncp-iam-authenticator-api-authentication-key-value)
    - [Use command ncp-iam-authenticator create-kubeconfig](#use-command-ncp-iam-authenticator-create-kubeconfig)
- [IAM authentication user management](#iam-authentication-user-management)
    - [Add IAM user to cluster](#add-iam-user-to-cluster)

## Guide
NAVER Cloud Platform Official Guide provides more detailed and friendly guides.
- [PUB](https://guide.ncloud-docs.com/docs/en/k8s-iam-auth-ncp-iam-authenticator#) (Multilingual support including English)
- [FIN](https://guide-fin.ncloud-docs.com/docs/k8s-iam-auth-ncp-iam-authenticator) (Only Korean)
- [GOV](https://guide-gov.ncloud-docs.com/docs/k8s-iam-auth-ncp-iam-authenticator) (Only Korean)

## Installation
1. Specify version, os and arch you want to use.
    ```bash
    export version="1.0.6" # available versions can be found in GitHub Releases.
    export os="darwin" # darwin, linux, windows
    export arch="amd64" # amd64, arm64
    ```
2. Download the `ncp-iam-authenticator` binary.
    - macOS, Linux
    ```bash
    curl -o ncp-iam-authenticator -L https://github.com/NaverCloudPlatform/ncp-iam-authenticator/releases/download/v${version}/ncp-iam-authenticator_${os}_${arch}
    ```
    - windows (PowerShell)
    ```bash
    curl -o ncp-iam-authenticator -L https://github.com/NaverCloudPlatform/ncp-iam-authenticator/releases/download/v${version}/ncp-iam-authenticator_windows_amd64.exe
    ```
3. (Optional) You can use SHA-256 SUM to check the downloaded binary file.
    1. Check the sum of SHA-256 of the `ncp-iam-authenticator` binary file.
        - macOS, Linux
            ```bash
            openssl sha1 -sha256 ncp-iam-authenticator
            ```
        - windows (PowerShell)
            ```bash
            Get-FileHash ncp-iam-authenticator.exe
            ```
    2. Download SHA-256 SUM.
        ```bash
        curl -o ncp-iam-authenticator.sha256 -L https://github.com/NaverCloudPlatform/ncp-iam-authenticator/releases/download/v${version}/ncp-iam-authenticator_${version}_SHA256SUMS
        ```
    3. Checks if two values match.
4. Set permission and Path
    - macOS, Linux
        1. Add the execution permission to the binary.
            ```bash
            chmod +x ./ncp-iam-authenticator
            ```
        2. Create `$HOME/bin/ncp-iam-authenticator`, and add to `$PATH`.
            ```bash
            mkdir -p $HOME/bin && cp ./ncp-iam-authenticator $HOME/bin/ncp-iam-authenticator && export PATH=$PATH:$HOME/bin
            ```
        3. Add `PATH` to the shell profile.
           - bash
               ```bash
               echo 'export PATH=$PATH:$HOME/bin' >> ~/.bash_profile
               ```
           - zsh
               ```bash
               echo 'export PATH=$PATH:$HOME/bin' >> ~/.zshrc
               ```
    - windows
        1. Create a new folder, such as C:\bin.
        2. Copy the execution file ncp-iam-authenticator.exe to the new folder.
        3. Edit the user or system PATH environment variable to add the new folder to PATH.
        4. Close the PowerShell terminal, and open a new terminal to import the new PATH variable.
5. Test if the `ncp-iam-authenticator` binary works normally.
    ```
    ncp-iam-authenticator help
    ```

## Creating IAM authentication kubeconfig

You can create a kubeconfig through `ncp-iam-authenticator`, or manually create a kubeconfig that uses `ncp-iam-authenticator`, for IAM cluster authentication in Kubernetes Service.

### Set ncp-iam-authenticator API authentication key value

An API authentication key value must first be set up to use `ncp-iam-authenticator`.  
You can get the API authentication key from **[My Page]** > **[Manage account]** > **[Manage authentication key]**  
Set the API key in OS environment variable or configure file. ( OS environment variable takes priority over the configure file.)

- OS environment variable
    ```bash
    export NCLOUD_ACCESS_KEY=ACCESSKEYACCESSKEYAC
    export NCLOUD_SECRET_KEY=SECRETKEYSECRETKEYSECRETKEYSECRETKEYSECR
    export NCLOUD_API_GW=https://ncloud.apigw.ntruss.com
    ```
- The configure file in the user environment home directory's .ncloud folder
    ```bash
    $ cat ~/.ncloud/configure
    [DEFAULT]
    ncloud_access_key_id = ACCESSKEYACCESSKEYAC
    ncloud_secret_access_key = SECRETKEYSECRETKEYSECRETKEYSECRETKEYSECR
    ncloud_api_url = https://ncloud.apigw.ntruss.com
    
    [project]
    ncloud_access_key_id = ACCESSKEYACCESSKEYAC
    ncloud_secret_access_key = SECRETKEYSECRETKEYSECRETKEYSECRETKEYSECR
    ncloud_api_url = https://ncloud.apigw.ntruss.com
    ```

### Use command ncp-iam-authenticator create-kubeconfig
1. Confirm if `ncp-iam-authenticator` has been installed.
2. Use the `ncp-iam-authenticator create-kubeconfig` command to create a kubeconfig for the cluster.
    ```bash
    ncp-iam-authenticator create-kubeconfig --region <region-code> --clusterUuid <cluster-uuid> > kubeconfig.yaml
    ```
    - region-code : Cluster Region code
      ex) KR, SGN
    - cluster-uuid: Cluster UUID
    - If you specify a profile of the NCLOUD CLI configure file with the `--profile` option, then the profile will be used for authentication when the `kubectl` command is executed.
3. Test the `kubectl` command with the `kubeconfig` file created.
    ```bash
    $ kubectl get namespaces --kubeconfig kubeconfig.yaml
    NAME                    STATUS   AGE
    default                 Active   1h
    kube-node-lease         Active   1h
    kube-public             Active   1h
    kube-system             Active   1h
    kubernetes-dashboard    Active   1h
    ```

## IAM authentication user management
When you create a Kubernetes Service cluster, the `SubAccount account that created the cluster` and `main account` will automatically be included in the `system:masters` group in the cluster's RBAC configuration. This configuration is not shown in the cluster information or ConfigMap. In order to give permissions to use a cluster to an IAM user, `ncp-auth` ConfigMap must be registered to the `kube-system` namespace.  
The configuration can be set up after `ncp-iam-authenticator` has been installed and the kubeconfig is created.

### Add IAM user to cluster
1. A `kubectl` credential must already be set up with the `IAM user who created the cluster` or `main account`.
2. Create `ncp-auth` ConfigMap.
    ```
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: ncp-auth
      namespace: kube-system
    data:
      mapSubAccounts: |
        - subAccountIdNo: <iam-user-idno>
          username: <username>
          groups:
            - <groups>
    ```
3. ConfigMap's IAM user parameters are as below.
    - subaccountIdNo: ID number of the IAM user to be added, as can be confirmed from the IAM console
    - username: username to map on the IAM user within Kubernetes
    - groups: list of groups to map users within Kubernetes For more details, refer to [Default roles and role bindings](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#default-roles-and-role-bindings).
4. Check if the IAM user, or the Kubernetes user or user group with a role mapped, is bound to a Kubernetes role by `RoleBinding` or `ClusterRoleBinding`. For more information, refer to [Using RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) in the Kubernetes documents.
- Permission to view resources in all namespaces - The group name is `full-access-group`, and this needs to be mapped to the IAM user groups from `ncp-auth` ConfigMap.
    ```bash
    $ cat <<EOF | kubectl --kubeconfig $KUBE_CONFIG apply -f -
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
     name: full-access-clusterrole
    rules:
    - apiGroups:
      - ""
      resources:
      - nodes
      - namespaces
      - pods
      verbs:
      - get
      - list
    - apiGroups:
      - apps
      resources:
      - deployments
      - daemonsets
      - statefulsets
      - replicasets
      verbs:
      - get
      - list
    - apiGroups:
      - batch
      resources:
      - jobs
      verbs:
      - get
      - list
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
    name: full-access-binding
    subjects:
    - kind: Group
     name: full-access-group
     apiGroup: rbac.authorization.k8s.io
    roleRef:
     kind: ClusterRole
     name: full-access-clusterrole
     apiGroup: rbac.authorization.k8s.io
    EOF
    ```

- Permission to view resources for a specific namespace - The namespace set to the file is `default`, so please specify the namespace you want and modify the result. The group name is `restricted-access-group`, and this needs to be set to IAM user's groups in the `ncp-auth` ConfigMap.
    ```
    $ cat <<EOF | kubectl --kubeconfig $KUBE_CONFIG apply -f -
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      name: restricted-access-clusterrole
    rules:
    - apiGroups:
      - ""
      resources:
      - nodes
      - namespaces
      verbs:
      - get
      - list
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
      name: restricted-access-clusterrole-binding
    subjects:
    - kind: Group
      name: restricted-access-group
      apiGroup: rbac.authorization.k8s.io
    roleRef:
      kind: ClusterRole
      name: restricted-access-clusterrole
      apiGroup: rbac.authorization.k8s.io
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: Role
    metadata:
      namespace: default
      name: restricted-access-role
    rules:
    - apiGroups:
      - ""
      resources:
      - pods
      verbs:
      - get
      - list
    - apiGroups:
      - apps
      resources:
      - deployments
      - daemonsets
      - statefulsets
      - replicasets
      verbs:
      - get
      - list
    - apiGroups:
      - batch
      resources:
      - jobs
      verbs:
      - get
      - list
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: RoleBinding
    metadata:
      name: restricted-access-role-binding
      namespace: default
    subjects:
    - kind: Group
      name: restricted-access-group
      apiGroup: rbac.authorization.k8s.io
    roleRef:
      kind: Role
      name: restricted-access-role
      apiGroup: rbac.authorization.k8s.io
    EOF
    ```
