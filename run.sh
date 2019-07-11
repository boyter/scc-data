#!/bin/bash
#find json/ | xargs -n 1 -P 8 ./parse.py | ./summer.py
go build
./scc-data
python convert_json.py
rm scc-data