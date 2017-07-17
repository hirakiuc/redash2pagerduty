#!/bin/bash

WRK=/usr/local/bin/wrk

$WRK -s ./post.lua -c 3 -t 3 http://localhost:8080/
