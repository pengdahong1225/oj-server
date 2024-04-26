cmake_minimum_required(VERSION 3.23)

set(TARGET core-wrapper)
set(CMAKE_CXX_STANDARD 20)

if (WIN32)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} /Zc:__cplusplus")
endif ()

project(${TARGET})

include_directories(core)
include_directories(core/src)
include_directories(wrapper)
include_directories(common/json)
include_directories(judgeclient)

include_directories(test)
include_directories(test/src)

aux_source_directory(. SRCS)
aux_source_directory(core SRCS)
aux_source_directory(core/src SRCS)
aux_source_directory(wrapper SRCS)
aux_source_directory(common SRCS)
aux_source_directory(judgeclient SRCS)

aux_source_directory(test SRCS)
aux_source_directory(test/src SRCS)

add_executable(${TARGET} ${SRCS})
