
A toolbox for debugging Docker container and Kubernetes with visual Web UI.

<br>

![kube-debug](docs/imgs/kube-debug-logo.jpg)

<br>



# [1] Quick start?

<br>

## 1.1 Download kube-install package file

<br>
You can download the `kube-debug-*.tgz` package from https://github.com/cloudnativer/kube-debug/releases. <br>

For example, we have downloaded the `kube-debug-x86-v0.1.0.tgz` package.<br>

```
# cd ~/
# wget https://github.com/cloudnativer/kube-debug/releases/download/v0.1.0/kube-debug-x86-v0.1.0.tgz
# tar -zxvf kube-debug-x86-v0.1.0.tgz
# cd ~/kube-debug/
```

<br>

## 1.2 Initialize the kube-debug environment:

<br>
Execute the following command to initialize the local debug environment: <br>

```
# cd ~/kube-debug/
# ./kube-debug -init
```

<br>

## 1.3 Start debugging

<br>
You can use `kube-debug` to debug the local host, the local container, any kubernetes node and any kubernetes pod of any namespace.<br>

For example, We can use `kube-debug -container <container id or container name>` to debug the local container. let's debug the local container with ID `9a64c7a0d6bd` , We can perform the following command operations:<br>

```
# cd ~/kube-debug/
# ./kube-debug -container "9a64c7a0d6bd"
```

After the command is executed, the following information will be displayed:

```

xxxxxxxxxx

```

Now you can use a web browser to access `xxxxx` and enter the visual web ui for debugging.

<br>

![kube-debug](docs/imgs/kube-debug-ui-01.jpg)

<br>
Here we can easily do visual debugging, Automatically complete the command line, Rich visual debugging tools:
<br>

![kube-debug](docs/imgs/kube-debug-ui-02.jpg)

<br>

In addition, we can also use the `yyyy` command to directly log in to the debugging container for debugging.

```
xxxxxxxxx

```

<br>
<br>
<br>




# [2] How to debug?

<br>

## 2.1 Debug the local host
<br>

We can use `kube-debug -localhost` to debug the local host. Let's debug the local host , We can perform the following command operations:<br>

```
# cd ~/kube-debug/
# ./kube-debug -localhost
```

After the command is executed, the following information will be displayed:

```

xxxxxxxxxx

```

Now you can use a web browser to access `xxxxx` and enter the Visual Web UI for debugging. At the same time, you can also use the `yyyy` command to directly log in to the debugging container for debugging.
<br>
<br>

## 2.2 Debug the local container

<br>
We can use `kube-debug -container <container id or container name>` to debug the local container. Let's debug the local container with ID `9a64c7a0d6bd` , We can perform the following command operations:<br>

```
# cd ~/kube-debug/
# ./kube-debug -container "9a64c7a0d6bd"
```

After the command is executed, the following information will be displayed:

```

xxxxxxxxxx

```

Now you can use a web browser to access `xxxxx` and enter the Visual Web UI for debugging. At the same time, you can also use the `yyyy` command to directly log in to the debugging container for debugging.
<br>
<br>

## 2.3 Debug the any kubernetes node

<br>
We can use `kube-debug -node <kubernetes node IP>` to debug any kubernetes node. Let's debug the kubernetes node with IP `192.168.1.12` , We can perform the following command operations:<br>

```
# cd ~/kube-debug/
# ./kube-debug -node "192.168.1.12"
```

After the command is executed, the following information will be displayed:

```

xxxxxxxxxx

```

Now you can use a web browser to access `xxxxx` and enter the Visual Web UI for debugging. At the same time, you can also use the `yyyy` command to directly log in to the debugging container for debugging.
<br>
<br>
## 2.4 Debug the any kubernetes pod

<br>
We can use `kube-debug -pod <pod name> -namespace <namespace> -kubeconfig <kubeconfig file>` to debug any kubernetes pod. Let's debug the `test-6bfb69dc64-hdblq` pod in the `testns` namespace , We can perform the following command operations:<br>

```
# cd ~/kube-debug/
# ./kube-debug -pod "test-6bfb69dc64-hdblq" -namespace "testns" -kubeconfig "/etc/kubernetes/pki/kubectl.kubeconfig"
```

After the command is executed, the following information will be displayed:

```

xxxxxxxxxx

```

Now you can use a web browser to access `xxxxx` and enter the Visual Web UI for debugging. At the same time, you can also use the `yyyy` command to directly log in to the debugging container for debugging.
<br>
<br>





# [3] Cleaning up the debug environment

<br>

We can clean up the local machine's debug environment by executing `kube-debug -clear` . The following command will automatically clean up the residual temporary cache file, debug container, debug process and other information.

```
# cd ~/kube-debug/
# ./kube-debug -clear
```

Now the debug environment of the local machine is cleaned up.
<br>
<br>


# [4] Parameter introduction

<br>
The parameters about kube-debug can be viewed using the `kube-debug -help` command. You can also <a href="docs/parameters.md">see more detailed param
eter introduction here</a>.<br>
<br>
<br>

# [5] How to build it?

<br>
The build can be completed automatically by executing the `make` command.You can also <a href="docs/build.md">see more detailed build instructions here</a>.<br>
<br>
<br>
<br>
