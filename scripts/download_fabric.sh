#!/bin/bash

# Author: Cuiting Shi
# Email: ct.shi@hnair.com

REPOSITORY=hyperledger

IMAGES=("fabric-peer" "fabric-membersrvc")
TAG="x86_64-0.6.1-preview"

BASEIMAGE="fabric-baseimage"
BASETAGS=("x86_64-0.2.0" "x86_64-0.2.1" "x86_64-0.2.2")

echo ======================================
echo FABRIC IMAGES:
for i in ${IMAGES[@]}; do
  echo docker pull $REPOSITORY/$i:$TAG
  docker pull $REPOSITORY/$i:$TAG
  echo
done

echo ======================================
echo FABRIC CHAINCODE BASE IMAGES:
echo pull base image
for i in ${BASETAGS[@]}; do
  echo docker pull $REPOSITORY/$BASEIMAGE:$i
  docker pull $REPOSITORY/$BASEIMAGE:$i
  echo 
done
