cmake_minimum_required(VERSION 3.20)

set(TARGET judge-service)
set(CMAKE_CXX_STANDARD 20)

if (WIN32)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} /Zc:__cplusplus")
endif ()

project(${TARGET})

include_directories(../library)
include_directories(../proto)

include_directories(./app)
include_directories(./handler)
aux_source_directory(./app SRCS)
aux_source_directory(./handler SRCS)
aux_source_directory(. SRCS)

add_executable(${TARGET} ${SRCS})