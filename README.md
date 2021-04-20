
#A toolbox for debugging Docker container and Kubernetes with Web UI.

<br>

![kube-debug](docs/imgs/kube-debug-logo.jpg)

<br>

  Initialize the kube-debug environment 
    kube-debug -init 
  Debug the local host 
    kube-debug -localhost 
  Debug the target container (container ID is '9a64c7a0d6bd') and open the debug port of tcp-38080 on the local host 
    kube-debug -container "9a64c7a0d6bd" -hostport 38080 
  Query the container ID and kubernetes node IP of 'test-6bfb69dc64-hdblq' pod in 'testns' namespace 
    kube-debug -show "test-6bfb69dc64-hdblq" -namespace "testns" -kubeconfig "/etc/kubernetes/pki/kubectl.kubeconfig" 


Use the browser to access http://hostIP:3080

<br>

![kube-debug](docs/imgs/kube-debug-ui.jpg)

<br>



