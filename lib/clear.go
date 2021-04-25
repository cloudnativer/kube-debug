package kdlib

import (
    "fmt"
    "os"
)


func ClearDebuggEnv(hostUsername string){
    CheckSoft("iptables")
    CheckSoft("docker")
    fmt.Println("Cleaning up temporary cache files...")
    if hostUsername == "root" {
        os.Remove("/tmp/kube-debug*")
    } else {
        fmt.Println("\nNotice:  If you do not use the root account, you may need to enter the sudo permission password of the local host,")
        ShellOutput("sudo rm -rf /tmp/kube-debug*")
    }
    fmt.Println("Cleaning up debug container process...")
    _,err1 := ShellOutput("sudo docker ps | grep 'kube-debug' | awk '{print $1}' | xargs sudo docker stop >/dev/null 2>&1")
    if err1 != nil {
        fmt.Println("\nThe container list is empty.\n")
    }
    fmt.Println("Cleaning up port forwarding rule...")
     _,err2 := ShellOutput(" iptablesList=`sudo iptables -S -t nat | grep \"kube-debug-\" | sed 's/-A/-t nat -D /g' | sed 's/$/&\\;/g' ` && iptablesNum=`echo $iptablesList | awk -F';' '{print NF-1}'` && if [ $iptablesNum -ge 1 ]; then  for ((i=1;i<=$iptablesNum;i++)); do  echo $iptablesList | cut -d \";\" -f $i | xargs sudo iptables ; done; fi ")
    if err2 != nil {
        fmt.Println("\nWarning:  Port forwarding rule cleaning failed, please manually execute the following command to clean up:\n iptables -S -t nat | grep \"kube-debug-\" | sed 's/-A/iptables -t nat -D /g'")
    }
    fmt.Println("\nLocal debugging environment cleaning completed! \n")
}



