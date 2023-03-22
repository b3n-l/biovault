package main

import (
	"biovault/lib"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strings"
)

func main() {
	var writeStdOut = flag.Bool("stdout", false, "Write data to stdout")
	var filePath = flag.String("f", "", "Filename to read or write")
	var zero = flag.Bool("zero", false, "Zero sector before writing")
	var read = flag.Bool("r", false, "Read from tag")
	var sector = flag.Int("s", 1, "Sector to read/write")
	var write = flag.Bool("w", false, "Write to tag")
	var encrypt = flag.Bool("e", false, "Encrypt?")
	var decrypt = flag.Bool("d", false, "Decrypt?")
	flag.Parse()

	if *read && *write {
		log.Fatalln("You can't read and write at the same time")
	}

	// Reading output to stdout
	if *read {
		doRead(*sector, *decrypt, *writeStdOut, *filePath)
		return
	} else if *write {
		if *zero {
			lib.WriteSector(*sector, make([]byte, 1024))
		}
		doWrite(*sector, *encrypt, *filePath)
		return
	}
	flag.PrintDefaults()
}

func getPassword() []byte {
	fmt.Println("Enter password:")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatalln("Unable to read password")
	}
	return password
}

func doWrite(sector int, encrypt bool, filePath string) {
	log.Println(fmt.Sprintf("Writing %s to sector %d", filePath, sector))
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Unable to read supplied file:", err)
	}
	var dataToWrite []byte
	var password []byte
	if encrypt {
		password = getPassword()
		fmt.Println("Confirm password:")
		password2 := getPassword()
		if strings.Compare(string(password), string(password2)) != 0 {
			log.Fatalln("Passwords didn't match")
		}
		file := lib.PadToMultipleOf8(file)
		dataToWrite, err = lib.EncryptBytes(file, password)
		if err != nil {
			log.Fatalln("Error encrypting data:", err)
		}
	} else {
		dataToWrite = file
	}
	log.Println("Writing bytes: ", hex.EncodeToString(dataToWrite))
	// Prepend the length that we are writing as two bytes, this allows multiple records
	lengthSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthSlice, uint16(len(dataToWrite)))
	fmt.Println("Press the Enter Key to start writing...")
	fmt.Scanln()
	lib.WriteSector(sector, lib.PadToMultipleOf8(append(lengthSlice, dataToWrite...)))
	log.Println("Sector written successfully")

	// Verify write once we have written
	readBack := lib.Dump(sector)
	sourceHash := sha256.Sum256(dataToWrite)
	readLength := binary.BigEndian.Uint16(readBack[16:18])
	readBackTruncated := readBack[18 : readLength+18]

	readBackHash := sha256.Sum256(readBackTruncated)
	if sourceHash != readBackHash {
		log.Fatalln("Error comparing source and written hashes, got ", hex.EncodeToString(sourceHash[:]), hex.EncodeToString(readBackHash[:]))
	} else {
		log.Println("Write verified")
	}

	// Test that the data can be decrypted again with the same key
	if encrypt {
		decrypted, err := lib.DecryptBytes(readBackTruncated, password)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(decrypted))
	}
}

func doRead(sector int, decrypt bool, writeStdOut bool, filePath string) {
	log.Println(fmt.Sprintf("Dumping sector %d", sector))
	data := lib.Dump(sector)

	readLength := binary.BigEndian.Uint16(data[16:18])

	// We want to strip off the first 16 bytes, plus the 2 byte length header
	dataTruncated := data[18 : readLength+18]

	if decrypt {
		var err error
		dataTruncated, err = lib.DecryptBytes(dataTruncated, getPassword())
		if err != nil {
			log.Fatalln(err)
		}
	}

	if writeStdOut {
		fmt.Println(string(dataTruncated))
	}
	if filePath != "" {
		file, err := os.Create(filePath)
		defer file.Close()
		if err != nil {
			log.Fatalln("Unable to open output file", err)
		}
		_, err = file.Write(data)
		if err != nil {
			log.Fatalln("Unable to write file")
		}
	}
	return
}
