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

    "golang.org/x/crypto/ssh/terminal"
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
    iptablesPort,_ := ShellOutput("sudo iptables -S -t nat | grep \"dport "+strconv.Itoa(port)+" -m\" ")
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
    cmd,err := ShellOutput("command -v "+softname+" >/dev/null 2>&1 || { echo 'notok'; }")
    CheckErr(err)
    if cmd != "" {
        panic("\""+softname+"\" Component not installed, please manually install component or contact administrator!")
    }
}

func CheckSshLoginSudo(nodeIP string, hostUsername string, hostHomedir string)(string, string){
    sshLoginOut,_ := ShellOutput("ssh "+hostUsername+"@"+nodeIP+" -o PreferredAuthentications=publickey -o StrictHostKeyChecking=no \"ls\" &> /dev/null && if [ $? -eq 0 ]; then echo \"ok\" ; fi")
    if sshLoginOut != "ok\n" {
        _, err := os.Stat(hostHomedir+"/.ssh/id_rsa")
        if err != nil {
            if os.IsNotExist(err) {
                ShellExecute("ssh-keygen -t rsa -P '' -f "+hostHomedir+"/.ssh/id_rsa")
            }
        }
        fmt.Println("Warning: If you don't get through the local ssh key to the target k8s-node ("+nodeIP+"), please wait and enter the password for target k8s-node!\nPlease wait a moment...")
        ShellExecute("ssh-copy-id -p 22 "+hostUsername+"@"+nodeIP+" >/dev/null 2>&1")
    }
    fmt.Println("Checking target k8s-node ("+nodeIP+"), please wait...")
    var sudoStr string
    var sshStr string = "ssh "
    if hostUsername != "root" {
        sshSudoOut,_ := ShellOutput("ssh "+hostUsername+"@"+nodeIP+" -o PreferredAuthentications=publickey -o StrictHostKeyChecking=no \"sudo ls \" &> /dev/null && if [ $? -eq 0 ]; then echo \"ok\" ; fi")
        if sshSudoOut != "ok\n" {
            fmt.Println("\nWarning: If you do not use the root account, please make sure that the account you are using has sudo permission without entering a password ! \n")
            fmt.Printf("[sudo] password for "+nodeIP+": ")
            passwd,_ := terminal.ReadPassword(0)
            passwdStr := string(passwd[:])
            fmt.Println("\n")
            sudoStr = "echo \""+passwdStr+"\" | sudo -S"
        } else {
            sudoStr = "sudo"
        }
        _,err := ShellOutput("ssh "+hostUsername+"@"+nodeIP+" \" "+sudoStr+" ls \" ")
        if err != nil {
            sshStr = "ssh -tt "
        }
    }
    return sudoStr,sshStr
}

func GenerateRemoteCheck(nodeIP string, hostUsername string, dir string,port int) {
    checkShellFile, err1 := os.Create(dir+"/kube-debug-init.sh")
    CheckErr(err1)
    defer checkShellFile.Close()
    checkShellFile.WriteString(" if [ \"`command -v iptables >/dev/null 2>&1 || { echo 'notok'; }`\" == \"notok\" ]; then \n    echo \"iptables Component not installed, please manually install component or contact administrator!\" && exit 1 \n fi \n if [ \"`command -v docker >/dev/null 2>&1 || { echo 'notok'; }`\" == \"notok\" ]; then \n    echo \"docker Component not installed, please manually install component or contact administrator!\" && exit 1 \n fi \n if [ \"`command -v lsof >/dev/null 2>&1 || { echo 'notok'; }`\" == \"notok\" ]; then \n    echo \"lsof Component not installed, please manually install component or contact administrator!\" && exit 1 \n fi \n if [ `lsof -i:"+strconv.Itoa(port)+" | wc -l` -ne 0 ] || [ `iptables -S -t nat | grep -w '"+strconv.Itoa(port)+"' | wc -l ` -ne 0 ] ; then \n    echo 'debugport "+strconv.Itoa(port)+" is taken, please select another port!' && exit 1 \n fi \n docker load < /tmp/kube-debug-container-image.tar  >/dev/null 2>&1 \n ")
    err2 := os.Chmod(dir+"/kube-debug-init.sh", 0755)
    CheckErr(err2)
    ShellExecute("scp "+dir+"/kube-debug* "+hostUsername+"@"+nodeIP+":/tmp/")
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
 
func ShellExecute(shellfile string)(error){
    cmd := exec.Command("sh", "-c", shellfile)
    //fmt.Println(cmd)
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
    //fmt.Println(cmd)
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



