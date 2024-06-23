//
// Created by Messi on 24-4-29.
//

#include "LengthHeaderCodec.h"
#include "common/logger/rlog.h"

int LengthHeaderCodec::decode(const muduo::string &data, SSJudgeRequest &msg) {
    int32_t length = std::atoi(data.substr(0, kHeaderLen).c_str());
    muduo::string body = data.substr(kHeaderLen);
    if (length != body.size()) {
        LOG_ERROR("LengthHeaderCodec::onMessage -> length[%d] not equal body size[%d]", length, body.size());
        return -1;
    }
    if (!msg.ParseFromString(body)) {
        LOG_ERROR("LengthHeaderCodec::onMessage -> decode from string error");
        return -1;
    }
    return 0;
}

int LengthHeaderCodec::encode(muduo::net::Buffer &buffer, const SSJudgeResponse &msg) {
    muduo::string body;
    if (!msg.SerializeToString(&body)) {
        LOG_ERROR("LengthHeaderCodec::onMessage -> encode from string error");
        return -1;
    }
    buffer.append(body.data(), body.size());
    int32_t len = static_cast<int32_t>(body.size());
    int32_t be32 = muduo::net::sockets::hostToNetwork32(len);
    buffer.prepend(&be32, sizeof be32);

    return 0;
}
