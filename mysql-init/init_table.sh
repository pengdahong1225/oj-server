#!/bin/bash

MYSQL=/usr/bin/mysql

db_host=127.0.0.1
db_port=3306
db_user=root
db_pwd=1225pdhGkl
db_name=gvb_server

SOURCE_FILE=db.sql

# 建库
$MYSQL -h$db_host -P$db_port -u$db_user -p$db_pwd -e "create database if not exists $db_name default character set = utf8mb4;"

# 建表
$MYSQL -h$db_host -P$db_port -u$db_user -p$db_pwd -e "source $SOURCE_FILE;"

if [ $? -eq 0 ]
then
	echo "table_create done."
else
	echo "table_create failed."
fi
