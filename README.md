
#A toolbox for debugging Docker container and Kubernetes with visual Web UI.

<br>

![kube-debug](docs/imgs/kube-debug-logo.jpg)

<br>


```
  (1) Initialize the kube-debug environment: 
          kube-debug -init 
  (2) Debug the local host: 
          kube-debug -localhost 
  (3) Debug the target container (container ID is '9a64c7a0d6bd') on the local host: 
          kube-debug -container "9a64c7a0d6bd" 
  (4) Debug the target k8s-node host (IP is 192.168.1.13), and the custom debug listening port is 38080: 
          kube-debug -node "192.168.1.13" -hostport 38080 
  (5) Debug the pod 'test-6bfb69dc64-hdblq' in the 'testns' namespace: 
          kube-debug -pod "test-6bfb69dc64-hdblq" -namespace "testns" -kubeconfig "/etc/kubernetes/pki/kubectl.kubeconfig" 
  (6) Clean up the local host debugging environment: 
          kube-debug -clear 
```


Use the browser to access http://hostIP:3080

<br>

![kube-debug](docs/imgs/kube-debug-ui-01.jpg)

<br>

<br>

![kube-debug](docs/imgs/kube-debug-ui-02.jpg)

<br>


<br>
The parameters about kube-debug can be viewed using the `kube-debug -help` command. You can also <a href="docs/parameters.md">see more detailed param
eter introduction here</a>.<br>
<br>


