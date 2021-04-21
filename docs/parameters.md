
```

Usage of kube-debug: kube-debug [COMMAND] { [OBJECT] [ARGS]... } 

COMMAND: 
  init          Initialize the kube-debug environment. 
  localhost     Debug the local host. 
  container     Set the target container ID or container name to be debugged. 
  pod           Set the kubernetes pod name to query.
  node          Set the kubernetes node IP to query. 
  clear         Clean up the local host debugging environment. 
  version       View software version information. 
  help          View usage help information. 

OBJECT: 
  hostport      Custom debug kubernetes node in the local debugging listening port (the default is TCP 3080 port). 
  namespace     Set the namespace of kubernetes pod to be queried. 
  kubeconfig    Set the kubeconfig file path of kubernetes cluster to be queried. 

EXAMPLE: 
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

