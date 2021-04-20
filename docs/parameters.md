Usage of kube-debug: kube-debug [COMMAND] { [OBJECT] [ARGS]... } 

COMMAND: 
  init          Initialize the kube-debug environment. 
  localhost     Debug the local host. 
  container     Set the target container ID or container name to be debugged. 
  show          Set the kubernetes pod name to query. 

OBJECT: 
  hostport      Custom debug container in the local debugging listening port (the default is TCP 3080 port). (default 3080) 
  namespace     Set the namespace of kubernetes pod to be queried. 
  kubeconfig    Set the kubeconfig file path of kubernetes cluster to be queried. 

EXAMPLE: 
  Initialize the kube-debug environment 
    kube-debug -init 
  Debug the local host 
    kube-debug -localhost 
  Debug the target container (container ID is '9a64c7a0d6bd') and open the debug port of tcp-38080 on the local host 
    kube-debug -container "9a64c7a0d6bd" -hostport 38080 
  Query the container ID and kubernetes node IP of 'test-6bfb69dc64-hdblq' pod in 'testns' namespace 
    kube-debug -show "test-6bfb69dc64-hdblq" -namespace "testns" -kubeconfig "/etc/kubernetes/pki/kubectl.kubeconfig" 


