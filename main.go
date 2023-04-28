package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// go run main.go run cmd args
func main() {

	fmt.Println(os.Args[1])
	switch os.Args[1] {
	case "run":
		//creates a new namespace where you can manipulate the hostname
		run()

	case "child":
		child()
	default:
		fmt.Println("wrong command")

	}
}
func run() {
	//starts a new process and runs the command in the new process
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	//creates a new process and allocates the unix timesharing system namespace
	cmd.SysProcAttr = &syscall.SysProcAttr{

		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
		//used to issolate mounts
		Unshareflags: syscall.CLONE_NEWNS,
	}
	must(cmd.Run())
}

func child() {
	fmt.Printf("running %v from processID %v \n", os.Args[2:], os.Getpid())
	// cg()
	//runs the actual
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	//mount 
	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot("/home/brayo/rootfs"))
	must(syscall.Chdir("/"))
	//to do
	cg()

	cmd2 := exec.Command("pwd")
	cmd2.Stderr = os.Stderr
	cmd2.Stdin = os.Stdin
	cmd2.Stdout = os.Stdout
	cmd2.Run()

	//PID file
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	must(cmd.Run())
	must(syscall.Unmount("proc", 0))
}

func cg() {
	//contains information about the resources accessible to this particular namespace
	cgroups := "/sys/fs/cgroup"
	//pids is the limit on the number of processes
	pids := filepath.Join(cgroups, "pids")

	must(os.MkdirAll(filepath.Join(pids, "brayo"), 0755))

	//it limits the number of processes to 20
	must(ioutil.WriteFile(filepath.Join(pids, "brayo/pids.max"), []byte("20"), 0700))

	must(ioutil.WriteFile(filepath.Join(pids, "brayo/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "brayo/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}


func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//fork and excec
// /proc is a pseudo file system, a pseudo file is a files ystem tha facilitates communication between the userspace and kernel space
// ps goes to /proc to find out the running process
// initially all the processes are viewed as part of the hosts processes even the ones in the newly created namespace
// we have to create a new process namespace with a different filesystem root(with its own copy of slash proc and filesystem).
