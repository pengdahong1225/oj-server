//
// Created by Messi on 24-4-28.
//

#include "AppServer.h"
#include "LengthHeaderCodec.h"
#include "judge.pb.h"
#include "handler/HandlerProxy.h"
#include "common/logger/rlog.h"

using std::placeholders::_1;
using std::placeholders::_2;
using std::placeholders::_3;

AppServer::AppServer(muduo::net::EventLoop *loop, const muduo::net::InetAddress &listenAddr)
        : server_(loop, listenAddr,
                  "judge-service") {
    // 主要传递给主线程loop的Acceptor
    server_.setConnectionCallback(
            std::bind(&AppServer::onConnection, this, _1));
    // 传递给主线程loop的Acceptor 和 新连接
    server_.setMessageCallback(
            std::bind(&AppServer::onMessage, this, _1, _2, _3));
}

void AppServer::start() {
    server_.start();
}

void AppServer::onConnection(const muduo::net::TcpConnectionPtr &conn) {
    LOG_INFO("[%s] -> [%s] is %s", conn->peerAddress().toIpPort().c_str(), conn->localAddress().toIpPort().c_str(),
             conn->connected() ? "UP" : "DOWN");
}

/*
 * 将消息交给HandlerProxy处理
 * 每个子线程loop收到消息都会回调这个函数，将消息传递到应用层来
 * 虽然都是同一个函数，但在每个线程空间都有独立调用栈
 */

/**
 * 循环读取包头，直到读取完一个完成的包头 -> 再去解码包头
 * 根据包头提供的长度，循环读取body，直到读取完一个完整的body -> 再去解码body
 */
void AppServer::onMessage(const muduo::net::TcpConnectionPtr &conn, muduo::net::Buffer *buf, muduo::Timestamp time) {
//    muduo::string data = "";
//    while (buf->readableBytes() > 0)
//        data.append(buf->retrieveAllAsString());
//    buf->retrieveAll();
    muduo::string data(buf->retrieveAllAsString());
    LOG_INFO("data received: %s", data.c_str());
    // 协议解析
    SSJudgeRequest request;
    if (LengthHeaderCodec::decode(data, request) != 0) {
        return;
    }

    // 处理
    SSJudgeResponse response = HandlerProxy::handle(request);

    LOG_INFO("response: %s", response.DebugString().c_str());

    // 返回
    muduo::net::Buffer buffer;
    if (LengthHeaderCodec::encode(buffer, response) == 0) {
        conn->send(&buffer);
    }
}
