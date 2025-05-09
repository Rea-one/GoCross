package gocross

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MNConn interface {
	Init(host string, num int)
	Get() *minio.Client
	ReGet(*minio.Client)
	creat_conn() *minio.Client
}

type mnConn struct {
	host_     string
	conn_num_ int
	conn_     mQueue[*minio.Client]
}

func (tar *mnConn) creat_conn() *minio.Client {
	accessKey, _ := generateRandomKey(16) // 生成一个16字节的Access Key
	secretKey, _ := generateRandomKey(32) // 生成一个32字节的Secret Key
	conn, err := minio.New(tar.host_, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("mnConn 创建连接失败: %v", err)
	}
	return conn
}

func (tar *mnConn) Init(host string, num int) {
	tar.host_ = host
	tar.conn_num_ = num
	for tar.conn_.Size() < tar.conn_num_ {
		tar.conn_.Push(tar.creat_conn())
	}
}

func (tar *mnConn) Get() *minio.Client {
	if tar.conn_.Size() < tar.conn_num_ {
		tar.conn_.Push(tar.creat_conn())
	}
	defer tar.conn_.Pop()
	return tar.conn_.The()
}

func (tar *mnConn) ReGet(conn *minio.Client) {
	tar.conn_.Push(conn)
}

func generateRandomKey(length int) (string, error) {
	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
