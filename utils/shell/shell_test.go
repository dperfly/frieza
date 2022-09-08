package shell

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestExecShell(t *testing.T) {
	ExecShell("ping www.baidu.com")
}

func TestCommandAndOutputChan(t *testing.T) {
	ch := make(chan string, 10)
	ctx, cancel := context.WithCancel(context.Background())
	go CommandAndOutputChan(ctx, ch, "ping www.baidu.com")
	for {
		select {
		case out := <-ch:
			fmt.Println(out)
		case <-time.After(5 * time.Second):
			cancel()
			return
		}
	}

}
