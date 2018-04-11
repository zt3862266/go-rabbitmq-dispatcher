#!/bin/bash


echo "kill old process ..."
kill -INT `cat run/grd.pid`
exit 0