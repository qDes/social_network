#!/bin/bash
#title          : Launch 2 MYSQLnodes2
#description            : This script creates 2 MySQL containers and launches them assuming the data /log/backup/conf.d are present in /opt2/mysql/<node_prefix><nodenumber> folders.
#author         : Avinash Barnwal
#date           : 22092016
#version        : 0.1
#usage          : bash Launch2MySqlnodes
#=============================================================================

DB_NAME=mydata
ROOT_PASS=roo235t
MYSQL_IMAGE='mysql:latest'
NODE_PREFIX=mysql

docker run --name ${NODE_PREFIX}1 \
       -e MYSQL_ROOT_PASSWORD=$ROOT_PASS \
       -e MYSQL_DATABASE=$DB_NAME -dit \
       -v ./mysql/${NODE_PREFIX}1/conf.d:/etc/mysql/mysql.conf.d/ \
       -v ./mysql/${NODE_PREFIX}1/data:/var/lib/mysql \
       -v ./mysql/${NODE_PREFIX}1/log:/var/log/mysql \
       -v ./mysql/${NODE_PREFIX}1/backup:/backup \
       -p 3306 \
       -h ${NODE_PREFIX}1 $MYSQL_IMAGE

NODE1_PORT=$(docker inspect --format='{{(index (index .NetworkSettings.Ports "3306/tcp") 0).HostPort}}' ${NODE_PREFIX}1)
 # https://docs.docker.com/engine/reference/commandline/inspect/

docker run --name ${NODE_PREFIX}2 \
       -e MYSQL_ROOT_PASSWORD=$ROOT_PASS \
       -e MYSQL_DATABASE=$DB_NAME -dit \
       --link ${NODE_PREFIX}1:${NODE_PREFIX}1cl \
       -v ./mysql/${NODE_PREFIX}2/conf.d:/etc/mysql/mysql.conf.d/ \
       -v ./mysql/${NODE_PREFIX}2/data:/var/lib/mysql \
       -v ./mysql/${NODE_PREFIX}2/log:/var/log/mysql \
       -v ./mysql/${NODE_PREFIX}2/backup:/backup \
       -p 3306 \
       -h ${NODE_PREFIX}2 $MYSQL_IMAGE

NODE2_IP=$(docker inspect --format '{{ .NetworkSettings.IPAddress }}' ${NODE_PREFIX}2)
# This would add the second node's IP in the Host file of mysql first node.
docker exec -i ${NODE_PREFIX}1 sh -c 'echo '$NODE2_IP ${NODE_PREFIX}2 ${NODE_PREFIX}2' >> /etc/hosts';