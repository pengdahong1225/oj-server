#include "../Inspector.h"
#include "../../EventLoop.h"
#include "../../EventLoopThread.h"

using namespace muduo;
using namespace muduo::net;

int main()
{
  EventLoop loop;
  EventLoopThread t;
  Inspector ins(t.startLoop(), InetAddress(12345), "test");
  loop.loop();
}

