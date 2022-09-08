package slave

import (
	"Frieza/slave/server"
	"Frieza/utils/ip"
	"context"
	"log"
	"os"
)

func StartClient(ctx context.Context, address string) {
	for {
		// 由于ctx chan的缘故，每次都需要new，否则直接ctx.done()始终是close状态
		slaveServer := server.NewSlaveServer(ctx, ip.GetOutBoundIp(), address)
		go func() {
			select {
			case <-slaveServer.Quit:
				log.Println("[Frieza-client] exit...")
				os.Exit(0)
			}
		}()
		slaveServer.Run()
	}

}
