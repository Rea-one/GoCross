package main

import (
	"net"
	"time"
)

type Listener interface {
	Init(chan task, chan task)
	Start()
	Stop()
}

type listener struct {
	rcv_id_pool_ mQueue[int]
	address_     string
	receivers_   *mList[receiver]
	wait_req_    chan []bool
	release_     chan int
	quit_        chan struct{}
	listener_    net.Listener
	ipasser_     chan task
	opasser_     chan task
	current_rs   int
	max_rs_      int
}

func (tar *listener) Init(ip chan task, op chan task) {
	tar.current_rs = 0
	tar.max_rs_ = 4
	tar.ipasser_ = ip
	tar.opasser_ = op
	tar.rcv_id_pool_.Init()
	for i := 0; i < tar.max_rs_; i++ {
		tar.rcv_id_pool_.Push(i)
	}
}

func (tar *listener) Start() {
	var err error

	tar.listener_, err = net.Listen("tcp", tar.address_)
	if err != nil {
		// 处理监听错误
		return
	}
	defer tar.listener_.Close()
	tar.quit_ = make(chan struct{})
	go func() {
		for {
			select {
			case <-tar.quit_:
				return
			default:
				if tar.current_rs >= tar.max_rs_ {
					time.Sleep(time.Second * 1)
					continue
				}
				// 接受连接
				conn, err := tar.listener_.Accept()
				if err != nil {
					// 处理连接错误
					continue
				}

				// 创建接收者
				var rcv receiver
				rcv.Init(conn, tar.ipasser_, tar.opasser_)
				var node mListNode[receiver]
				node.Init(rcv)
				rcv.Start()

				// 添加接收者到链表
				tar.receivers_.Push_tail(&node)
			}
		}
	}()
}

// 停止监听方法
func (tar *listener) Stop() {
	close(tar.quit_)
	if tar.listener_ != nil {
		tar.listener_.Close()
	}
}

// func sort_receivers(receivers *[]receiver, left int, right int) {
// 	if left >= right {
// 		return
// 	}
// 	var l = left
// 	var r = right
// 	for l <= r {
// 		for receivers[l].conn_.RemoteAddr().String() < receivers[(left+right)/2].conn_.RemoteAddr().String() {
// 			l++
// 		}
// 		for receivers[r].conn_.RemoteAddr().String() > receivers[(left+right)/2].conn_.RemoteAddr().String() {
// 			r--
// 		}
// 		if l <= r {
// 			receivers[l], receivers[r] = receivers[r], receivers[l]
// 			l++
// 			r--

// 		}
// 	}
// 	sort_receivers(receivers, left, r)
// 	sort_receivers(receivers, l, right)
// }
