package server

import (
	"Frieza/constant"
	"Frieza/master"
	"Frieza/master/factory"
	"Frieza/master/slaves"
	proto "Frieza/pb"
	"Frieza/utils/md5"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func init() {
	factory.Register(&slaves.Slaves{
		SlavesMap: make(map[string]*master.Slave),
	})
}
func NewSlaves() (master.Slaves, error) {
	return factory.New()
}

func StartServer(address string) {
	g := grpc.NewServer()
	proto.RegisterFriezaServer(g, &Server{})
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic("监听错误:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("启动错误:" + err.Error())
	}
}

type Server struct {
	proto.UnimplementedFriezaServer
}

func ServerRecv(streams proto.Frieza_BidirectionalStreamServer, quit chan struct{}, name string, slaves master.Slaves) {
	for {
		recv, err := streams.Recv()
		if err != nil { //"rpc error: code = Canceled desc = context canceled"
			quit <- struct{}{}
			return
		}
		_ = slaves.UpdateIp(name, recv.IP)
		// 将接受到的状态同步到slaves的status
		switch recv.Status {
		case constant.Idle, constant.Stopped:
			err := slaves.UpdateStatus(name, constant.Idle)
			if err != nil {
				fmt.Println(err)
			}
		case constant.Starting:
			err := slaves.UpdateStatus(name, constant.Starting)
			if err != nil {
				fmt.Println(err)
			}
		case constant.Running:
			err := slaves.UpdateStatus(name, constant.Running)
			if err != nil {
				fmt.Println(err)
			}
		case constant.Failed:
			err := slaves.UpdateStatus(name, constant.Failed)
			if err != nil {
				fmt.Println(err)
			}
		default:
			panic(fmt.Sprintf("%d is not found", recv.Status))
		}

	}
}

func (s Server) BidirectionalStream(streams proto.Frieza_BidirectionalStreamServer) error {
	// 接收输入指令并处理
	slaves, err := NewSlaves()
	if err == master.ErrNotFound {
		log.Println(err)
	}
	//slave唯一值标识
	uuid := md5.GetMD5Hash(time.Now().String())
	// 用于控制退出的chan
	quit := make(chan struct{}, 0)

	// 如果不存在的name则创建，理论上这里不用判断，肯定成功
	if !slaves.IsExist(uuid) {
		_ = slaves.Create(uuid)
	}

	// 接收信息打印内容
	go ServerRecv(streams, quit, uuid, slaves)

	slaveCmdChan, err := slaves.Get(uuid)
	if err == master.ErrExist {
		log.Println(err)
		return err
	}
	// 监听quit，slave, 3s自动检测PingBackCmd
	for {
		select {
		case <-quit:
			_ = slaves.Delete(uuid)
			return nil

		case cmd := <-slaveCmdChan:
			err := streams.Send(&proto.Response{
				Cmd: cmd,
			})
			if err != nil {
				return err
			}

		// 自动3s检测客户端状态
		case <-time.After(time.Second * 3):
			err := streams.Send(&proto.Response{
				Cmd: constant.PingBackCmd,
			})
			if err != nil {
				_ = slaves.Delete(uuid)
				return err
			}
		}
	}
}
