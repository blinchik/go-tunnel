//Package tunnel facilitate the work with remote machines over ssh tunnels
package tunnel

import (
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

//FirstClient will create SSH client with remote machine
func FirstClient(user string, key string, address string, port string, method string) *ssh.Client {

	var signer ssh.Signer

	if method == "local" {

		signer = SignLocalKey(key)

	}

	if method == "key" {

		// Create the Signer for this private key.
		signer, err = ssh.ParsePrivateKey([]byte(key))
		if err != nil {
			log.Fatal(err)
		}
	}

	GetKnownhostsPublic(user, signer, address, port)

	config := &ssh.ClientConfig{
		Timeout: (5 * time.Second),

		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: HostKeyCallback,
	}

	client, err := ssh.Dial("tcp", address+":"+port, config)

	if err != nil {
		log.Println("FirstClient - ssh dial err: ", err)
	}

	return client
}

//TargetClient will create SSH client with remote machine through the middle machine
func TargetClient(cleint *ssh.Client, user string, key string, address string, port string, method string) *ssh.Client {

	var signer ssh.Signer

	if method == "local" {

		signer = SignLocalKey(key)

	}

	if method == "key" {

		// Create the Signer for this private key.
		signer, err = ssh.ParsePrivateKey([]byte(key))
		if err != nil {
			log.Fatal(err)
		}
	}

	GetKnownhostsPrivate(cleint, user, signer, address, port)

	conn, err := cleint.Dial("tcp", address+":"+port)
	if err != nil {
		log.Println("TargetClient - ssh dial err: ", err)
	}

	config := &ssh.ClientConfig{
		Timeout: (5 * time.Second),
		User:    user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: HostKeyCallback,
	}

	ncc, chans, reqs, err := ssh.NewClientConn(conn, address+":"+port, config)
	if err != nil {
		log.Println("NewClientConn ", err)

	}

	newCleint := ssh.NewClient(ncc, chans, reqs)

	return newCleint

}
