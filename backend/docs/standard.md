# grpc服务
返回的err统一使用grpc提供的status.Errorf

数据不存在要和err区分开

# gateway
返回的err统一使用proto中定义的Error