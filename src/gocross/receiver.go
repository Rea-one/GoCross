package gocross

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net"
	"now/sqlmap"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type Receiver interface {
	Init(int, net.Conn, chan int, *mnConn, chan sqlmap.Task, chan sqlmap.Task)
	Start()
	Stop()
	GetIP() string
}

type receiver struct {
	id_           int
	counter_      int
	conn_         net.Conn
	mnConn_       *mnConn
	stop_         bool
	release_      chan int
	ipasser_      chan sqlmap.Task
	opasser_      chan sqlmap.Task
	feedback_     feedback
	lastHeartbeat time.Time // 记录最后一次收到心跳的时间
}

func (tar *receiver) Init(id int, conn net.Conn, mnc *mnConn,
	release chan int, ip chan sqlmap.Task, op chan sqlmap.Task) {
	// socket 连接
	tar.conn_ = conn
	// 未使用的minio的连接
	tar.mnConn_ = mnc
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
	log.Printf("receiver已启动 id: %d", tar.id_)
}

// 这里只有少数情况是在这处理，多数在worker中处理
func (tar *receiver) read() {
	buf := make([]byte, 64*1024) // 增大缓冲区以适应图像
	for !tar.stop_ {
		// 设置读取超时
		tar.conn_.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := tar.conn_.Read(buf)
		if err != nil {
			log.Printf("读取数据失败: %v", err)
			tar.Stop() // 主动调用 Stop() 释放资源
			break
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		mess := string(data)
		switch {
		case mess == "nomore":
			// 处理终止信号
			tar.Stop()
		case mess == "pong":
			// 收到心跳响应，更新 lastHeartbeat
			tar.lastHeartbeat = time.Now()
		case strings.HasPrefix(mess, "image "):
			// 提取 ImageID
			imageID := strings.TrimPrefix(mess, "image ")
			log.Printf("准备接收图片: %s", imageID)
			n, err := tar.conn_.Read(buf)
			if err != nil {
				log.Printf("读取图片失败: %v", err)
				continue
			}
			imageData := make([]byte, n)
			copy(imageData, buf[:n])
			log.Printf("接收图片成功: %s", imageID)
			now := tar.mnConn_.Get()
			defer tar.mnConn_.ReGet(now)
			// 上传图片到 MinIO
			bucketName := "images"                    // 指定存储桶名称
			objectName := imageID                     // 使用 imageID 作为对象名
			contentType := "application/octet-stream" // 可根据实际情况调整

			policy := `{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": "*",
						"Action": "s3:GetObject",
						"Resource": "arn:aws:s3:::images/*"
					}
				]
			}`
			// 确保 bucket 存在
			err = now.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
			if err != nil {
				// 检查 bucket 是否已存在
				exists, errBucketExists := now.BucketExists(context.Background(), bucketName)
				if !exists || errBucketExists != nil {
					log.Printf("无法创建或访问存储桶: %v", err)
					continue
				}
			}
			now.SetBucketPolicy(context.Background(), bucketName, policy)

			// 上传文件
			_, err = now.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(imageData), int64(len(imageData)), minio.PutObjectOptions{ContentType: contentType})
			if err != nil {
				log.Printf("上传图片到 MinIO 失败: %v", err)
			} else {
				log.Printf("图片 %s 已上传至 MinIO", imageID)
			}
			tar.feedback_ = feedback{
				Image:     "https://" + tar.mnConn_.host_ + "/images/" + imageID,
				Message:   "图片已上传至 MinIO",
				Receiver:  "",
				Sender:    "",
				State:     "image",
				Timestamp: "",
			}
			feedbackBytes, err := json.Marshal(tar.feedback_)
			if err != nil {
				log.Printf("序列化 feedback 失败: %v", err)
				return
			}
			tar.conn_.Write(feedbackBytes)
		default:
			// 处理普通查询
			new_task := sqlmap.Task{
				Deadline: time.Now().Add(time.Second * 10),
				Message:  mess,
			}
			tar.ipasser_ <- new_task
		}

		// 检查心跳是否超时
		if time.Since(tar.lastHeartbeat) > 15*time.Second {
			log.Printf("心跳超时，关闭 receiver ID: %d", tar.id_)
			tar.Stop()
		}
	}
}

// 数据返回采用json格式，与接收格式不同
func (tar *receiver) write() {
	// 启动心跳检测
	go tar.waitBeat()
	for !tar.stop_ {
		select {
		case task := <-tar.opasser_:
			// 根据类型处理
			switch task.Ttype {
			case "pass":
				tar.write_pass(task)
			case "add friend":
				tar.write_addFriend(task)
			case "response add friend":
				tar.write_resAddFriend(task)
			default:
				tar.write_single(task)
			}
		default:
			time.Sleep(time.Millisecond * 300)
		}
	}
}

func (tar *receiver) write_single(task sqlmap.Task) {
	feedbackBytes, err := json.Marshal(feedback{
		At:        task.At,
		Sender:    task.Sender,
		Receiver:  task.Receiver,
		State:     task.GetState(),
		Timestamp: task.TimeStamp,
		Message:   task.Message,
		Image:     task.ImageURL,
	})
	if err != nil {
		log.Printf("序列化 feedback 失败: %v", err)
		return
	}
	tar.conn_.Write(feedbackBytes)
}

func (tar *receiver) write_pass(task sqlmap.Task) {
	feedbackBytes, err := json.Marshal(feedback{
		At:        task.At,
		Sender:    task.Sender,
		Receiver:  task.Receiver,
		State:     task.GetState(),
		Timestamp: task.TimeStamp,
		Message:   task.Message,
		Image:     task.ImageURL,
	})
	if err != nil {
		log.Printf("序列化 feedback 失败: %v", err)
		return
	}
	tar.conn_.Write(feedbackBytes)
}
func (tar *receiver) write_addFriend(task sqlmap.Task) {
	feedbackBytes, err := json.Marshal(feedback{
		At:        task.At,
		Sender:    task.Sender,
		Receiver:  task.Receiver,
		State:     task.GetState(),
		Timestamp: task.TimeStamp,
	})
	if err != nil {
		log.Printf("序列化 feedback 失败: %v", err)
		return
	}
	tar.conn_.Write(feedbackBytes)
}

func (tar *receiver) write_resAddFriend(task sqlmap.Task) {
	feedbackBytes, err := json.Marshal(feedback{
		At:        task.ImageID,
		Sender:    task.Sender,
		Receiver:  task.Receiver,
		State:     task.GetState(),
		Timestamp: task.TimeStamp,
	})
	if err != nil {
		log.Printf("序列化 feedback 失败: %v", err)
		return
	}
	tar.conn_.Write(feedbackBytes)
}
func (tar *receiver) isStopped() bool {
	return tar.stop_
}

func (tar *receiver) waitBeat() {
	feedbackBytes, err := json.Marshal(feedback{
		State: "ping",
	})
	for !tar.isStopped() {
		time.Sleep(5 * time.Second) // 每5秒发送一次心跳
		if tar.isStopped() {
			return
		}

		if err != nil {
			log.Printf("序列化 feedback 失败: %v", err)
			return
		}
		tar.conn_.Write(feedbackBytes)
	}
}

func (tar *receiver) Stop() {
	tar.conn_.Close()
	tar.release_ <- tar.id_
	tar.stop_ = true
}

func (tar *receiver) GetIP() string {
	return tar.conn_.RemoteAddr().String()
}
