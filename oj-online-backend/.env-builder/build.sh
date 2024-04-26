#!/bin/bash

if [ $# -ne 1 ]; then
    echo "Usage: $0 <input need build path>"
    exit
fi

path = $*

sudo docker run -v $(path):/root/builder -it -u root builder-core /bin/bash