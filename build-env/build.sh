if [ $# -ne 1 ]; then
    echo "Usage: $0 <input need build path>"
    exit
fi

sudo docker run -v $*:/root/builder -it -u root registry.cn-chengdu.aliyuncs.com/aliyun_pdh_space/builder:v1.0.0 /bin/bash
