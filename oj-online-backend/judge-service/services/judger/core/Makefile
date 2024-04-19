CC  = gcc
GCCVER := $(shell $(CC) -dumpversion | awk -F. '{ print $$1"."$$2}' )

CFLAGS = -g -Wall -Werror -O3 -std=c99 -lseccomp # libseccomp是一个用于 Linux 系统的 syscall 过滤器库

INC	= -I./src -I./src/rules
SRCS = $(wildcard ./src/*.c) $(wildcard ./src/rules/*.c)

COMPILE_LIB_HOME = .
DYNAMIC_NAME = libjudger.so
STATIC_NAME = libjudger.a
DYNAMIC_LIB	= $(COMPILE_LIB_HOME)/$(DYNAMIC_NAME)
STATIC_LIB = $(COMPILE_LIB_HOME)/$(STATIC_NAME)

all: $(DYNAMIC_LIB) $(STATIC_LIB)

$(DYNAMIC_LIB): $(SRCS:.c=.o)
	$(CC) -pg -o $@ $^ $(CFLAGS)

$(STATIC_LIB): $(SRCS:.c=.o)
	@ar cr $@ $^

%.o: %.c
	$(CC) $(CFLAGS) $(INC) -c -pg -o $@ $<

clean:
	rm -f *.o .po *.so *.d $(SRCS:.c=.o) $(SRCS:.c=.d) $(DYNAMIC_LIB) $(STATIC_LIB) $(DYNAMIC_NAME) $(STATIC_NAME)
