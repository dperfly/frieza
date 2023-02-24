package server

import (
	"Frieza/constant"
	"Frieza/jmeter"
	proto "Frieza/pb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

type SlaveServer struct {
	// 检查jmeter-slave的状态
	// 如果进程存在说明是在工作，否则就是stopped状态
	// 判断jmeter-server.log文件内容, 包含started,说明正在运行中 Running
	// 最后一行包含Created remote object说明刚启动 或者 包含finished说明运行完毕， 两者都处于空闲状态 Idle
	Status  int32
	Message string

	Cancel context.CancelFunc
	ctx    context.Context

	Quit   chan struct{}
	output chan string

	slaveIP string
	address string

	streams proto.Frieza_BidirectionalStreamClient
	conn    *grpc.ClientConn
}

func NewSlaveServer(ctx context.Context, slaveIP string, address string) *SlaveServer {
	c, cancel := context.WithCancel(ctx)
	return &SlaveServer{
		Status:  constant.Starting,
		ctx:     c,
		Cancel:  cancel,
		output:  make(chan string, 1),
		slaveIP: slaveIP,
		address: address,
		Quit:    make(chan struct{}, 1),
	}

}

func (s *SlaveServer) UpdateSlaveStatus() {
	for {
		select {
		case out := <-s.output:
			if strings.HasPrefix(out, "Starting") {
				log.Println("[jmeter-server]", out)
				//空闲状态接收到开始信息后处于运行中
				s.Status = constant.Running
			} else if strings.HasPrefix(out, "Finished") && s.Status == constant.Running {
				log.Println("[jmeter-server]", out)
				//处于运行中，运行结束后恢复到空闲状态
				s.Status = constant.Idle
			} else if strings.HasPrefix(out, "Created remote object") && s.Status == constant.Starting {
				log.Println("[jmeter-server]", out)
				// 首次启动成功后进入空闲状态
				s.Status = constant.Idle
			} else if strings.HasPrefix(out, "An error occurred") && s.Status == constant.Starting {
				// 启动失败
				log.Println("[jmeter-server]", out)

				//set rmi
				if !jmeter.IsOpenServerSSL() {
					err := jmeter.SetServerSSL()
					if err != nil || !jmeter.IsOpenServerSSL() {
						s.Status = constant.Failed
						continue
					}
				}
				// Set the network card used
				//if !jmeter.IsSetServerHostname() {
				//	err := jmeter.SetServerHostname()
				//	if err != nil || !jmeter.IsSetServerHostname() {
				//		s.Status = constant.Failed
				//		continue
				//	}
				//}

				s.Cancel()
				// 一 错误处理尝试恢复
				// 1.Listen failed on port: 0
				// 2.org.apache.log.Logger
			} else if strings.HasPrefix(out, "Neither the JAVA_HOME") && s.Status == constant.Starting {
				log.Println("[jmeter-server]", out)
				log.Println("try run >>>", jmeter.InstallJDKCmdCentOS)
				err := jmeter.InstallJDK()
				if err != nil || !jmeter.IsFindJavaEnv() {
					log.Println("install jdk failed")
					s.Status = constant.Failed
					continue
				}
				log.Println("install jdk success")
				log.Println("restart...")
				s.Cancel()
			}

		case <-s.ctx.Done():
			return

		case <-time.After(time.Second * 3):
		}
	}
}

func (s *SlaveServer) isFailedStop() {
	// 碰到failed 停止
	for {
		select {
		case <-time.After(3 * time.Second):
			switch s.Status {
			case constant.Failed:
				log.Println("[Frieza-client] Failed ; stopping server...")
				log.Println("[Frieza-client] The program is about to exit...")
				s.Quit <- struct{}{}
				s.Cancel()
				return
			}
		case <-s.ctx.Done():
			return
		}

	}
}

func (s *SlaveServer) interruptCheck() {
	c := make(chan os.Signal, 1) //  fix : misuse of unbuffered os.Signal channel as argument to signal.Notify
	signal.Notify(c, os.Interrupt)
	select {
	case sig := <-c:
		log.Printf("[Frieza-client] Got %s signal. Aborting...\n", sig)
		jmeter.New().StopJmeterThread()
		//jmeter.StopJmeterThread()
		s.Quit <- struct{}{}
		s.Cancel()
	}
}

func (s *SlaveServer) Done() <-chan struct{} {
	return s.ctx.Done()
}

func (s *SlaveServer) String() string {
	switch s.Status {
	case constant.Stopped:
		return "Stopped"
	case constant.Failed:
		return "Failed"
	case constant.Idle:
		return "Idle"
	case constant.Starting:
		return "Starting"
	default:
		panic("slave status not found")
	}
}

func (s *SlaveServer) dial() {
	log.Println("The service is starting. Please wait...")
	conn, err := grpc.Dial(s.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		s.Cancel()
		return
	}
	s.conn = conn
}

func (s *SlaveServer) reCheckIsOk() {
	client := proto.NewFriezaClient(s.conn)
	for {
		select {
		case <-time.After(3 * time.Second):
			streams, err := client.BidirectionalStream(context.Background())
			if err != nil {
				//log.Println(err)
				log.Println("Attempting to reconnect to master...")
			} else {
				log.Println("find Frieza master success...")
				s.streams = streams
				return
			}
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *SlaveServer) recvRunCmd() {
	//jm := jmeter.New()
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			res, err := s.streams.Recv()
			if err != nil {
				log.Println("master status err...", err)
				s.Cancel()
			}
			switch res.Cmd {
			case constant.StopJmeterRunCmd:
				// Running 时才会触发
				if s.Status == constant.Running {
					// 1. 停止 2.启动
					//jm.StopJmeterThread()
					//s.Status = constant.Stopped
					//jm.StartSlave(s.ctx, s.output, s.slaveIP)
					//s.Status = constant.Idle
					//log.Println("restart jmeter-server success...wait master send jmeter cmd.")

					// 直接重启slave服务即可，重启后自动连接上master
					s.Cancel()
				} else {
					log.Printf("StopJmeterRunCmd failed ,because of slave status is %s\n", s)
				}

			case constant.PingBackCmd:
				// ping back 不用处理，直接发送现有的内容
				// pass
			case constant.ServerQuit:
				s.Cancel()
			}

			err = s.streams.Send(&proto.Requests{
				IP:      s.slaveIP,
				Status:  s.Status,
				Message: s.Message,
			})
			if err != nil {
				log.Println(err)
			}
		}
	}

}

func (s SlaveServer) close() {
	var err error
	err = s.conn.Close()
	if err != nil {
		log.Println(err)
	}
	err = s.streams.CloseSend()
	if err != nil {
		log.Println(err)
	}
}

func (s *SlaveServer) pingToFriezaMaster() {
	jm := jmeter.New()
	// 首次启动jmeter-server,发送ping请求
	// TODO 这里应该不需要检测ctx.Done()
	select {
	case <-time.After(1 * time.Second):
		jm.StartSlave(s.ctx, s.output, s.slaveIP)
		err := s.streams.Send(&proto.Requests{
			IP:      s.slaveIP,
			Status:  s.Status,
			Message: s.Message,
		})
		if err != nil {
			log.Println(err)
		} else {
			return
		}
	case <-s.ctx.Done():
		log.Println("[Frieza-client] close pingToFriezaMaster func")
		return
	}
}

func (s *SlaveServer) Run() {
	s.dial()
	s.reCheckIsOk()
	s.pingToFriezaMaster()
	defer s.close()

	//监听中断ctrl+c
	go func() {
		s.interruptCheck()
	}()
	// 根据output输出内容更新状态
	sy := sync.WaitGroup{}
	sy.Add(3)
	go func() {
		s.UpdateSlaveStatus()
		defer sy.Done()
		log.Println("[Frieza-client] close UpdateSlaveStatus func")
	}()
	//监听master命令
	go func() {
		s.recvRunCmd()
		defer sy.Done()
		log.Println("[Frieza-client] close recvRunCmd func")

	}()
	// 监听失败信息
	go func() {
		s.isFailedStop()
		defer sy.Done()
		log.Println("[Frieza-client] close isFailedStop func")
	}()

	sy.Wait()
	log.Println("[Frieza-client] close ALL func ok")
}
