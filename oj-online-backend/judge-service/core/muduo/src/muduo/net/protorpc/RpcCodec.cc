// Copyright 2010, Shuo Chen.  All rights reserved.
// http://code.google.com/p/muduo/
//
// Use of this source code is governed by a BSD-style license
// that can be found in the License file.

// Author: Shuo Chen (chenshuo at chenshuo dot com)

#include "RpcCodec.h"

#include "../../base/Logging.h"
#include "../Endian.h"
#include "../TcpConnection.h"

#include "muduo/net/protorpc/rpc.pb.h"
#include "google-inl.h"

using namespace muduo;
using namespace muduo::net;

namespace
{
  int ProtobufVersionCheck()
  {
    GOOGLE_PROTOBUF_VERIFY_VERSION;
    return 0;
  }
  int dummy __attribute__ ((unused)) = ProtobufVersionCheck();
}

namespace muduo
{
namespace net
{
const char rpctag [] = "RPC0";
}
}
