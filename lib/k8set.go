package kdlib

import (
    "fmt"
    "strconv"
    "strings"
    "context"
    "path/filepath"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)


// k8sClientSet
func k8sClientSet(kubeConfig string) *kubernetes.Clientset {
    config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(kubeConfig))
    if err != nil {
        panic(err)
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }
    return clientset
}

// Get Pod containerID and IP
func GetPod(podName string, nameSpace string, kubeConfig string) (string, string, string) {
    ctx := context.TODO()
    getresult, err := k8sClientSet(kubeConfig).CoreV1().Pods(nameSpace).Get(ctx, podName, metav1.GetOptions{})
    if err != nil {
        panic(err)
    }
    return getresult.Status.HostIP,getresult.Status.PodIP,getresult.Status.ContainerStatuses[0].ContainerID[9:]
}

func RunRemoteContainer(hostUsername string, nodeIP string, podIP string, containerID string, kdContainerName string, version string, debugPort int, sudoStr string, sshStr string) {
    err1 := ShellExecute(sshStr+hostUsername+"@"+nodeIP+" \" "+sudoStr+" /tmp/kube-debug-init.sh \" ")
    CheckErr(err1)
    containerStr := "host"
    containerCmd := ""
    if ( containerID != "" && podIP != nodeIP ) {
        containerStr = "container:"+containerID
    } else {
        containerCmd = "kube-debug-ttyd -p "+strconv.Itoa(debugPort)+" bash"
    }
    err2 := ShellExecute(sshStr+hostUsername+"@"+nodeIP+" \" "+sudoStr+" docker run --rm -itd --name \""+kdContainerName+"\" --net "+containerStr+" --pid "+containerStr+" --privileged cloudnativer/kube-debug:"+version+" "+containerCmd+" \" ")
    CheckErr(err2)
    if ( podIP != "" && podIP != nodeIP ) {
        err3 := ShellExecute("ssh "+hostUsername+"@"+nodeIP+" \" "+sudoStr+" iptables -t nat -A PREROUTING -p tcp -m tcp --dport "+strconv.Itoa(debugPort)+" -j DNAT --to-destination "+podIP+":3080 -m comment --comment \""+kdContainerName+"\" && "+sudoStr+" iptables -t nat -A POSTROUTING -d "+podIP+"/32 -p tcp -m tcp --sport 3080 -j SNAT --to-source "+nodeIP+" -m comment --comment \""+kdContainerName+"\" \"")
        CheckErr(err3)
    }
    fmt.Println("\nDebug environment started successfully!")
}

func RunLocalContainer(containerID string, kdContainerName string, hostIP string, debugPort int, version string) {
    if ( containerID != "" && hostIP != "") {
        err1 := ShellExecute("sudo docker run --rm -itd --name "+kdContainerName+" --net container:"+containerID+" --pid container:"+containerID+" --privileged cloudnativer/kube-debug:"+version+" kube-debug-ttyd -p "+strconv.Itoa(debugPort)+" bash")
        CheckErr(err1)
        kdContainerIP,err2 := ShellOutput("sudo docker exec -i "+kdContainerName+" hostname -i")
        CheckErr(err2)
        kdContainerIP = strings.Replace(kdContainerIP, "\n", "", -1)
        if kdContainerIP == "" {
            panic("Can not get container IP address, it may be a system permission problem!")
        } else {
            if !( kdContainerIP == hostIP || kdContainerIP == "127.0.0.1 ::1" || kdContainerIP == "127.0.0.1" ) {
                err3 := ShellExecute("sudo iptables -t nat -A PREROUTING -p tcp -m tcp --dport "+strconv.Itoa(debugPort)+" -j DNAT --to-destination "+kdContainerIP+":"+strconv.Itoa(debugPort)+" -m comment --comment \""+kdContainerName+"\" && sudo iptables -t nat -A POSTROUTING -d "+kdContainerIP+"/32 -p tcp -m tcp --sport "+strconv.Itoa(debugPort)+" -j SNAT --to-source "+hostIP+" -m comment --comment \""+kdContainerName+"\"")
                CheckErr(err3)
            }
        }
    } else {
        err4 := ShellExecute("sudo docker run --rm -itd --name kube-debug-localhost --net host --pid host --privileged cloudnativer/kube-debug:"+version+" kube-debug-ttyd -p "+strconv.Itoa(debugPort)+" bash")
        CheckErr(err4)
    }
    fmt.Println("\nDebug environment started successfully!")
}


