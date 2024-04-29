//
// Created by Messi on 24-4-29.
//

#ifndef JUDGE_SERVICE_LENGTHHEADERCODEC_H
#define JUDGE_SERVICE_LENGTHHEADERCODEC_H

#include "muduo/net/Buffer.h"

class LengthHeaderCodec {
    const static size_t kHeaderLen = sizeof(int32_t); // 4字节包头
public:
    static int parse(const muduo::string& data,);
};


#endif //JUDGE_SERVICE_LENGTHHEADERCODEC_H
