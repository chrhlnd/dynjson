#!/bin/bash


echo "---------------------------" >> ./bench_history.txt
echo $(date) >> ./bench_history.txt
go test --bench='.*' >> ./bench_history.txt
