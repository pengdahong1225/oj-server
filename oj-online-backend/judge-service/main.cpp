//
// Created by Messi on 24-4-26.
//

#include "app/AppServer.h"
#include "muduo/base/Logging.h"
#include "muduo/base/LogFile.h"
#include "muduo/net/EventLoop.h"
#include <unistd.h>
#include <libgen.h>

std::unique_ptr<muduo::LogFile> g_logFile;

void outputFunc(const char *msg, int len) {
//    fwrite(msg, 1, len, stdout);
    g_logFile->append(msg, len);
}

void flushFunc() {
//    fflush(stdout);
    g_logFile->flush();
}

int main(int argc, char *argv[]) {
    // 初始化日志
    char name[256] = {'\0'};
    strncpy(name, argv[0], sizeof name - 1);
    g_logFile.reset(new muduo::LogFile(::basename(name), 200 * 1000));
    muduo::Logger::setOutput(outputFunc);
    muduo::Logger::setFlush(flushFunc);
    muduo::Logger::setLogLevel(muduo::Logger::DEBUG);

    LOG_INFO << "master pid = " << getpid();

    // server
    muduo::net::EventLoop loop;
    muduo::net::InetAddress listenAddr(9020);
    AppServer server(&loop, listenAddr);

    // 主线程loop，负责监听端口
    server.start();
    loop.loop();

    return 0;
}