package main

import "os"
import "fmt"
import "os/exec"
import "syscall"
import "io/ioutil"
import "path/filepath"

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("unknown command")
	}
}

func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	// HACK: namespaceを分離した新しいプロセスの作成とそのプロセス内でのhostname等の設定を同時に行うことができないた、え
	// まずはnamespaceを分離した新しいプロセスを生成し、そのプロセス内でhostname等の設定を行うようにする
	// /proc/self/exe は現在実行中のバイナリを指す特殊なパス
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// 標準入出力を親と共有（対話型にするため）
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 指定したnamespaceを分離する
	cmd.SysProcAttr = &syscall.SysProcAttr {
		// UTS namespace: hostname, domainname
		// PID namespace: プロセスID空間の分離
		// NS namespace: マウントポイントの分離（マウントテーブルがコピーされるが共有状態は残る）
		// Cloneflagsはビットマスクで複数指定を行う
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		// マウントテーブルの共有を解除
		Unshareflags: syscall.CLONE_NEWNS,
	}
	// プロセスを実行
	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cg()

	// ホスト名を変更
	syscall.Sethostname([]byte("container"))
	// rootファイルシステムを変更
	syscall.Chroot("/vagrant/ubuntu-fs")
	syscall.Chdir("/")
	// proc（procfsの仮想デバイス）を./procにprocfsという種類のファイルシステムとしてマウント
	// source=proc, target=proc, fstype=proc
	syscall.Mount("proc", "proc", "proc", 0, "")

	defer syscall.Unmount("proc", 0)

	// 新しいプロセスの生成を準備
	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	// 標準入出力を親と共有（対話型にするため）
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// プロセスを実行
	must(cmd.Run())
}

func cg() {
	cgroups := "/sys/fs/cgroup/"
	err := os.Mkdir(filepath.Join(cgroups, "container"), 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	// プロセス数の上限を設定（20プロセス）
    must(os.WriteFile(filepath.Join(cgroups, "container/pids.max"), []byte("20"), 0700))
    // メモリ上限を設定（256MB）
    memLimit := 256 * 1024 * 1024
	must(os.WriteFile(filepath.Join(cgroups, "container/memory.max"), []byte(fmt.Sprintf("%d", memLimit)), 0700))
    // CPU を 50% に制限 (quota=50ms, period=100ms)
	must(os.WriteFile(filepath.Join(cgroups, "container/cpu.max"), []byte("50000 100000"), 0700))
	// このcgroupに現在のプロセスを所属させる
	must(ioutil.WriteFile(filepath.Join(cgroups, "container/cgroup.procs"), []byte(fmt.Sprintf("%d", os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
