# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation

#!/bin/bash

helpPrint()
{
   echo ""
   echo "Usage: $0 -m method -i resourceId"
   echo -e "\t-m Simulated HTTP Method. It could be POST, GET, DEL,PATCH  "
   echo -e "\t-i ResourceId for the method to operate. It should be afServiceId"
   exit 1 # Exit with help
}


while getopts "m:i:" opt
do
   case "$opt" in
      m ) method="$OPTARG" ;;
      i ) resourceId="$OPTARG" ;;
      ? ) helpPrint ;; # Print help
   esac
done


if [ -z "$method" ] || [ -z "$resourceId" ]
then
   echo "Some input parameters empty";
   helpPrint
fi

echo "Running with input parameters:"
echo "$method"
echo "$resourceId"

case $method in 
   "POST") curl -vvv -X POST -i "Content-Type: application/json" --data @./json/POST001.json http://localhost:8080/ngcoam/v1/af/services;;
   "GET") curl -vvv http://localhost:8080/ngcoam/v1/af/services/$resourceId;;
   "PATCH") curl -vvv -X PATCH -i "Content-Type: application/json" --data @./json/PATCH001.json http://localhost:8080/ngcoam/v1/af/services/$resourceId;;
   "DEL") curl -vvv -X DELETE http://localhost:8080/ngcoam/v1/af/services/$resourceId;;
   *) echo "Wrong method!";;
esac
