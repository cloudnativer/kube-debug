package kdlib

import (
    "fmt"
)


func InitDebugEnv(dir string, debugPort int){
    CheckFileExist(dir+"/kube-debug-container-image.tar")
    CheckSoft("iptables")
    CheckSoft("docker")
    if debugPort != 0 {
    CheckPortExist(debugPort)
    }
    err := ShellExecute("sudo docker load < "+dir+"/kube-debug-container-image.tar")
    CheckErr(err)
    fmt.Println("\nkube-debug Debugging environment initialization completed!\n")
}



