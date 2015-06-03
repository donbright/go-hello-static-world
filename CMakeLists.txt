# Copyright 2012 Douglas Linder
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# This file will help to build and statically-link a C language library
# file (.a on un*x, .lib on win*) to the Go language. It creates a 
# .go "bridge" file automatically from a template 'go.in' file, by 
# modifying the paths in the template file.
#


cmake_minimum_required (VERSION 2.8)
project (Hello)

# Set up library file name. In spirit of Go "Export" syntax, use 
# capitalization for the library name.
set(Hello "Hello")

# Debug mode?
set(CMAKE_BUILD_TYPE Debug)

# Flags
if(NOT WIN32)
  set(CMAKE_C_FLAGS "-Wall -pedantic -Wimplicit-function-declaration -Wreturn-type -Wstrict-prototypes -Wextra -fPIC")
  set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Werror")
  set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -std=c99")
endif()

# Collect source files
file(GLOB_RECURSE SOURCES src/*.c*)

# libary headers
include_directories(${PROJECT_SOURCE_DIR}/include)

# Build library files (.a files) under 'lib' subdir
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${PROJECT_BINARY_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${PROJECT_BINARY_DIR}/lib)
#set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${PROJECT_BINARY_DIR}/bin)

# Build library files
add_library(${Hello} ${SOURCES})
target_link_libraries(${Hello})

# Create .go bridge file using 'go.in' template. Mostly by modifying paths
# to point to the proper include directory and the library files

set(INFILE ${PROJECT_SOURCE_DIR}/demo/src/bridge/bridge.go.in)
set(OUTFILE ${PROJECT_BINARY_DIR}/src/bridge/bridge.go)
set(C_INCLUDE_PATH ${PROJECT_SOURCE_DIR}/include)
set(C_LIBRARY_PATH $<TARGET_FILE:Hello>)

message(STATUS "Creating .go bridge file ${OUTFILE}")
message(STATUS "from template input file ${INFILE}" )
message(STATUS "replacing C_INCLUDE_PATH with ${C_INCLUDE_PATH}")
message(STATUS "replacing C_LIBRARY_PATH with ${C_LIBRARY_PATH}")

configure_file( ${INFILE} ${OUTFILE} )

