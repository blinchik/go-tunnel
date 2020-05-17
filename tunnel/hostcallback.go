package tunnel

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

// HostKeyCallback creates a host key callback from the given OpenSSH host key files. It assumes that the file is located at "<home>/.ssh/known_hosts"
func HostKeyCallback(dialAddr string, addr net.Addr, key ssh.PublicKey) error {

	kh.New(filepath.FromSlash(fmt.Sprintf("%s/.ssh/known_hosts", home)))
	if err != nil {
		log.Fatal("could not create hostkeycallback function: ", err)
	}

	return nil
}

// Harvest dial to the target machine and add the key from ssh response to the known hosts. It assumes that your known_hosts file has an aws header as a comment, i.e. #AWS
func Harvest(dialAddr string, addr net.Addr, key ssh.PublicKey) error {

	returnedKey := fmt.Sprintf("%s %s %s\n", strings.Split(dialAddr, ":")[0], key.Type(), base64.StdEncoding.EncodeToString(key.Marshal()))

	KnownHostsPath := filepath.FromSlash(filepath.FromSlash(fmt.Sprintf("%s/.ssh/known_hosts", home)))

	file, err := os.OpenFile(KnownHostsPath, os.O_RDWR, 0644)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	KnownHosts, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(string(KnownHosts), returnedKey) {
		log.Println("\nhost already known\t", dialAddr)
	} else {

		KnownHostsUpdated := strings.Replace(string(KnownHosts), "#AWS", fmt.Sprintf("#AWS\n%s", returnedKey), -1)

		file, err = os.OpenFile(KnownHostsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

		w := bufio.NewWriter(file)

		_, err = w.WriteString(KnownHostsUpdated)

		w.Flush()

		fmt.Printf("### add %s %s %s\n", strings.Split(dialAddr, ":")[0], key.Type(), base64.StdEncoding.EncodeToString(key.Marshal()))

	}

	return nil
}

// GetKnownhostsPublic will harvest the known hosts from machines that located in public subnets
func GetKnownhostsPublic(user string, signer ssh.Signer, address string, port string) {

	config := &ssh.ClientConfig{
		Timeout: (5 * time.Second),

		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: Harvest,
	}

	client, err := ssh.Dial("tcp", address+":"+port, config)

	if err != nil {
		log.Println(err)
	}

	defer client.Close()

}

// GetKnownhostsPrivate will harvest the known hosts from machines that located in private subnets over SSH tunnel
func GetKnownhostsPrivate(cleint *ssh.Client, user string, signer ssh.Signer, address string, port string) {

	conn, err := cleint.Dial("tcp", address+":"+port)
	if err != nil {
		log.Println(err)
	}

	config := &ssh.ClientConfig{
		Timeout: (5 * time.Second),

		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: Harvest,
	}

	ncc, chans, reqs, err := ssh.NewClientConn(conn, address+":"+port, config)
	if err != nil {
		log.Println(err)
	}

	newCleint := ssh.NewClient(ncc, chans, reqs)

	defer newCleint.Close()

}
