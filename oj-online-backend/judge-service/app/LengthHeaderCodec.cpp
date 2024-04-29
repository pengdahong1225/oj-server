//
// Created by Messi on 24-4-29.
//

#include "LengthHeaderCodec.h"
#include "muduo/base/Logging.h"

int LengthHeaderCodec::parse(const std::string &data) {
    int32_t length = std::atoi(data.substr(0, kHeaderLen).c_str());
    muduo::string body = data.substr(kHeaderLen);
    if (length != body.size()) {
        LOG_ERROR << "LengthHeaderCodec::onMessage"
                  << " -> "
                  << "the length of package is error";
        return -1;
    }


    return 0;
}
