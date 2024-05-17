//
// Created by Messi on 24-4-29.
//

#ifndef JUDGE_SERVICE_LENGTHHEADERCODEC_H
#define JUDGE_SERVICE_LENGTHHEADERCODEC_H

#include "proto/judge.pb.h"
#include "muduo/net/Buffer.h"

class LengthHeaderCodec {
    const static size_t kHeaderLen = sizeof(int32_t); // 4字节包头
public:
    static int decode(const muduo::string& data, SSJudgeRequest& msg);
    static int encode(muduo::net::Buffer& buffer, const SSJudgeResponse& msg);
};


#endif //JUDGE_SERVICE_LENGTHHEADERCODEC_H
