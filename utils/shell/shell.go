package shell

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func CommandAndOutputChan(ctx context.Context, ch chan string, cmd string) error {
	// 1.传入Chan 和命令, 执行后通过chan实时获取内容。
	// 2.判断ch中是否存在某个关键字
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", cmd) // windows
	} else {
		c = exec.Command("bash", "-c", cmd) // mac or linux
	}
	stdout, err := c.StdoutPipe()
	if err != nil {
		log.Println("run jmeter-slave error", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			// 其实这段去掉程序也会正常运行，只是我们就不知道到底什么时候Command被停止了，而且如果我们需要实时给web端展示输出的话，这里可以作为依据 取消展示
			select {
			// 检测到ctx.Done()之后停止读取
			case <-ctx.Done():
				log.Println("[Frieza-client] close CommandAndOutputChan func")
				return
			default:
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					return
				}
				// 去掉换行符
				readString = strings.ReplaceAll(readString, "\n", "")
				ch <- readString
			}
		}
	}(&wg)
	err = c.Start()
	if err != nil {
		log.Println("run jmeter-slave error", err)
		return err
	}
	wg.Wait()
	return nil
}

func ExecShell(s string) (string, error) {
	// 执行完毕输出所有的内容
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", s) // windows
	} else {
		cmd = exec.Command("bash", "-c", s) // mac or linux
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	return out.String(), err
}
