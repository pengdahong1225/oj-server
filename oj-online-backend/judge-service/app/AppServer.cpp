//
// Created by Messi on 24-4-28.
//

#include "AppServer.h"
#include "muduo/base/Logging.h"

using std::placeholders::_1;
using std::placeholders::_2;
using std::placeholders::_3;

AppServer::AppServer(muduo::net::EventLoop *loop, const muduo::net::InetAddress &listenAddr) : server_(loop, listenAddr,
                                                                                                       "judge-service") {
    server_.setConnectionCallback(
            std::bind(&AppServer::onConnection, this, _1));
    server_.setMessageCallback(
            std::bind(&AppServer::onMessage, this, _1, _2, _3));
}

void AppServer::start() {
    server_.start();
}

void AppServer::onConnection(const muduo::net::TcpConnectionPtr &conn) {
    LOG_INFO << "EchoServer - " << conn->peerAddress().toIpPort() << " -> "
             << conn->localAddress().toIpPort() << " is "
             << (conn->connected() ? "UP" : "DOWN");
}

void AppServer::onMessage(const muduo::net::TcpConnectionPtr &conn, muduo::net::Buffer *buf, muduo::Timestamp time) {
    muduo::string msg(buf->retrieveAllAsString());
    LOG_INFO << conn->name() << " echo " << msg.size() << " bytes, "
             << "data received at " << time.toString();
    conn->send(msg);
}
