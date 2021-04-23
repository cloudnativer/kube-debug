package kdlib

import (
    "fmt"
    "net"
    "os"
    "os/exec"
    "io"
    "io/ioutil"
    "log"
    "strconv"
    "strings"
)



func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}

func CheckParamString(paramname string, param string) {
    if param == "" {
        panic("\""+paramname+"\" Parameter cannot be empty, please check!")
    }
}

func CheckDebugPort(debugPort int) {
    if debugPort == 0 {
        panic("debugport cannot be empty, please check!")
    }
}

func CheckPortExist(port int) {
    address := fmt.Sprintf("%s:%d", "0.0.0.0", port)
    //Determine whether the port has been monitored locally.
    listener, err := net.Listen("tcp", address)
    if err != nil {
        panic("debugport "+address+" is taken, please select another port!")
    }
    defer listener.Close()
    //Determine whether the port already exists in the iptables rule.
    iptablesPort,_ := ShellOutput("iptables -S -t nat | grep \"dport "+strconv.Itoa(port)+" -m\" ")
    if iptablesPort != "" {
        panic("debugport "+address+" is taken, please select another port!")
    }
}

func CheckFileExist(filename string) {
    //Judge whether the file is generated
    _, err := os.Stat(filename)
    if err != nil {
        panic("\""+filename+"\" File does not exist, please restore the file manually or contact the administrator!")
    }
    if os.IsNotExist(err) {
        panic("\""+filename+"\" File does not exist, please restore the file manually or contact the administrator!")
    }
}

func CheckSoft(softname string) {
    cmd,err := exec.Command("sh","-c","type "+softname+" >/dev/null 2>&1 || { echo 'n'; }").Output()
    CheckErr(err)
    if string(cmd) != "" {
        if string(cmd[0]) == "n" {
            panic("\""+softname+"\" Component not installed, please manually install component or contact administrator!")
        }
    }
}

func CheckSshLogin(nodeIP string, hostUsername string, hostHomedir string){
    out,_ := ShellOutput("ssh "+nodeIP+" -o PreferredAuthentications=publickey -o StrictHostKeyChecking=no \"ls\" &> /dev/null && if [ $? -eq 0 ]; then echo \"ok\" ; fi")
    if out != "ok\n" {
        _, err := os.Stat(hostHomedir+"/.ssh/id_rsa")
        if err != nil {
            if os.IsNotExist(err) {
                ShellExecute("ssh-keygen -t rsa -P '' -f "+hostHomedir+"/.ssh/id_rsa")
            }
        }
        fmt.Println("Please wait and enter the password of the target k8s-node host,")
        ShellExecute("ssh-copy-id -p 22 "+hostUsername+"@"+nodeIP+" >/dev/null 2>&1")
    }
}

func ShellAsynclog(reader io.ReadCloser) error {
    cache := "" //Cache less than one line of log information
    buf := make([]byte, 2048)
    for {
        num, err := reader.Read(buf)
        if err != nil && err!=io.EOF{
            return err
        }
        if num > 0 {
            b := buf[:num]
            s := strings.Split(string(b), "\n")
            line := strings.Join(s[:len(s)-1], "\n") //Take out the whole line of log
            fmt.Printf("%s%s\n", cache, line)
            cache = s[len(s)-1]
        }
    }
    return nil
}
 
func ShellExecute(shellfile string) error {
    cmd := exec.Command("sh", "-c", shellfile)
    stdout, _ := cmd.StdoutPipe()
    stderr, _ := cmd.StderrPipe()
    if err := cmd.Start(); err != nil {
        log.Printf("Error starting command: %s......", err.Error())
        return err
    }
    go ShellAsynclog(stdout)
    go ShellAsynclog(stderr)
    if err := cmd.Wait(); err != nil {
        log.Printf("Error waiting for command execution: %s......", err.Error())
        return err
    }
    return nil
}

func ShellOutput(strCommand string)(string, error){
    cmd := exec.Command("/bin/bash", "-c", strCommand) 
    stdout, _ := cmd.StdoutPipe()
    if err := cmd.Start(); err != nil{
        return "",err
    }
    out_bytes, _ := ioutil.ReadAll(stdout)
    stdout.Close()
    if err := cmd.Wait(); err != nil {
        return "",err
    }
    return string(out_bytes),nil
}

func ShowHelp(){
    fmt.Println("Usage of kube-debug: kube-debug [COMMAND] { [OBJECT] [ARGS]... } \n\nCOMMAND: \n  init          Initialize the kube-debug environment. \n  localhost     Debug the local host(Listen to TCP-3080 port by default, and the debugport can be modified by '-debug' parameter). \n  container     Set the target container ID or container name to be debugged. \n  pod           Set the kubernetes pod name to query.\n  node          Set the kubernetes node IP to query. \n  clear         Clean up the local host debugging environment. \n  version       View software version information. \n  help          View usage help information. \n\nOBJECT: \n  debugport     Set the debug listening port on the host. \n  namespace     Set the namespace of kubernetes pod to be queried. \n  kubeconfig    Set the kubeconfig file path of kubernetes cluster to be queried. \n\nEXAMPLE: \n  (1) Initialize the kube-debug environment: \n          kube-debug -init \n  (2) Debug the local host: \n          kube-debug -localhost \n  (3) Debug the target container (container ID is '9a64c7a0d6bd') on the local host, and set the debug listening port is TCP-38080: \n          kube-debug -container \"9a64c7a0d6bd\" -debugport 38080 \n  (4) Debug the target k8s-node host (IP is 192.168.1.13), and set the debug listening port is TCP-38081: \n          kube-debug -node \"192.168.1.13\" -debugport 38081 \n  (5) Debug the pod 'test-6bfb69dc64-hdblq' in the 'testns' namespace, and set the debug listening port is TCP-38082: \n          kube-debug -pod \"test-6bfb69dc64-hdblq\" -namespace \"testns\" -kubeconfig \"/etc/kubernetes/pki/kubectl.kubeconfig\" -debugport 38082 \n  (6) Clean up the local host debugging environment: \n          kube-debug -clear \n")
}

func ShowVersion(){
    fmt.Println("Version 0.1.0 \nRelease Date: 4/23/2021 \n")
}



