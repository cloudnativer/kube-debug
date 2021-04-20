package main

import (
    "fmt"
    "flag"
    "strconv"
    "kube-debug/lib"

)

func main() {

    //Set variable
    const VERSION string = "v0.1.0"
    var show string
    var namespace string
    var kubeconfig string
    var container string
    var hostPort int = 3080

    //Setting the flag of command parameters
    versionFlag := flag.Bool("version",false,"View software version information.")
    helpFlag := flag.Bool("help",false,"View usage help information.")
    initFlag := flag.Bool("init",false,"Initialize the kube-debug environment.")
    hostportFlag := flag.Int("hostport",hostPort,"Custom debug container in the local debugging listening port (the default is TCP 3080 port).")
    localhostFlag := flag.Bool("localhost",false,"Debug the local host.")
    flag.StringVar(&container,"container","","Set the target container ID or container name to be debugged.")
    flag.StringVar(&show,"show","","Set the kubernetes pod name to query.")
    flag.StringVar(&namespace,"namespace","","Set the namespace of kubernetes pod to be queried.")
    flag.StringVar(&kubeconfig,"kubeconfig","","Set the kubeconfig file path of kubernetes cluster to be queried.")
    flag.Parse()
    if *hostportFlag != hostPort {
        hostPort = *hostportFlag
    }



    switch {

       case *initFlag :
           kdlib.CheckFileExist("./images/kube-debug-container-image.tar")
           kdlib.CheckSoft("docker")
           kdlib.ShellExecute("docker load < ./images/kube-debug-container-image.tar")

       case show != "" :
           kdlib.CheckParam("namespace",namespace)
           kdlib.CheckParam("kubeconfig",kubeconfig)
           kdlib.CheckFileExist(kubeconfig)
           nodeIP,containerID := kdlib.GetPod(show, namespace, kubeconfig)
           fmt.Println("Kubernetes Pod Name : "+show+"\nKubernetes Node IP : "+nodeIP+"\nContainer ID : "+containerID+"\n\n---------------------------------------------------------------------------------------\n")
           fmt.Println("Notice! Please go to '"+nodeIP+"' host to execute the following command to start debugging : \n           kube-debug -container \""+containerID+"\" -hostport <your custom debug port> \n")

       case container != "" :
           kdlib.CheckSoft("docker")
           kdlib.ShellExecute("docker run --rm -itd --net container:"+container+" --pid container:"+container+" --privileged -p "+strconv.Itoa(hostPort)+":3080 cloudnativer/kube-debug:"+VERSION)

       case *localhostFlag :
           kdlib.CheckSoft("docker")
           kdlib.ShellExecute("docker run --rm -itd --net host --pid host --privileged cloudnativer/kube-debug:"+VERSION+" ttyd -p "+strconv.Itoa(hostPort)+" bash")

       case *versionFlag :
           kdlib.ShowVersion()

       case *helpFlag :
           kdlib.ShowHelp()

       default:
           fmt.Println("Notice: the command parameter you entered is wrong. Please refer to the help document below and try again after checking! \n")
           kdlib.ShowHelp()

    }


}



