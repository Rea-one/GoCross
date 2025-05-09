package gocross

import (
	"log"
	"net"
	"now/sqlmap"
	"time"
)

type Receiver interface {
	Init(int, net.Conn, chan int, chan sqlmap.Task, chan sqlmap.Task)
	Start()
	Stop()
	GetIP() string
}

type receiver struct {
	id_      int
	counter_ int
	conn_    net.Conn
	stop_    bool
	release_ chan int
	ipasser_ chan sqlmap.Task
	opasser_ chan sqlmap.Task
}

func (tar *receiver) Init(id int, conn net.Conn,
	release chan int, ip chan sqlmap.Task, op chan sqlmap.Task) {
	tar.conn_ = conn
	tar.id_ = id
	tar.stop_ = false
	tar.release_ = release
	tar.ipasser_ = ip
	tar.opasser_ = op
	tar.counter_ = 1
	log.Print("receiver 初始化完成")
}

func (tar *receiver) Start() {
	go tar.write()
	go tar.read()
}

func (tar *receiver) read() {
	buf := make([]byte, 64*1024) // 增大缓冲区以适应图像
	for !tar.stop_ {
		n, err := tar.conn_.Read(buf)
		if err != nil {
			tar.conn_.Close()
			break
		}
		data := make([]byte, n)
		copy(data, buf[:n])

		mess := string(data)
		switch {
		case mess == "nomore":
			// 处理终止信号
		case strings.HasPrefix(mess, "image "):
			// 提取 ImageID
			imageID := strings.TrimPrefix(mess, "image ")
			log.Printf("准备接收图片: %s", imageID)

			// 下次 Read 接收图像数据
			go func(id string) {
				n, err := tar.conn_.Read(buf)
				if err != nil {
					log.Printf("接收图片失败：%v", err)
					return
				}
				imgData := make([]byte, n)
				copy(imgData, buf[:n])

				// 上传到 MinIO
				url, err := uploadToMinio(id, imgData)
				if err != nil {
					log.Printf("上传到 MinIO 失败：%v", err)
					return
				}

				// 构造 Task
				new_task := sqlmap.Task{
					ID:       float64(tar.id_) / float64(tar.counter_),
					State:    "image uploaded",
					Ttype:    "pass",
					Deadline: time.Now().Add(time.Second * 30),
					ImageID:  url
				}
				tar.ipasser_ <- new_task
			}(imageID)

		default:
			// 处理普通查询
			new_task := sqlmap.Task{
				ID:       float64(tar.id_) / float64(tar.counter_),
				State:    "cross received",
				Ttype:    "pass",
				Deadline: time.Now().Add(time.Second * 10),
				Query:    mess,
			}
			tar.ipasser_ <- new_task
		}
	}
}

func (tar *receiver) write() {
	for !tar.stop_ {
		select {
		case task := <-tar.opasser_:
			switch task.GetState() {
			case "nomore":
				log.Printf("%v 号接收者关闭连接中", tar.id_)
				tar.conn_.Close()
				tar.release_ <- tar.id_
			case "request image":
				// 预告图片
				tar.conn_.Write([]byte("image " + task.ImageID))
				tar.conn_.Write(task.Image)
			default:
				tar.conn_.Write([]byte(task.GetState()))
			}
		default:
			time.Sleep(time.Millisecond * 300)
		}
	}
}

func (tar *receiver) Stop() {
	tar.stop_ = true
}

func (tar *receiver) GetIP() string {
	return tar.conn_.RemoteAddr().String()
}
