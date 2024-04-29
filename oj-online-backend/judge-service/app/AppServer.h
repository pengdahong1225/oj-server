//
// Created by Messi on 24-4-28.
//

#ifndef JUDGE_SERVICE_APPSERVER_H
#define JUDGE_SERVICE_APPSERVER_H

#include "muduo/net/TcpServer.h"

class AppServer {
public:
    AppServer(muduo::net::EventLoop *loop,
              const muduo::net::InetAddress &listenAddr);

    void start();  // calls server_.start();

private:
    void onConnection(const muduo::net::TcpConnectionPtr &conn);

    void onMessage(const muduo::net::TcpConnectionPtr &conn,
                   muduo::net::Buffer *buf,
                   muduo::Timestamp time);

private:
    muduo::net::TcpServer server_;
};


#endif //JUDGE_SERVICE_APPSERVER_H
