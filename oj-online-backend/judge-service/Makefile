CXX = g++
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

# judger-wrapper
JUDGER_INCLUDE = ../library
JUDGER_LIBRARY = ../library/libs
JUDGER_LIB = -L$(JUDGER_LIBRARY) -ljudger

# app
TARGET = JudgeService
INC	+= -I$(MUDUO_INCLUDE) -I$(JUDGER_INCLUDE) -I./app -I./handler
LIB += $(MUDUO_LIB) $(JUDGER_LIB)

SRC += $(wildcard *.cpp) $(wildcard ./app/*.cpp) $(wildcard ./handler/*.cpp)

# build
$(TARGET): $(SRC)
	$(CXX) $(CXXFLAGS) $(INC) -c -o $@ $^ $(LIB)

all: $(TARGET)
clean:
	rm -rf $(TARGET)

.PHONY: all clean
