#include <stdio.h>
#include <stdlib.h>
#include <pcap.h>
#include <string.h>
#define PCAP_SAVEFILE "./pcap_savefile"

int packets = 0;   /* running count of packets read in */
pcap_t *pl = NULL;
pcap_t *po = NULL;

int
usage(char *progname)
{
        printf("Usage: %s, [<savefile name>]\n", basename(progname));
        exit(7);
}

/********************************************************************
 * 
 *******************************************************************/
void pktPrint(const unsigned char *pkt, int pkt_len)
{
	int ii;
	
	printf("Pkt len: %d\n", pkt_len);
	for (ii=0; ii < pkt_len; ii++) {
		if (!(ii%16))
			printf("\n");
		printf(" %02x", pkt[ii]);
	}
	printf("\n");
}

/********************************************************************
 * 
 *******************************************************************/
pcap_t *pl_init(char *iface)
{
	pcap_t *pl;
	char pl_errbuf[PCAP_ERRBUF_SIZE];

	pl = pcap_open_live(iface, 1600, 0, 1, pl_errbuf);
	if (pl_errbuf[0] != '\0') {
		fprintf(stderr, "%s", pl_errbuf);
	}
	return pl;
}

/********************************************************************
 * 
 *******************************************************************/
int pl_send(pcap_t *pl, const unsigned char *pkt, int pkt_len)
{
	int rc;

	if (pcap_inject(pl, pkt, pkt_len) == -1) {
		pcap_perror(pl, 0);
		return -1;
	}

	printf("Sent success, len: %d\n", pkt_len);
	return pkt_len;
}


/********************************************************************
 * 
*******************************************************************/

int pl_recv(pcap_t *pl, const unsigned char *pkt, int pkt_len){

	//for debug purpose only print the packet 
	pktPrint(pkt, pkt_len);	

	if (isEPCPkt(pkt) && isEPCExpectedPkt(pkt)) {
		return 0;
	}

        return -1;
}

void pl_recv_handler(unsigned char *none, const struct pcap_pkthdr *pktHdr,
					const unsigned char *pkt)
{
	printf("%s: \n", __func__);
	static int pkt_count =0;
	int rc = 0;
	

	pkt_count++;

	rc = pl_recv(NULL, pkt, pktHdr->caplen);

	printf("DEBUG: pl_recv() return value: %d\n", rc);	

	if ((pkt_count == 100) || (rc == 0)) {
		printf("existing from pcap_loop()...\n");
		pcap_breakloop(pl);
	}
	
	return;
}
/*****************************************************************

 * *******************************************************************/
void
po_reader_handler(u_char *user, const struct pcap_pkthdr *hdr, const u_char *po_pkt)
{
        int offset = 26; /* 14 bytes for MAC header +
                          * 12 byte offset into IP header for IP addresses
                          *                           */

        if (hdr->caplen < 30) {
                /* captured data is not long enough to extract IP address */
                fprintf(stderr,
                        "Error: not enough captured packet data present to extract IP addresses.\n");
                return;
        }
                
	//Check this packet is from from enb
	if (!isEnbPkt(po_pkt)) {
		printf("DEBUG: drop this PO packet\n");
	}

	//send pkt to EPC
	pl_send(pl_po_pkt, hdr->caplen); 

	//wait for EPC response
	pl_recv_start();

	return;
}

void pl_recv_start()
{
	printf("DEBUG: Enter to pcap_live_recv() loop\n");
	if (pcap_loop(pl, 50, pl_recv_handler, NULL) < 0) {
		printf("DEBUG: pcap_loop exited... : %s\n", pcap_geterr(pl));
	}		
	printf("DEBUG: Exit from pcap_live_recv() loop...\n");
}

int
main(int argc, char **argv)
{
        char filename[80];       /* name of savefile to read packet data from */
        char errbuf[PCAP_ERRBUF_SIZE];  /* buffer to hold error text */
        char prestr[80];         /* prefix string for errors from pcap_perror */

	if (argc < 1)
                usage(argv[0]);
        if (argc >= 2) 
                strcpy(filename,argv[1]);
        else
                strcpy(filename, PCAP_SAVEFILE);
	
	pl = pl_init("enp24s0f0");
	//pl = pl_open_live("enp24s0f0", BUFSIZ, 0, -1, errbuf); 
	if (pl == NULL) {
		printf ("pl_init() failed\n");
		return 0;
	}

        if (!(po = pcap_open_offline(filename, errbuf))) {
                fprintf(stderr,
                        "Error in opening savefile, %s, for reading: %s\n",
                        filename, errbuf);
                exit(2);
        }

        if (pcap_dispatch(po, 0, &po_reader_handler, (char *)0) < 0) {
                sprintf(prestr,"Error reading packets");
                pcap_perror(po,prestr);
                exit(4);
        }		

        printf("\nTotal no. of Packets sent : %d\n", packets);


	
        //pcap_close(po);
        pcap_close(pl);
	printf("capture complete.\n");

}
