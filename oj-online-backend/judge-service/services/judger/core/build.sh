#!/bin/bash

sudo docker run -v /root/server/oj-online-server/oj-online-backend/judge-service/judgecore:/root/builder -it -u root builder-core /bin/bash