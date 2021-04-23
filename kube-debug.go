package main

import (
    "fmt"
    "flag"
    "os"
    "os/user"
    "strconv"
    "strings"
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
           kdlib.CheckFileExist(currentdir+"/kube-debug-container-image.tar")
           kdlib.CheckSoft("iptables")
           kdlib.CheckSoft("docker")
           if debugPort != 0 {
               kdlib.CheckPortExist(debugPort)
           } 
           err := kdlib.ShellExecute("sudo docker load < "+currentdir+"/kube-debug-container-image.tar")
           kdlib.CheckErr(err)
           fmt.Println("\nkube-debug Debugging environment initialization completed!\n")

       case container != "" :
           kdContainerName := "kube-debug-container-"+container
           kdlib.CheckSoft("iptables")
           kdlib.CheckSoft("docker")
           kdlib.CheckDebugPort(debugPort)
           kdlib.CheckPortExist(debugPort)
           err1 := kdlib.ShellExecute("sudo docker run --rm -itd --name "+kdContainerName+" --net container:"+container+" --pid container:"+container+" --privileged cloudnativer/kube-debug:"+VERSION)
           kdlib.CheckErr(err1)
           kdContainerIP,err2 := kdlib.ShellOutput("sudo docker exec -i "+kdContainerName+" hostname -i")
           kdlib.CheckErr(err2)
           kdContainerIP = strings.Replace(kdContainerIP, "\n", "", -1)
           err3 := kdlib.ShellExecute("sudo iptables -t nat -A PREROUTING -p tcp -m tcp --dport "+strconv.Itoa(debugPort)+" -j DNAT --to-destination "+kdContainerIP+":3080 -m comment --comment \""+kdContainerName+"\" && sudo iptables -t nat -A POSTROUTING -d "+kdContainerIP+"/32 -p tcp -m tcp --sport 3080 -j SNAT --to-source "+hostIP+" -m comment --comment \""+kdContainerName+"\"")
          kdlib.CheckErr(err3)
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
           kdlib.CheckSshLogin(nodeIP, hostUsername, hostHomedir)
           kdlib.ShellExecute("scp ./kube-debug* "+hostUsername+"@"+nodeIP+":/tmp/")
           if hostUsername != "root" {
               fmt.Println("\nWarning:  If you do not enter the root account, please make sure that the account you are using has sudo permission without entering a password ! \n")
           }
           err1 := kdlib.ShellExecute("ssh "+hostUsername+"@"+nodeIP+" \"/tmp/kube-debug -init -debugport "+strconv.Itoa(debugPort)+" && sudo docker run --rm -itd --name "+kdContainerName+" --net container:"+containerID+" --pid container:"+containerID+" --privileged cloudnativer/kube-debug:"+VERSION+"\"")
           kdlib.CheckErr(err1)
           err2 := kdlib.ShellExecute("ssh "+hostUsername+"@"+nodeIP+" \"sudo iptables -t nat -A PREROUTING -p tcp -m tcp --dport "+strconv.Itoa(debugPort)+" -j DNAT --to-destination "+podIP+":3080 -m comment --comment \"kdContainerName\" && sudo iptables -t nat -A POSTROUTING -d "+podIP+"/32 -p tcp -m tcp --sport 3080 -j SNAT --to-source "+nodeIP+" -m comment --comment \"kdContainerName\" \"")
           kdlib.CheckErr(err2)
           fmt.Println("\nNotice: You can now enter the debugging interface in the following two ways:\n        (1) Using a web browser to access http://k8s-node's_IP:"+strconv.Itoa(debugPort)+" Debug! (Recommended URL: http://"+nodeIP+":"+strconv.Itoa(debugPort)+")\n        (2) Login to the target k8s-node host ("+nodeIP+"), debugging with commands: docker exec -it "+kdContainerName+" /bin/bash \n")

       case node != "" :
           kdContainerName := "kube-debug-node-"+node
           kdlib.CheckDebugPort(debugPort)
           kdlib.CheckIP(node)
           kdlib.CheckSoft("ssh-copy-id")
           kdlib.CheckSoft("scp")
           kdlib.CheckSoft("ssh")
           fmt.Println("Preparing kube-debug environment, Please wait... \n")
           kdlib.CheckSshLogin(node, hostUsername, hostHomedir)
           kdlib.ShellExecute("scp ./kube-debug* "+hostUsername+"@"+node+":/tmp/")
           if hostUsername != "root" {
               fmt.Println("\nWarning:  If you do not use the root account, please make sure that the account you are using has sudo permission without entering a password ! \n")
           }
           err := kdlib.ShellExecute("ssh "+hostUsername+"@"+node+" \"/tmp/kube-debug -init -debugport "+strconv.Itoa(debugPort)+" && sudo docker run --rm -itd --name "+kdContainerName+" --net host --pid host --privileged cloudnativer/kube-debug:"+VERSION+" kube-debug-ttyd -p "+strconv.Itoa(debugPort)+" bash\" ")
           kdlib.CheckErr(err)
           fmt.Println("\nNotice: You can now enter the debugging interface in the following two ways:\n        (1) Using a web browser to access http://k8s-node's_IP:"+strconv.Itoa(debugPort)+" Debug! (Recommended URL: http://"+node+":"+strconv.Itoa(debugPort)+")\n        (2) Login to the target k8s-node host ("+node+"), debugging with commands: docker exec -it "+kdContainerName+" /bin/bash \n")

       case *localhostFlag :
           kdlib.CheckDebugPort(debugPort)
           kdlib.CheckPortExist(debugPort)
           kdlib.CheckSoft("docker")
           err := kdlib.ShellExecute("sudo docker run --rm -itd --name kube-debug-localhost --net host --pid host --privileged cloudnativer/kube-debug:"+VERSION+" kube-debug-ttyd -p "+strconv.Itoa(debugPort)+" bash")
           kdlib.CheckErr(err)
           fmt.Println("\nNotice: You can now enter the debugging interface in the following two ways:\n        (1) Using a web browser to access http://Localhost's_IP:"+strconv.Itoa(debugPort)+" Debug! (Recommended URL: http://"+hostIP+":"+strconv.Itoa(debugPort)+")\n        (2) Use the command to debug directly on the local host: docker exec -it kube-debug-localhost /bin/bash \n")

       case *clearFlag :
           kdlib.CheckSoft("iptables")
           kdlib.CheckSoft("docker")
           fmt.Println("Cleaning up temporary cache files...")
           if hostUsername == "root" {
               os.Remove("/tmp/kube-debug*")
           } else {
               fmt.Println("\nNotice:  If you do not use the root account, you may need to enter the sudo permission password of the local host,")
               kdlib.ShellOutput("sudo rm -rf /tmp/kube-debug*")
           }
           fmt.Println("Cleaning up debug container process...")
           _,err1 := kdlib.ShellOutput("sudo docker ps | grep 'kube-debug' | awk '{print $1}' | xargs sudo docker stop >/dev/null 2>&1")
           if err1 != nil {
               fmt.Println("\nThe container list is empty.\n")
           }
           fmt.Println("Cleaning up port forwarding rule...")
           _,err2 := kdlib.ShellOutput(" iptablesList=`iptables -S -t nat | grep \"kube-debug-\" | sed 's/-A/-t nat -D /g' | sed 's/$/&\\;/g' ` && iptablesNum=`echo $iptablesList | awk -F';' '{print NF-1}'` && if [ $iptablesNum -ge 1 ]; then  for ((i=1;i<=$iptablesNum;i++)); do  echo $iptablesList | cut -d \";\" -f $i | xargs iptables ; done; fi ")

           if err2 != nil {
               fmt.Println("\nWarning:  Port forwarding rule cleaning failed, please manually execute the following command to clean up:\n iptables -S -t nat | grep \"kube-debug-\" | sed 's/-A/iptables -t nat -D /g'")
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



