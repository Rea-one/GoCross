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
	id_       int
	counter_  int
	conn_     net.Conn
	mnConn_   *mnConn
	stop_     bool
	release_  chan int
	ipasser_  chan sqlmap.Task
	opasser_  chan sqlmap.Task
	feedback_ feedback
}

func (tar *receiver) Init(id int, conn net.Conn, mnc *mnConn,
	release chan int, ip chan sqlmap.Task, op chan sqlmap.Task) {
	tar.conn_ = conn
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
				State:     "success",
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
			// switch task.GetState() {
			// case "nomore":
			// 	log.Printf("%v 号接收者关闭连接中", tar.id_)
			// 	tar.Stop()
			// case "request image":
			// 	// 预告图片
			// 	tar.conn_.Write([]byte("image " + task.ImageID))
			// 	tar.conn_.Write(task.Image)
			// default:
			// 	tar.conn_.Write([]byte(task.GetState()))
			// }
			tar.feedback_ = feedback{
				Image:     task.ImageID,
				Message:   task.Message,
				Receiver:  task.Sender,
				Sender:    task.Receiver,
				State:     task.GetState(),
				Timestamp: task.TimeStamp,
			}
			feedbackBytes, err := json.Marshal(tar.feedback_)
			if err != nil {
				log.Printf("序列化 feedback 失败: %v", err)
				return
			}
			tar.conn_.Write(feedbackBytes)
		default:
			time.Sleep(time.Millisecond * 300)
		}
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
