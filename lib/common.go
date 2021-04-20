package kdlib

import (
    "fmt"
    "os"
    "os/exec"
    "io"
//    "io/ioutil"
    "log"
    "strings"
)



func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}

func CheckParam(paramname string, param string) {
    if param == "" {
         panic("\""+paramname+"\" Parameter cannot be empty, please check!")
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

func ShowHelp(){
    fmt.Println("Usage of kube-debug: kube-debug [COMMAND] { [OBJECT] [ARGS]... } \n\nCOMMAND: \n  init          Initialize the kube-debug environment. \n  localhost     Debug the local host. \n  container     Set the target container ID or container name to be debugged. \n  show          Set the kubernetes pod name to query. \n\nOBJECT: \n  hostport      Custom debug container in the local debugging listening port (the default is TCP 3080 port). (default 3080) \n  namespace     Set the namespace of kubernetes pod to be queried. \n  kubeconfig    Set the kubeconfig file path of kubernetes cluster to be queried. \n\nEXAMPLE: \n  Initialize the kube-debug environment \n    kube-debug -init \n  Debug the local host \n    kube-debug -localhost \n  Debug the target container (container ID is '9a64c7a0d6bd') and open the debug port of tcp-38080 on the local host \n    kube-debug -container \"9a64c7a0d6bd\" -hostport 38080 \n  Query the container ID and kubernetes node IP of 'test-6bfb69dc64-hdblq' pod in 'testns' namespace \n    kube-debug -show \"test-6bfb69dc64-hdblq\" -namespace \"testns\" -kubeconfig \"/etc/kubernetes/pki/kubectl.kubeconfig\" \n")
}

func ShowVersion(){
    fmt.Println("Version 0.1.0 \nRelease Date: 4/23/2021 \n")
}



