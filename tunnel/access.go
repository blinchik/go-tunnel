package tunnel

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func AddSshPubKey(client *ssh.Client, user string, keyName string) {

	homeName := fmt.Sprintf("/home/%s", user)

	localDir := fmt.Sprintf("./pub/%s", keyName)
	targetDir := fmt.Sprintf("./%s", keyName)

	command := []string{
		fmt.Sprintf("sudo su; echo \"`cat ./%s`\" >> %s/.ssh/authorized_keys", keyName, homeName),
		fmt.Sprintf("rm ./%s", keyName),
		"exit",
	}

	TransferFiles(client, localDir, targetDir)
	ExecuteCommands(client, command)

}

func DeleteSshPubKey(client *ssh.Client, keyTag string) {

	command := []string{
		fmt.Sprintf("sed --in-place '/%s/d' .ssh/authorized_keys", keyTag),
		"exit",
	}

	ExecuteCommands(client, command)

}
