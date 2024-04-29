//
// Created by Messi on 24-4-29.
//

#ifndef JUDGE_SERVICE_HANDLERPROXY_H
#define JUDGE_SERVICE_HANDLERPROXY_H

#include "judge.pb.h"

class HandlerProxy {
public:
    static JudgeResponse handle(JudgeRequest& request);
};


#endif //JUDGE_SERVICE_HANDLERPROXY_H
