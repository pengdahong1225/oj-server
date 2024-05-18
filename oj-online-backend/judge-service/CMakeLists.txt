cmake_minimum_required(VERSION 3.20)

set(TARGET judge-service)
set(CMAKE_CXX_STANDARD 20)

if (WIN32)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} /Zc:__cplusplus")
endif ()

project(${TARGET})

include_directories(../library)
include_directories(../proto)
include_directories(.)

aux_source_directory(./common SRCS)
aux_source_directory(./common/logger SRCS)
aux_source_directory(./wrapper SRCS)
aux_source_directory(./app SRCS)
aux_source_directory(./handler SRCS)
aux_source_directory(./judger SRCS)
aux_source_directory(./judger/judgeclient SRCS)
aux_source_directory(. SRCS)
add_executable(${TARGET} ${SRCS})