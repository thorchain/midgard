#!/bin/sh

# Tags created by passing version 1.2.3
# 1
# 1.2
# 1.2.3

echo " -t $1:$(echo $2 | awk -F '.' '{print $1}') -t $1:$(echo $2 | awk -F '.' '{print $1"."$2}') -t $1:$2 "
