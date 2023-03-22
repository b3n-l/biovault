# Biovault

This is a rewrite of [Flamebarke's Biovault](https://github.com/flamebarke/biovault) in Go.
The aims of this rewrite was to eliminate calls to external applications like openssl, xxd etc. wherever possible.

I also wanted to be able to take arbitrary files and write them to the xSIID implant.

## Features:
- Writing arbitrary files to sector 0 or sector 1 of the NTAG I2C/xSIID implant.
- Write verification using SHA256 checksums
- Optional AES256GCM encryption/decryption with PKFBD2 and 500,000 iterations 
- Intermediate files are generated in temp folders and cleaned up afterwards

## Usage
```
  -d    Decrypt?
  -e    Encrypt?
  -f string
        Filename to read or write
  -r    Read from tag
  -s int
        Sector to read/write (default 1)
  -stdout
        Write data to stdout
  -w    Write to tag
  -zero
        Zero sector before writing
```