package main

import (
    "fmt"
    "flag"
    "os"
    "os/user"
    "strconv"
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
    var hostPort int = 3080
    var hostUsername string
    var hostHomedir string
    var currentdir string

    //Setting the flag of command parameters
    versionFlag := flag.Bool("version",false,"View software version information.")
    helpFlag := flag.Bool("help",false,"View usage help information.")
    initFlag := flag.Bool("init",false,"Initialize the kube-debug environment.")
    hostportFlag := flag.Int("hostport",hostPort,"Custom debug container in the local debugging listening port (the default is TCP 3080 port).")
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
    if *hostportFlag != hostPort {
        hostPort = *hostportFlag
    }
    u, err := user.Current()
    hostUsername = u.Username
    hostHomedir = u.HomeDir
    path, err := os.Executable()
    kdlib.CheckErr(err)
    currentdir = filepath.Dir(path)


    switch {

       case *initFlag :
           kdlib.CheckFileExist(currentdir+"/kube-debug-container-image.tar")
           kdlib.CheckSoft("docker")
           err := kdlib.ShellExecute("sudo docker load < "+currentdir+"/kube-debug-container-image.tar")
           kdlib.CheckErr(err)
           fmt.Println("\nkube-debug Debugging environment initialization completed!\n")

       case container != "" :
           kdContainerName := "kube-debug-container-"+container
           kdlib.CheckSoft("docker")
           err := kdlib.ShellExecute("sudo docker run --rm -itd --name "+kdContainerName+" --net container:"+container+" --pid container:"+container+" --privileged cloudnativer/kube-debug:"+VERSION+" kube-debug-ttyd -p "+strconv.Itoa(hostPort)+" bash")
           kdlib.CheckErr(err)
           fmt.Println("\nNotice: You can now enter the debugging interface in the following two ways:\n        (1) Using a web browser to access http://Container's IP:"+strconv.Itoa(hostPort)+" Debug! (you can use the 'docker inspect' command to view the container's IP) \n        (2) Use the command to debug directly on the local host: docker exec -it "+kdContainerName+" /bin/bash \n")

       case pod != "" :
           kdlib.CheckSoft("ssh-copy-id")
           kdlib.CheckSoft("scp")
           kdlib.CheckSoft("ssh")
           kdContainerName := "kube-debug-pod-"+pod
           kdlib.CheckParam("namespace",namespace)
           kdlib.CheckParam("kubeconfig",kubeconfig)
           kdlib.CheckFileExist(kubeconfig)
           nodeIP,podIP,containerID := kdlib.GetPod(pod, namespace, kubeconfig)
           kdlib.CheckIP(nodeIP)
           kdlib.CheckIP(podIP)
           fmt.Println("Preparing kube-debug environment, Please wait! \n")
           kdlib.CheckSshLogin(nodeIP, hostUsername, hostHomedir)
           kdlib.ShellExecute("scp ./kube-debug* "+hostUsername+"@"+nodeIP+":/tmp/")
           if hostUsername != "root" {
           fmt.Println("\nWarning:  If you do not enter the root account, please make sure that the account you are using has sudo permission without entering a password ! \n")
           }
           err := kdlib.ShellExecute("ssh "+hostUsername+"@"+nodeIP+" \"/tmp/kube-debug -init && sudo docker run --rm -itd --name "+kdContainerName+" --net container:"+containerID+" --pid container:"+containerID+" --privileged cloudnativer/kube-debug:"+VERSION+"  kube-debug-ttyd -p "+strconv.Itoa(hostPort)+" bash\" ")
           kdlib.CheckErr(err)
           fmt.Println("\nNotice: You can now enter the debugging interface in the following two ways:\n        (1) Using a web browser to access http://k8s-pod's IP:"+strconv.Itoa(hostPort)+" Debug! (Recommended URL: http://"+podIP+":"+strconv.Itoa(hostPort)+")\n        (2) Login to the target k8s-node host ("+nodeIP+"), debugging with commands: docker exec -it "+kdContainerName+" /bin/bash \n")

       case node != "" :
           kdContainerName := "kube-debug-node-"+node
           kdlib.CheckIP(node)
           kdlib.CheckSoft("ssh-copy-id")
           kdlib.CheckSoft("scp")
           kdlib.CheckSoft("ssh")
           fmt.Println("Preparing kube-debug environment, Please wait! \n")
           kdlib.CheckSshLogin(node, hostUsername, hostHomedir)
           kdlib.ShellExecute("scp ./kube-debug* "+hostUsername+"@"+node+":/tmp/")
           if hostUsername != "root" {
           fmt.Println("\nWarning:  If you do not use the root account, please make sure that the account you are using has sudo permission without entering a password ! \n")
           }
           err := kdlib.ShellExecute("ssh "+hostUsername+"@"+node+" \"/tmp/kube-debug -init && sudo docker run --rm -itd --name "+kdContainerName+" --net host --pid host --privileged cloudnativer/kube-debug:"+VERSION+" kube-debug-ttyd -p "+strconv.Itoa(hostPort)+" bash\" ")
           kdlib.CheckErr(err)
           fmt.Println("\nNotice: You can now enter the debugging interface in the following two ways:\n        (1) Using a web browser to access http://k8s-node's IP:"+strconv.Itoa(hostPort)+" Debug! (Recommended URL: http://"+node+":"+strconv.Itoa(hostPort)+")\n        (2) Login to the target k8s-node host ("+node+"), debugging with commands: docker exec -it "+kdContainerName+" /bin/bash \n")

       case *localhostFlag :
           kdlib.CheckSoft("docker")
           err := kdlib.ShellExecute("sudo docker run --rm -itd --name kube-debug-localhost --net host --pid host --privileged cloudnativer/kube-debug:"+VERSION+" kube-debug-ttyd -p "+strconv.Itoa(hostPort)+" bash")
           kdlib.CheckErr(err)
           fmt.Println("\nNotice: You can now enter the debugging interface in the following two ways:\n        (1) Using a web browser to access http://Localhost's IP:"+strconv.Itoa(hostPort)+" Debug! (Recommended URL: http://"+hostIP+":"+strconv.Itoa(hostPort)+")\n        (2) Use the command to debug directly on the local host: docker exec -it kube-debug-localhost /bin/bash \n")

       case *clearFlag :
           kdlib.CheckSoft("docker")
           fmt.Println("Cleaning up temporary cache files.")
           if hostUsername == "root" {
               os.Remove("/tmp/kube-debug*")
           } else {
               fmt.Println("\nNotice:  If you do not use the root account, you may need to enter the sudo permission password of the local host,")
               kdlib.ShellOutput("sudo rm -rf /tmp/kube-debug*")
           }
           _,err := kdlib.ShellOutput("sudo docker ps | grep 'kube-debug' | awk '{print $1}' | xargs sudo docker stop >/dev/null 2>&1")
           if err != nil {
           fmt.Println("\nThe container list is empty.\n")
           }
           fmt.Println("\nLocal debugging environment cleaning completed! \n")

       case *versionFlag :
           kdlib.ShowVersion()

       case *helpFlag :
           kdlib.ShowHelp()

       default:
           fmt.Println("Notice: the command parameter you entered is wrong. Please refer to the help document below and try again after checking! \n")
           kdlib.ShowHelp()

    }


}



