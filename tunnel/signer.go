package tunnel

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// SignLocalKey create signatures from private keys under the local path ~/.ssh that verify against a public key.
func SignLocalKey(KeyName string) ssh.Signer {

	pathToPrivateKey := filepath.FromSlash(fmt.Sprintf("%s/.ssh/%s", home, KeyName))

	key, err := ioutil.ReadFile(pathToPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal(err)
	}

	return signer

}
