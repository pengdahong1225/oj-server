cmake_minimum_required(VERSION 3.17)

set(TARGET judge-service)
set(CMAKE_CXX_STANDARD 20)

if (WIN32)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} /Zc:__cplusplus")
endif ()

project(${TARGET})

include_directories(src)
include_directories(src/app)
include_directories(src/processer)
include_directories(../../proto)

aux_source_directory(./src SRCS)
aux_source_directory(./src/app SRCS)
aux_source_directory(./src/processer SRCS)

add_executable(${TARGET} ${SRCS})
