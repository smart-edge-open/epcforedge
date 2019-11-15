# Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/bash

helpPrint()
{
   echo ""
   echo "Usage: $0 -m method -i resourceId"
   echo -e "\t-m Simulated HTTP Method. It could be POST, GET, DEL,PATCH or DELDNN "
   echo -e "\t-i ResourceId for the method to operate. It should be afId"
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
   "POST") curl -vvv -X POST -i "Content-Type: application/json" --data @./json/POST001.json http://localhost:8080/oam/v1/af/services;;
   "GET") curl -vvv http://localhost:8080/oam/v1/af/services/$resourceId;;
   "PATCH") curl -vvv -X PATCH -i "Content-Type: application/json" --data @./json/PATCH001.json http://localhost:8080/oam/v1/af/services/$resourceId;;
   "DEL") curl -vvv -X DELETE http://localhost:8080/oam/v1/af/services/$resourceId;;
   "DELDNN") curl -vvv -X DELETE http://localhost:8080/oam/v1/af/services/$resourceId/locationServices/by_dnai;;
   *) echo "Wrong method!";;
esac
