#!/bin/sh
docker run -it --rm alexeiled/stress-ng --vm-bytes $(awk '/MemFree/{printf "%d\n", $2 * 1;}' < /proc/meminfo)k --vm-keep -m 10

