#!/bin/sh

region=`curl -s http://169.254.169.254/latest/dynamic/instance-identity/document|grep region|awk -F\" '{print $4}'`

join=""
servers=$(elb-discovery --load-balancer-name consul --region ${region} --count 0 --private-ip-only 2>&1)
ret_servers=$?

len=$(echo ${servers} | jq length 2>&1)
ret_ips=$?

[ $ret_servers == 0 ] && [ $ret_ips == 0 ] &&
    join_as="-join $(echo ${servers} | jq -r .[0] 2>&1)"

consul agent -data-dir /tmp/consul ${join}