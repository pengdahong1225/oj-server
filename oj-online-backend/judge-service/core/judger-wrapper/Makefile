CC  = gcc
CXX = g++

CFLAGS = -g -Wall -Werror -O3 -std=c99
CXXFLAGS = -std=c++2a -Wall -g -pipe -rdynamic -fno-strict-aliasing -Wno-unused-function -Wno-sign-compare -fpermissive -Wno-invalid-offsetof

#LIB = ./core/libcore.a
LINK = -lseccomp # libseccomp是一个用于 Linux 系统的 syscall 过滤器库

INC	+= -I./core -I./core/rules -I./wrapper -I./common -I./common/json -I./judgeclient

C_SRC += $(wildcard ./core/*.c) $(wildcard ./core/rules/*.c)
CPP_SRC += $(wildcard ./common/*.cpp) $(wildcard ./judgeclient/*.cpp) $(wildcard ./wrapper/*.cpp)

C_OBJ = $(patsubst %.c,%.o,$(C_SRC))
CPP_OBJ = $(patsubst %.cpp,%.o,$(CPP_SRC))

OBJ = $(C_OBJ) $(CPP_OBJ)

COMPILE_LIB_HOME = .
DYNAMIC_NAME = libjudger.so
STATIC_NAME = libjudger.a
DYNAMIC_LIB	= $(COMPILE_LIB_HOME)/$(DYNAMIC_NAME)
STATIC_LIB = $(COMPILE_LIB_HOME)/$(STATIC_NAME)

all: $(DYNAMIC_LIB) $(STATIC_LIB)

$(DYNAMIC_LIB): $(OBJ)
	$(CXX) -pg -o $@ $^ $(CXXFLAGS) $(LINK)

$(STATIC_LIB): $(OBJ)
	@ar cr $@ $^

%.o: %.c
	$(CC) $(CFLAGS) $(INC) -c -pg -o $@ $< $(LINK)

%.o: %.cpp
	$(CXX) $(CXXFLAGS) $(INC) -c -pg -o $@ $< $(LINK)

clean:
	rm -rf $(OBJ) $(DYNAMIC_LIB) $(STATIC_LIB) $(DYNAMIC_NAME) $(STATIC_NAME)
