include ../opinionated.mk

mlpx.h: mlpx.so
.PHONY: mlpx.h

mlpx.so: mlpx.go
> go build -o "$@" -buildmode=c-shared $^

mlpx.a: mlpx.so
> ar rcs $@ $<

test:
> sh ./run_tests.sh
.PHONY: test

clean:
> rm -f mlpx.so mlpx.h mlpx.a
> rm -f test/mlpx.so test/mlpx.h test/a.out
.PHONY: clean
