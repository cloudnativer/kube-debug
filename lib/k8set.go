package kdlib

import (
        "context"
        "path/filepath"

        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/client-go/kubernetes"
        "k8s.io/client-go/tools/clientcmd"
)


//k8sClientSet
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

// Get Pod containerID
func GetPod(podName string, nameSpace string, kubeConfig string) (string, string, string) {
        ctx := context.TODO()
        getresult, err := k8sClientSet(kubeConfig).CoreV1().Pods(nameSpace).Get(ctx, podName, metav1.GetOptions{})
        if err != nil {
                panic(err)
        }
        //fmt.Println(getresult.Status.ContainerStatuses[0].ContainerID[9:])
        //fmt.Println(getresult.Status.HostIP)
        //containerID := "7af8850a10339c2b5799d627d6324d007008a214f32f39ec21b3304a10c01b23"
        return getresult.Status.HostIP,getresult.Status.PodIP,getresult.Status.ContainerStatuses[0].ContainerID[9:]
}



