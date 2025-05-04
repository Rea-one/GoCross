package gocross

import (
	"fmt"
	"net"
	"time"
)

type Listener interface {
	Init(chan task, *iomap, *string)
	Start()
	Stop()
}

type listener struct {
	rcv_id_pool_ mQueue[int]
	address_     string
	receivers_   *mList[*receiver]
	rcv_map_     *map[int]*mListNode[*receiver]
	io_map_      *iomap
	quit_        chan struct{}
	listener_    net.Listener
	release_     chan int
	signal_      chan string
	current_rs   int
	max_rs_      int
}

func (tar *listener) Init(signal chan string, iom *iomap, host *string) {
	tar.current_rs = 0
	tar.max_rs_ = 4
	tar.address_ = *host
	tar.signal_ = signal
	tar.io_map_ = iom
	tar.rcv_id_pool_.Init()
	for i := range tar.max_rs_ {
		tar.rcv_id_pool_.Push(i)
	}
}

func (tar *listener) Start() {
	var err error

	tar.listener_, err = net.Listen("tcp", tar.address_)
	fmt.Print("正在监听", tar.listener_.Addr().String())
	if err != nil {
		// 处理监听错误
		return
	}
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
				fmt.Print("有连接了")
				IP := conn.RemoteAddr().String()
				id := tar.rcv_id_pool_.The()
				tar.signal_ <- IP
				tar.rcv_id_pool_.Pop()
				tar.io_map_.Register(IP)
				// 创建接收者
				var rcv receiver
				rcv.Init(id, conn,
					tar.io_map_.GetIn(IP), tar.io_map_.GetOut(IP))
				var node mListNode[*receiver]
				node.Init(&rcv)
				(*tar.rcv_map_)[id] = &node
				rcv.Start()
				// 添加接收者到链表
				tar.receivers_.Push_tail(&node)
				tar.current_rs++
			}
			select {
			case rls := <-tar.release_:
				IP := (*tar.rcv_map_)[rls].Get().GetIP()
				tar.io_map_.Erase(IP)
				tar.rcv_id_pool_.Push(rls)
				tar.receivers_.Delete((*tar.rcv_map_)[rls])
				delete((*tar.rcv_map_), rls)
				tar.current_rs--
			default:

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
