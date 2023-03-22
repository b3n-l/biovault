package lib

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

func runCommand(name string, workingDir string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = workingDir
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func PadToMultipleOf8(input []byte) []byte {
	remainder := len(input) % 8
	if remainder == 0 {
		return input
	}
	padding := make([]byte, 8-remainder)
	return append(input, padding...)
}

// Dump the contents of the card to a temporary file and return the decoded bytes
func Dump(sector int) []byte {
	dir, err := os.MkdirTemp("", "")
	defer os.RemoveAll(dir)
	cmd := [...]string{"pm3", "-c", fmt.Sprintf("script run hf_i2c_plus_2k_utils -s %d -m d", sector)}

	out, err := runCommand(cmd[0], dir, cmd[1:]...)
	if err != nil {
		log.Fatalln(err)
	}
	// Find the generated hex file
	r := regexp.MustCompile("\\w{14}\\.hex")
	generatedFile, err := os.ReadFile(dir + "/" + r.FindString(out))
	if err != nil {
		log.Fatalln("Error reading data from temporary file")
	}
	byteOut, _ := hex.DecodeString(string(generatedFile))
	return byteOut
}

func WriteSector(sector int, data []byte) {
	tempFile, err := os.Create("tempfile-write")
	defer os.Remove(tempFile.Name())
	if err != nil {
		log.Fatalln("Could not create temporary file", err)
	}
	_, err = tempFile.Write(data)

	if err != nil {
		log.Fatalln("Could not write data to temporary file:", err)
	}
	err = tempFile.Sync()
	if err != nil {
		log.Fatalln("Could not sync file", err)
	}

	cmd := [...]string{"pm3", "-c", fmt.Sprintf("script run hf_i2c_plus_2k_utils -s %d -m f -f %s", sector, tempFile.Name())}
	_, err = runCommand(cmd[0], "", cmd[1:]...)
	if err != nil {
		log.Fatalln("Got error when writing file:", err)
	}
}
