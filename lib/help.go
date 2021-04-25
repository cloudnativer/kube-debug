package kdlib

import (
    "fmt"
)


func ShowHelp(){
    fmt.Println("Usage of kube-debug: kube-debug [COMMAND] { [OBJECT] [ARGS]... } \n\nCOMMAND: \n  init          Initialize the kube-debug environment. \n  localhost     Debug the local host(Listen to TCP-3080 port by default, and the debugport can be modified by '-debug' parameter). \n  container     Set the target container ID or container name to be debugged. \n  pod           Set the kubernetes pod name to query.\n  node          Set the kubernetes node IP to query. \n  clear         Clean up the local host debugging environment. \n  version       View software version information. \n  help          View usage help information. \n\nOBJECT: \n  debugport     Set the debug listening port on the host. \n  namespace     Set the namespace of kubernetes pod to be queried. \n  kubeconfig    Set the kubeconfig file path of kubernetes cluster to be queried. \n\nEXAMPLE: \n  (1) Initialize the kube-debug environment: \n          kube-debug -init \n  (2) Debug the local host: \n          kube-debug -localhost \n  (3) Debug the target container (container ID is '9a64c7a0d6bd') on the local host, and set the debug listening port is TCP-38080: \n          kube-debug -container \"9a64c7a0d6bd\" -debugport 38080 \n  (4) Debug the target k8s-node host (IP is 192.168.1.13), and set the debug listening port is TCP-38081: \n          kube-debug -node \"192.168.1.13\" -debugport 38081 \n  (5) Debug the pod 'test-6bfb69dc64-hdblq' in the 'testns' namespace, and set the debug listening port is TCP-38082: \n          kube-debug -pod \"test-6bfb69dc64-hdblq\" -namespace \"testns\" -kubeconfig \"/etc/kubernetes/pki/kubectl.kubeconfig\" -debugport 38082 \n  (6) Clean up the local host debugging environment: \n          kube-debug -clear \n")
}

func ShowVersion(){
    fmt.Println("Version 0.1.0 \nRelease Date: 4/23/2021 \n")
}



