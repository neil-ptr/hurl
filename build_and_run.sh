#! /bin/sh
# rebuild prog if necessary
make build
# run prog with some arguments
./hurl "$@"
