# Makefile to run the output of a Zig binary
#
# # Variables
ZIG_BINARY := zig-out/bin/hurl
# OUTPUT_FILE := output
#
# # Default target
# all: run

# Build the Zig binary
build:
	zig build

# Run the Zig binary
run: build
	./$(ZIG_BINARY)

# # Clean generated files
# clean:
#     rm -f $(ZIG_BINARY) $(OUTPUT_FILE)

