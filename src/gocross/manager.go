package gocross

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager interface {
	Start()
	Stop()
	Wait()
	serve()
	Init(chan string, *Checker, *cimess)
}

type manager struct {
	// 与listener的通信
	signal_ chan string
	// 映射表，保证worker与receiver通过ID或IP唯一访问
	checker_ *Checker

	conn_pool_ *pgxpool.Pool
	stop_      bool

	wkr_num_ int
	workers  *mList[*worker]
	release_ chan int
	wkr_map_ map[int]*mListNode[*worker]
}

func (tar *manager) Start() {
	go tar.serve()
	log.Printf("manager 开始运行")
}

func (tar *manager) Stop() {

}

func (tar *manager) Wait() {

}

func (tar *manager) Init(signal chan string, checker *Checker, mess *cimess) {
	config, _ := pgxpool.ParseConfig(mess.PGString())

	config.MinConns = 2
	config.MaxConns = 4

	tar.checker_ = checker
	tar.signal_ = signal
	tar.conn_pool_, _ = pgxpool.NewWithConfig(context.Background(), config)
	tar.wkr_num_ = 4
	tar.release_ = make(chan int, tar.wkr_num_)
	tar.wkr_map_ = make(map[int]*mListNode[*worker])
	for count := range tar.wkr_num_ {
		conn, _ := tar.conn_pool_.Acquire(context.Background())
		if tar.workers == nil {
			tar.workers = new(mList[*worker])

		}
		tar.workers.Push_tail(&mListNode[*worker]{})
		tar.workers.Tail().Init(new(worker))
		tar.workers.Tail().data.Init(count, "default",
			tar.checker_, tar.release_, conn)
		tar.wkr_map_[count] = tar.workers.Tail()
	}
	log.Printf("manager 初始化成功")
}

func (tar *manager) serve() {
	for _ = range tar.wkr_num_ {
		tar.workers.head.Get().Start()
		tar.workers.Move_tail(tar.workers.head)
	}
	tickerFast := time.NewTicker(10 * time.Millisecond)  // 高频轮询间隔
	tickerSlow := time.NewTicker(100 * time.Millisecond) // 低频等待间隔
	defer tickerFast.Stop()
	defer tickerSlow.Stop()

	useFastPolling := false

	for !tar.stop_ {
		if useFastPolling {
			select {
			case IP := <-tar.signal_:
				now := tar.workers.Head()
				// 更换worker的工作对象，缺乏保障机制，但在一般情况下此时worker处于等待状态
				now.data.Change(IP)
				tar.workers.Move_tail(now)
				tickerFast.Reset(10 * time.Millisecond) // 保持高频
			case rls := <-tar.release_:
				tar.workers.Move_head(tar.wkr_map_[rls])
			case <-tickerFast.C:
				// 继续高频轮询
			}
		} else {
			select {
			case IP := <-tar.signal_:
				useFastPolling = true
				now := tar.workers.Head()
				now.data.Change(IP)
				tar.workers.Move_tail(now)
			case rls := <-tar.release_:
				tar.workers.Move_head(tar.wkr_map_[rls])
			case <-tickerSlow.C:
				// 低频等待，防止空转
			}
		}
	}
}
