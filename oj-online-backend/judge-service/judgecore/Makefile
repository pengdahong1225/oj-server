CC  = gcc
GCCVER := $(shell $(CC) -dumpversion | awk -F. '{ print $$1"."$$2}' )

OPT     = -pipe -fno-ident -fPIC -shared -z defs
CFLAGS += $(OPT) -g -pg -Wall -D_GNU_SOURCE -funroll-loops -MMD -D_REENTRANT -Wno-invalid-offsetof
ifeq ($(MEMCHECK),1)
CFLAGS += -DMEMCHECK
endif

INC	= -I./src -I./src/rules
SRCS = $(wildcard ./src/*.c) $(wildcard ./src/rules/*.c)

COMPILE_LIB_HOME = .
DYNAMIC_NAME = core.so
STATIC_NAME = core.a
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
