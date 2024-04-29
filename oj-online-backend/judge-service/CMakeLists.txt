cmake_minimum_required(VERSION 3.26)

set(TARGET judge-service)
set(CMAKE_CXX_STANDARD 20)

if (WIN32)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} /Zc:__cplusplus")
endif ()

project(${TARGET})

include_directories(
        ../library
        ../judger-wrapper
)

include_directories(../proto)

include_directories(.)
include_directories(./app)
include_directories(./processer)
aux_source_directory(./ SRCS)
aux_source_directory(./app SRCS)
aux_source_directory(./processer SRCS)


add_executable(${TARGET} ${SRCS})
target_link_libraries(${TARGET} PRIVATE MUDUO_LIBRARY)