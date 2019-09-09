#pkt-replay

#Build
1. cd to pkt-replay directory
2. 'make' to build the code

#Run

Run by specifying the pcap file name(example file, test.pcap can be used for sample analysis) for transmission and for response, specify the protocol(ip/tcp/udp) and number of packets the you would like to capture live. 

#sample command to run pkt-replay

```
./pkt_replay test.pcap ip 10
```

