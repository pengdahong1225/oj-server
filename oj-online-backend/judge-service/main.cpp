//
// Created by Messi on 24-4-26.
//

#include "app/AppServer.h"
#include "muduo/net/EventLoop.h"
#include "common/logger/rlog.h"
#include <unistd.h>
#include <libgen.h>

int main(int argc, char *argv[]) {
    // 初始化日志
    LOG_INIT("/app/log", "judge-service", DEBUG);
    LOG_INFO("master pid = %d", getpid());

    // server
    muduo::net::EventLoop loop;
    muduo::net::InetAddress listenAddr(9020);
    AppServer server(&loop, listenAddr);

    // 主线程loop，负责监听端口
    server.start();
    loop.loop();

    return 0;
}