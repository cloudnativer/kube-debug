package main

import (
    "fmt"
    "flag"
    "os"
    "os/user"
//    "bufio"
    "strconv"
//    "strings"
    "path/filepath"
    "kube-debug/lib"
)

func main() {

    //Set variable
    const VERSION string = "v0.1.0"
    var pod string
    var namespace string
    var kubeconfig string
    var container string
    var node string
    var hostIP string
    var debugPort int = 0
    var hostUsername string
    var hostHomedir string
    var currentdir string

    //Setting the flag of command parameters
    versionFlag := flag.Bool("version",false,"View software version information.")
    helpFlag := flag.Bool("help",false,"View usage help information.")
    initFlag := flag.Bool("init",false,"Initialize the kube-debug environment.")
    debugportFlag := flag.Int("debugport",debugPort,"Custom debug container in the local debugging listening port (the default is TCP 3080 port).")
    localhostFlag := flag.Bool("localhost",false,"Debug the local host.")
    clearFlag := flag.Bool("clear",false,"Clean up the local host debugging environment.")
    flag.StringVar(&node,"node","","Set the kubernetes node IP to query.")
    flag.StringVar(&container,"container","","Set the target container ID or container name to be debugged.")
    flag.StringVar(&pod,"pod","","Set the kubernetes pod name to query.")
    flag.StringVar(&namespace,"namespace","","Set the namespace of kubernetes pod to be queried.")
    flag.StringVar(&kubeconfig,"kubeconfig","","Set the kubeconfig file path of kubernetes cluster to be queried.")
    flag.Parse()
    ip, err := kdlib.ExternalIP()
    kdlib.CheckErr(err)
    hostIP = ip.String()
    if *debugportFlag != debugPort {
        debugPort = *debugportFlag
    }
    u, err := user.Current()
    hostUsername = u.Username
    hostHomedir = u.HomeDir
    path, err := os.Executable()
    kdlib.CheckErr(err)
    currentdir = filepath.Dir(path)


    switch {

       case *initFlag :
           kdlib.InitDebugEnv(currentdir, debugPort)

       case container != "" :
           kdContainerName := "kube-debug-container-"+container
           kdlib.CheckSoft("iptables")
           kdlib.CheckSoft("docker")
           kdlib.CheckDebugPort(debugPort)
           kdlib.CheckPortExist(debugPort)
           kdlib.RunLocalContainer(container, kdContainerName, hostIP, debugPort, VERSION)
           fmt.Println("\nNotice: You can now enter the debugging interface in the following two ways:\n        (1) Using a web browser to access http://Localhost's_IP:"+strconv.Itoa(debugPort)+" Debug! (Recommended URL: http://"+hostIP+":"+strconv.Itoa(debugPort)+")\n        (2) Use the command to debug directly on the local host: docker exec -it "+kdContainerName+" /bin/bash \n")

       case pod != "" :
           kdContainerName := "kube-debug-pod-"+pod
           kdlib.CheckSoft("ssh-copy-id")
           kdlib.CheckSoft("scp")
           kdlib.CheckSoft("ssh")
           kdlib.CheckDebugPort(debugPort)
           kdlib.CheckParamString("namespace",namespace)
           kdlib.CheckParamString("kubeconfig",kubeconfig)
           kdlib.CheckFileExist(kubeconfig)
           nodeIP,podIP,containerID := kdlib.GetPod(pod, namespace, kubeconfig)
           kdlib.CheckIP(nodeIP)
           kdlib.CheckIP(podIP)
           fmt.Println("Preparing kube-debug environment, Please wait... \n")
           sudoStr := kdlib.CheckSshLoginSudo(nodeIP, hostUsername, hostHomedir)
           kdlib.GenerateRemoteCheck(nodeIP, hostUsername, currentdir, debugPort)
           kdlib.RunRemoteContainer(hostUsername,nodeIP,podIP,containerID,kdContainerName,VERSION,debugPort, sudoStr)
           fmt.Println("\nNotice: You can now enter the debugging interface in the following two ways:\n        (1) Using a web browser to access http://k8s-node's_IP:"+strconv.Itoa(debugPort)+" Debug! (Recommended URL: http://"+nodeIP+":"+strconv.Itoa(debugPort)+")\n        (2) Login to the target k8s-node host ("+nodeIP+"), debugging with commands: docker exec -it "+kdContainerName+" /bin/bash \n")

       case node != "" :
           kdContainerName := "kube-debug-node-"+node
           kdlib.CheckDebugPort(debugPort)
           kdlib.CheckIP(node)
           kdlib.CheckSoft("ssh-copy-id")
           kdlib.CheckSoft("scp")
           kdlib.CheckSoft("ssh")
           fmt.Println("Preparing kube-debug environment, Please wait... \n")
           sudoStr := kdlib.CheckSshLoginSudo(node, hostUsername, hostHomedir)
           kdlib.GenerateRemoteCheck(node, hostUsername, currentdir, debugPort)
           kdlib.RunRemoteContainer(hostUsername,node,"","",kdContainerName,VERSION,debugPort, sudoStr)
           fmt.Println("\nNotice: You can now enter the debugging interface in the following two ways:\n        (1) Using a web browser to access http://k8s-node's_IP:"+strconv.Itoa(debugPort)+" Debug! (Recommended URL: http://"+node+":"+strconv.Itoa(debugPort)+")\n        (2) Login to the target k8s-node host ("+node+"), debugging with commands: docker exec -it "+kdContainerName+" /bin/bash \n")

       case *localhostFlag :
           if debugPort == 0 {
               debugPort = 3080
           } 
           kdlib.CheckSoft("docker")
           kdlib.RunLocalContainer("", "", "", debugPort, VERSION)
           fmt.Println("\nNotice: You can now enter the debugging interface in the following two ways:\n        (1) Using a web browser to access http://Localhost's_IP:"+strconv.Itoa(debugPort)+" Debug! (Recommended URL: http://"+hostIP+":"+strconv.Itoa(debugPort)+")\n        (2) Use the command to debug directly on the local host: docker exec -it kube-debug-localhost /bin/bash \n")

       case *clearFlag :
           kdlib.ClearDebuggEnv(hostUsername)

       case *versionFlag :
           kdlib.ShowVersion()

       case *helpFlag :
           kdlib.ShowHelp()

       default:
           fmt.Println("Notice: the command parameter you entered is wrong. Please refer to the help document below and try again after checking! \n")
           kdlib.ShowHelp()

    }


}



