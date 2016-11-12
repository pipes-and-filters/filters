#!/bin/bash
echo "http request to localhost:8080 with post body:
=====================
hello world
goodbye
=====================

results:
====================="
curl -X POST -d $'hello world\ngoodbye' http://localhost:8080
echo "====================="
