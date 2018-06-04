# kr - Kube Reporter

This is a small client-side tool to get information about Kubernetes pods by running the `kubectl` command (in parallel) against different environments for a given project.

### How to use

build the program by running `go build kr.go`

Then execute the command `./kr myproject` where you have k8s pods named something like `myproject-42f3j4ifh43u`. 
If you did edit the `defaultProjectName` variable before compiling, then you can also avoid adding the project name and simply run `./kr` to get info about the project. 

What you get is as in the following picture:

![KR Screenshot](https://github.com/gabrielefalace/kr/blob/master/misc/kr_pic.jpg)
