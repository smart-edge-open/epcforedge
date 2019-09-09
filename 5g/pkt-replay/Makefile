CC		= gcc
LDFLAGS 	= -lpcap
SOURCES		= pkt_replay.c
OBJECTS		= $(SOURCES:.c=.o)
TARGET		= pkt_replay

all: $(SOURCES) $(TARGET)

$(TARGET): $(OBJECTS)
	$(CC) $(LDFLAGS) $(OBJECTS) -O $@

.c.o:
	$(CC) $< -o $@

clean:
	rm -rf $(OBJECTS) $(TARGET)
