include ../opinionated.mk

wavegen.h: wavegen.so
.PHONY: wavegen.h

wavegen.so: wavegen.go
> go build -o "$@" -buildmode=c-shared $^

wavegen.a: wavegen.so
> ar rcs $@ $<

clean:
> rm -f wavegen.so wavegen.h wavegen.a
.PHONY: clean
