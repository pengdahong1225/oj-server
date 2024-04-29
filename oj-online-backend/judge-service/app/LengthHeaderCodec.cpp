//
// Created by Messi on 24-4-29.
//

#include "LengthHeaderCodec.h"
#include "muduo/base/Logging.h"

int LengthHeaderCodec::decode(const muduo::string &data, JudgeRequest &msg) {
    int32_t length = std::atoi(data.substr(0, kHeaderLen).c_str());
    muduo::string body = data.substr(kHeaderLen);
    if (length != body.size()) {
        LOG_ERROR << "LengthHeaderCodec::onMessage"
                  << " -> "
                  << "the length of package is error";
        return -1;
    }
    if (!msg.ParseFromString(body)) {
        LOG_ERROR << "LengthHeaderCodec::onMessage"
                  << " -> "
                  << "decode from string error";
        return -1;
    }
    return 0;
}

int LengthHeaderCodec::encode(muduo::net::Buffer &buffer, const JudgeResponse &msg) {
    muduo::string body;
    if (!msg.SerializeToString(&body)) {
        LOG_ERROR << "LengthHeaderCodec::onMessage"
                  << " -> "
                  << "encode to string error";
        return -1;
    }
    buffer.append(body.data(), body.size());
    int32_t len = static_cast<int32_t>(body.size());
    int32_t be32 = muduo::net::sockets::hostToNetwork32(len);
    buffer.prepend(&be32, sizeof be32);

    return 0;
}
