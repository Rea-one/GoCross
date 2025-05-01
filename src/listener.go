package main

import (
	"net"
)

type Listener interface {
	Start()
}

type listener struct {
	address   string
	receivers []receiver
	listener  net.Listener
}

func (tar *listener) Start() {

	for {
		select {
		case net.Listen("tcp", tar.address):

		}
	}
}

// type listener struct {
// 	address   string
// 	receivers []receiver
// 	quit      chan struct{}  // 新增退出通知通道
// 	listener  net.Listener   // 缓存监听器实例
// }

// func (tar *listener) Start() {
// 	// 初始化监听器
// 	var err error
// 	tar.listener, err = net.Listen("tcp", tar.address)
// 	if err != nil {
// 		// 处理监听错误
// 		return
// 	}
// 	defer tar.listener.Close()

// 	// 初始化退出通道
// 	tar.quit = make(chan struct{})

// 	// 启动连接处理协程
// 	go func() {
// 		for {
// 			conn, err := tar.listener.Accept()
// 			if err != nil {
// 				select {
// 				case <-tar.quit:
// 					// 正常退出
// 					return
// 				default:
// 					// 处理异常
// 				}
// 			}
// 			// 处理连接
// 			for _, rcv := range tar.receivers {
// 				go rcv.Handle(conn)  // 异步处理
// 			}
// 		}
// 	}()
// }

// // 停止监听方法
// func (tar *listener) Stop() {
// 	close(tar.quit)
// 	if tar.listener != nil {
// 		tar.listener.Close()
// 	}
// }
