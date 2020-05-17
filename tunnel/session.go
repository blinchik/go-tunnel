package tunnel

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

// ExecuteCommands execute a list of shell commands on the client machine
func ExecuteCommands(cleint *ssh.Client, commands []string) {

	session, err := cleint.NewSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()
	// StdinPipe for commands
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Enable system stdout
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Start remote shell
	err = session.Shell()
	if err != nil {
		log.Fatal(err)
	}

	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Wait for session to finish
	err = session.Wait()
	if err != nil {
		log.Fatal(err)
	}

}
