# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation

#!/bin/bash

helpPrint()
{
   echo ""
   echo "Usage: $0 -t sanType -h {subj | *}"
   echo -e "\t-t SAN type: could be IP or DNS"
   echo -e "\t-h subject alternative names list: could be one or more IP addressess or domain names separated by space"
   exit 1 # Exit with help
}

Hostname_count=0
Hostname_flag=0
subjstr=""

while [ "$1" != "" ]; 
do
   case $1 in
      -t )
         if [ $Hostname_flag == 1 ]
         then
            echo "Incorrect Command Sequence"
            helpPrint
         fi
         shift 
         if [ "$1" != "" ]
         then
            sanType="$1"
         else
            echo "Missing argument for -t option"
            helpPrint
         fi
         ;;
      -h )
         if [ "$sanType" == "" ]
         then
            echo "Incorrect Command Sequence. Got -h option before -t option"
            helpPrint
         fi
         shift
         if [ "$1" != "" ]
         then
            subj1="$1"
            Hostname_flag=$((Hostname_flag+1))
            Hostname_count=$((Hostname_count+1))
            subjstr=$sanType"."$Hostname_count":""$1"
         else
            echo "Missing argument for -h option"
            helpPrint
         fi
         ;;
      ? ) helpPrint # Print help
         ;;
      * )
         if [ $Hostname_flag == 1 ]
         then
            Hostname_count=$((Hostname_count+1))
            subjstr+=","$sanType"."$Hostname_count":""$1"
         else
            echo "Incorrect Input"
            helpPrint
         fi
         ;;
  esac
  shift
done

if [ -z "$subj1" ] || [ -z "$subjstr" ] || [ -z "$sanType" ]
then
   echo "One of the input parameters missing"
   helpPrint
fi

if [ "$sanType" == IP ] || [ "$sanType" == DNS ]
then 
   echo "Input OK"
else
   echo "Wrong sanType"
   helpPrint
   exit 1
fi

echo "Running with input parameters:"
echo "$sanType"
#echo "$subj1"
echo "$subjstr"

ROOT_CA_NAME=OpenNESS-5G-Root

echo "Generating RootCA Key and Cert:"
openssl ecparam -genkey -name secp384r1 -out "root-ca-key.pem"
if (($?))
then 
   echo "RootCA key generation failed"
   exit 1
fi

openssl req -key "root-ca-key.pem" -new -x509 -days 90 -subj "/CN=$ROOT_CA_NAME" -out "root-ca-cert.pem" 
if (($?))
then 
   echo "RootCA cert generation failed"
   exit 1
fi

echo "Generating Server Key and Cert:"
openssl ecparam -genkey -name secp384r1 -out "server-key.pem"
if (($?))
then 
   echo "Servier key generation failed"
   exit 1
fi

openssl req -new -key "server-key.pem" -out "server-request.csr" -subj "/CN=$subj1"
if (($?))
then 
   echo "Servier CSR generation failed"
   exit 1
fi
rm -f extfile.cnf
#echo "subjectAltName = $sanType.1:$subj1,$sanType.2:$subj2,$sanType.3:$subj3,$sanType.4:$subj4" >> extfile.cnf
echo "subjectAltName = $subjstr" >> extfile.cnf
openssl x509 -req -extfile extfile.cnf -in "server-request.csr" -CA "root-ca-cert.pem" -CAkey "root-ca-key.pem" -days 90 -out "server-cert.pem" -CAcreateserial
if (($?))
then 
   echo "Servier cert generation failed"
   exit 1
fi

echo "Print CA Cert Pem:"
openssl x509 -in root-ca-cert.pem -text -noout
echo "Print Server Cert Pem:"
openssl x509 -in server-cert.pem -text -noout
echo "Successfully completed"
