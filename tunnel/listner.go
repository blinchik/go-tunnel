package tunnel

import (
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

const (
	sshPort = "22"
)

// LocalListner start listner on the local address and call Forward function
func LocalListner(cleintTarget *ssh.Client, localAddress string, remoteAddress string) {

	listener, err := net.Listen("tcp", localAddress)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {

		conn, err := listener.Accept()

		if err != nil {
			panic(err)
		}

		if conn == nil {
			panic("conn is nil")
		}

		Forward(cleintTarget, conn, remoteAddress)

	}

}

// Forward dial to the target client and start io.copy fron local-remote and remote-local
func Forward(cleintTarget *ssh.Client, localConn net.Conn, remoteAddress string) {

	remoteConn, err := cleintTarget.Dial("tcp", remoteAddress)
	if err != nil {

		localConn.Close()

		if remoteConn != nil {

			remoteConn.Close()

		}

		panic(err)

	}

	copyConn := func(writer, reader net.Conn) {

		if writer != nil {
			if reader != nil {

				_, err := io.Copy(writer, reader)

				if err != nil {

					log.Println("copyConn - err: ", err)

					writer.Close()
					reader.Close()

				}
			}
		}
		return

	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)

}

// ListenOverMiddle creates SSH tunnel over middle machine
func ListenOverMiddle(userBastion string, userTarget string, bastionKey string, targetKey string, bastionAddress string, TargetAddress, localListnerAddr string, remoteListnerAddr string) {

	clientBastion := FirstClient(userBastion, bastionKey, bastionAddress, sshPort, "local")
	cleintTarget := TargetClient(clientBastion, userTarget, targetKey, TargetAddress, sshPort, "local")
	LocalListner(cleintTarget, localListnerAddr, remoteListnerAddr)

}

// ListenDirect creates SSH tunnel directly to the remote machine
func ListenDirect(targetKey string, userTarget string, TargetAddress, localListnerAddr string, remoteListnerAddr string) {

	cleintTarget := FirstClient(userTarget, targetKey, TargetAddress, sshPort, "local")
	LocalListner(cleintTarget, localListnerAddr, remoteListnerAddr)

}
