
# 快速入手

main.go是程序入口
gocross是核心模块
gocross/models未使用
sql_map是辅助模块

## gocross 模块

sever作为统一接口，统筹gocross的所有模块
listener作为监听器，负责监听端口，接收请求，转发请求
manager作为数据库交互模块调度器

receiver作为接收模块
receiver从socket中接收数据，并转发到worker
worker作为处理模块
利用sql_map解析出sql语句，并通过worker与数据库交互
