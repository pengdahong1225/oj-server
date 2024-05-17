CC  = gcc
CXX = g++

CFLAGS = -g -Wall -Werror -O3 -std=c99
CXXFLAGS = -std=c++2a -g -O0 -Wall -Wextra \
		-Wconversion -Wno-unused-parameter \
		-Wold-style-cast -Woverloaded-virtual \
		-Wpointer-arith -Wshadow -Wwrite-strings \
		-rdynamic -fno-strict-aliasing -Wno-unused-function \
		-Wno-sign-compare -fpermissive -Wno-invalid-offsetof

# muduo
MUDUO_INCLUDE = ../library
MUDUO_LIBRARY = ../library/libs
MUDUO_LIB = -L$(MUDUO_LIBRARY) -lmuduo_net -lmuduo_base -lpthread

# proto
PROTO_DIR = ../proto
PROTO_INCLUDE = -I$(PROTO_DIR)
PROTO_LIB = -lprotobuf $(PROTO_DIR)/libproto.a

# app
TARGET = JudgeService
INC	+= -I$(MUDUO_INCLUDE) -I. -I./judger/core $(PROTO_INCLUDE)
LIB += $(MUDUO_LIB) $(JUDGER_LIB) -lseccomp $(PROTO_LIB)

C_SRC += $(wildcard ./judger/core/*.c) $(wildcard ./judger/core/rules/*.c)
CPP_SRC += $(wildcard *.cpp) $(wildcard ./app/*.cpp) $(wildcard ./handler/*.cpp)
CPP_SRC += $(wildcard ./common/*.cpp) $(wildcard ./judger/judgeclient/*.cpp) $(wildcard ./wrapper/*.cpp)

C_OBJ = $(patsubst %.c,%.o,$(C_SRC))
CPP_OBJ = $(patsubst %.cpp,%.o,$(CPP_SRC))

OBJ = $(C_OBJ) $(CPP_OBJ)

# build
$(TARGET): $(OBJ)
	$(CXX) $(CXXFLAGS) -o $@ $^ $(LIB)

%.o: %.c
	$(CC) $(CFLAGS) $(INC) -c -o $@ $<

%.o: %.cpp
	$(CXX) $(CXXFLAGS) $(INC) -c -o $@ $<

all: $(TARGET)
clean:
	rm -rf $(TARGET) $(OBJ)

.PHONY: all clean
