# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation

#!/bin/bash

helpPrint()
{
   echo ""
   echo "Usage: $0 -t sanType -n subj1 -m subj2"
   echo -e "\t-t SAN type: could be IP or DNS"
   echo -e "\t-n subject alternative name 1: could be IP address or domain name for NEF."
   echo -e "\t-m subject alternative name 2: could be IP address or domain name for AF."
   exit 1 # Exit with help
}


while getopts "t:n:m:" opt
do
   case "$opt" in
      t ) sanType="$OPTARG" ;;
      n ) subj1="$OPTARG" ;;
      m ) subj2="$OPTARG" ;;
      ? ) helpPrint ;; # Print help
   esac
done


if [ -z "$subj1" ] || [ -z "$subj2" ] || [ -z "$sanType" ]
then
   echo "Some input parameters empty"
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
echo "$subj1"
echo "$subj2"

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
echo "subjectAltName = $sanType.1:$subj1,$sanType.2:$subj2" >> extfile.cnf
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
