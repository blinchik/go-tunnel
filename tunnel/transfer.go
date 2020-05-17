package tunnel

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

// TransferFiles use SFTP to transfer files into clinet machine
func TransferFiles(client *ssh.Client, localDir string, targetDir string) {

	// create new SFTP client
	clientSFTP, err := sftp.NewClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer clientSFTP.Close()

	// create destination file
	dstFile, err := clientSFTP.Create(targetDir)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// create source file
	srcFile, err := os.Open(localDir)
	if err != nil {
		log.Fatal(err)
	}

	// copy source file to destination file
	bytesd, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d bytes copied\n", bytesd)
}
