#!/bin/sh

set -x
set -e

cd "$(dirname "$0")"
make clean
make
cp ./mlpx.so ./mlpx.h ./test
cd ./test

for f in test_*.c ; do
	cc "$f" ./mlpx.so
	./a.out
	rm -f a.out
done

rm -f ./mlpx.so
rm -f ./mlpx.h
