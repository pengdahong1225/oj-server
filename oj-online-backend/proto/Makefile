TOOL_PROTOC   := /usr/bin/protoc

LIB = libproto.a

CXXFLAGS += -std=c++2a -Wall -g -pipe -rdynamic -fno-strict-aliasing -Wno-unused-function -Wno-sign-compare -fpermissive

PB_OBJ = judge.o

INC_FLAGS =  -I.

$(LIB) :  $(PB_OBJ)
	ar rcs $@ $^
        
$(PB_OBJ) :%.o : %.proto
	$(TOOL_PROTOC) $< --cpp_out=./
		g++ $(CXXFLAGS) $(INC_FLAGS) -c -o $@ $(patsubst %.proto, %.pb.cc, $<)

clean:
        $(RM) $(LIB) $(DEP_DIR) -rf 

.PHONY: protbuf

