//
// Created by Messi on 24-4-26.
//

#include "app/AppServer.h"
#include "muduo/base/Logging.h"
#include "muduo/net/EventLoop.h"
#include <unistd.h>

void init() {
    muduo::Logger::setLogLevel(muduo::Logger::DEBUG);
    muduo::Logger::setLogFileName("app.log"); // 设置日志文件名
    muduo::Logger::initialize(); // 初始化日志系统
}

int main(int argc, char *argv[]) {
    LOG_INFO << "pid = " << getpid();

    // server
    muduo::net::EventLoop loop;
    muduo::net::InetAddress listenAddr(9020);
    AppServer server(&loop, listenAddr);

    // 主线程loop，负责监听端口
    server.start();
    loop.loop();

    return 0;
}