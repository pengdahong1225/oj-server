#!/bin/bash

current_path=$(pwd)

sudo docker run -v $(current_path):/root/builder -it -u root builder-core /bin/bash