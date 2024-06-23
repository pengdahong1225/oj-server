//
// Created by Messi on 24-4-29.
//

#include "LengthHeaderCodec.h"
#include "common/logger/rlog.h"
#include <arpa/inet.h>  // for ntohl

int LengthHeaderCodec::decode(const muduo::string &data, SSJudgeRequest& msg) {
    // 检查数据是否足够包含长度头
    if (data.size() < sizeof(uint32_t)) {
        LOG_ERROR("LengthHeaderCodec::onMessage -> length[%d] not a complete packet header size[%d]", data.size(), sizeof(uint32_t));
        return -1;  // 不满足一个完整的包头
    }

    // 读取并解析长度头
    uint32_t length = ntohl(*reinterpret_cast<const uint32_t*>(data.data()));

    // 检查数据是否足够包含消息体
    if (data.size() != sizeof(uint32_t) + length) {
        LOG_ERROR("LengthHeaderCodec::onMessage -> length[%d] not equal body size[%d]", length, data.size() - sizeof(uint32_t));
        return -1;  // 不满足一个完整的包体
    }

    // 解析消息
    if (!msg.ParseFromArray(data.data() + sizeof(uint32_t), length)) {
        LOG_ERROR("LengthHeaderCodec::onMessage -> decode from string error");
        return -1;  // 解析失败
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