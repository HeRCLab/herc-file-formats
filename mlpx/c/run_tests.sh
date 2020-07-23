#!/bin/sh

set -e

cd "$(dirname "$0")"
make clean
make
cp ./mlpx.so ./mlpx.h ./test
cd ./test

for f in test_*.c ; do
	set -x
	cc "$f" ./mlpx.so
	./a.out
	set +x
	rm -f a.out
done

rm -f ./mlpx.so
rm -f ./mlpx.h
