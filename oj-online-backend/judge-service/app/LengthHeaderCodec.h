//
// Created by Messi on 24-4-29.
//

#ifndef JUDGE_SERVICE_LENGTHHEADERCODEC_H
#define JUDGE_SERVICE_LENGTHHEADERCODEC_H

#include "judge.pb.h"
#include "muduo/net/Buffer.h"

class LengthHeaderCodec {
    const static size_t kHeaderLen = sizeof(int32_t); // 4字节包头
public:
    static int decode(const muduo::string& data, JudgeRequest& msg);
    static int encode(muduo::string& data, const JudgeRequest& msg);
};


#endif //JUDGE_SERVICE_LENGTHHEADERCODEC_H
