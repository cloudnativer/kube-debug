# How to build it?

<br>

First, make sure the docker and git commands are installed locally.

<br>

Then, use the following command to download the code locally:

```
# git clone https://github.com/cloudnativer/kube-debug.git
```

Enter `kube-debug` directory, executing the `make` command:

```
# cd kube-debug
# make
```

The above command will automatically complete the build operation. 
After building, We will see two files `kube-debug` and `kube-debug-container-image.tar` in the current directory. `kube-debug` file is binary executable file, `kube-debug-container-image.tar` file is kube-debug's container image offline package.
Now we can copy these two files to other machines, run the command `kube-debug -help` and start the debugging journey.
<br>
<br>
<br>
