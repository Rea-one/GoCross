package gocross

import (
	"log"
	"net"
	"time"
)

type Listener interface {
	Init(chan int, *iomap, *cimess)
	Start()
	Stop()
	serve()
	naction(net.Conn)
	daction()
}

type listener struct {
	rcv_id_pool_ mQueue[int]
	host_        string
	mn_host_     string
	receivers_   *mList[*receiver]
	// 本地表不使用指针保存
	rcv_map_   map[int]*mListNode[*receiver]
	io_map_    *iomap
	stop_      bool
	listener_  net.Listener
	release_   chan int
	signal_    chan string
	current_rs int
	max_rs_    int
	counter    int
	mns_num_   int
	mnConns_   mnConn
}

func (tar *listener) Init(signal chan string, iom *iomap, mess *cimess) {
	tar.current_rs = 0
	tar.max_rs_ = 4
	tar.mns_num_ = 4
	tar.host_ = mess.host_
	tar.mn_host_ = mess.mn_host_
	tar.signal_ = signal
	tar.io_map_ = iom
	tar.rcv_map_ = make(map[int]*mListNode[*receiver])
	tar.receivers_ = new(mList[*receiver])
	tar.rcv_id_pool_.Init()
	for i := range tar.max_rs_ {
		tar.rcv_id_pool_.Push(i)
	}

	lsn, err := net.Listen("tcp", tar.host_)
	if err != nil {
		log.Fatalf("初始化失败: %v\n", err)
	}
	tar.listener_ = lsn
	tar.mnConns_.Init(mess.pg_host_, tar.mns_num_)
	tar.release_ = make(chan int, 10)

	log.Printf("listener 初始化成功")
}

func (tar *listener) Start() {
	go tar.serve()
	log.Printf("listener 开始运行")
}

// 停止监听方法
func (tar *listener) Stop() {
	tar.stop_ = true
	if tar.listener_ != nil {
		tar.listener_.Close()
	}
}

func (tar *listener) serve() {
	for !tar.stop_ {
		tar.daction()

		if tar.current_rs >= tar.max_rs_ {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		conn, err := tar.listener_.Accept()
		if err != nil {
			log.Printf("连接接收失败: %v\n", err)
			continue
		}
		tar.counter++
		log.Printf("接收到新连接: %v\n", conn.RemoteAddr())
		tar.naction(conn)
	}
}

func (tar *listener) naction(conn net.Conn) {
	IP := conn.RemoteAddr().String()
	id := tar.rcv_id_pool_.The()
	tar.rcv_id_pool_.Pop()

	// 广播新连接信号
	tar.signal_ <- IP

	// 注册 IO 映射
	tar.io_map_.Register(IP)

	// 创建 receiver 实例
	rcv := &receiver{}
	rcv.Init(id, conn, &tar.mnConns_, tar.release_,
		tar.io_map_.GetIn(IP), tar.io_map_.GetOut(IP))

	// 创建链表节点
	node := new(mListNode[*receiver])
	node.Init(rcv)

	// 记录映射关系
	tar.rcv_map_[id] = node

	// 启动 receiver
	rcv.Start()

	// 加入链表
	tar.receivers_.Push_tail(node)
	tar.current_rs++
}

func (tar *listener) daction() {
	select {
	case rls := <-tar.release_:
		node := tar.rcv_map_[rls]
		if node == nil {
			return
		}
		rcv := node.Get()

		// 清理 IO 映射
		tar.io_map_.Erase(rcv.GetIP())

		// 回收 ID
		tar.rcv_id_pool_.Push(rls)

		// 从链表中删除
		tar.receivers_.Delete(node)

		// 删除 map 中的引用
		delete(tar.rcv_map_, rls)

		// 减少当前连接数
		tar.current_rs--

		log.Printf("回收 receiver ID: %d", rls)
	default:
		// 无回收请求则跳过
	}
}
