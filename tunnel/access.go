package tunnel

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func AddSshPubKey(client *ssh.Client, user string, keyString string) {

	homeName := fmt.Sprintf("/home/%s", user)

	command := []string{
		fmt.Sprintf("echo \"%s\" >> %s/.ssh/authorized_keys", keyString, homeName),
		"exit",
		"exit",
	}

	ExecuteCommands(client, command)

}

func DeleteSshPubKey(client *ssh.Client, keyTag string) {

	command := []string{
		fmt.Sprintf("sed --in-place '/%s/d' .ssh/authorized_keys", keyTag),
		"exit",
		"exit",

	}

	ExecuteCommands(client, command)

}
