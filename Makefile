chip8: chip8.go
	go build -o chip8 chip8.go

.PHONY: clean
clean:
	rm -f chip8 