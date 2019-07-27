#!/bin/bash
go build
./scc-data
python convert_json.py
rm scc-data