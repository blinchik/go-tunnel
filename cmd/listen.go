package main

import (
	"flag"
	"fmt"
	"log"
	"fmt"

	mEC2 "github.com/blinchik/go-aws/lib/manage-ec2"
	"github.com/blinchik/go-tunnel/tunnel"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	bastionFlag := flag.Bool("bastion", false, "connect through bastion")
	freetag := flag.Bool("freetag", false, "connect through bastion")
	giveip := flag.Bool("giveip", false, "connect through bastion")

	flag.Parse()

	if *bastionFlag {
		if *giveip {

			userBastion := flag.Arg(0)
			userTarget := flag.Arg(1)
			bastionKey := flag.Arg(2)
			targetKey := flag.Arg(3)
			bastionAddress := flag.Arg(4)
			TargetAddress := flag.Arg(5)
			localListnerAddr := flag.Arg(6)
			remoteListnerAddr := flag.Arg(7)

			flag.Args()

			go tunnel.Recoverer(300, "tunnel listen", func() {

				tunnel.ListenOverMiddle(userBastion, userTarget, bastionKey, targetKey, bastionAddress, TargetAddress, localListnerAddr, remoteListnerAddr)
			})

			select {}

		} else if *freetag {

			userBastion := flag.Arg(0)
			userTarget := flag.Arg(1)
			bastionKey := flag.Arg(2)
			targetKey := flag.Arg(3)
			freetagBastion := flag.Arg(4)
			freetagTarget := flag.Arg(5)
			ec2BastionName := flag.Arg(6)
			ec2TargetName := flag.Arg(7)
			localListnerAddr := flag.Arg(8)
			remoteListnerAddr := flag.Arg(9)

			descBastion := mEC2.DescribeByGeneralTag(freetagBastion, ec2BastionName)
			descTarget := mEC2.DescribeByGeneralTag(freetagTarget, ec2TargetName)

			fmt.Println(descTarget, descBastion)

			TargetAddress := descTarget.PrivateIpAddress
			bastionAddress := descBastion.PublicIp
			
			fmt.Println(TargetAddress,bastionAddress)
			fmt.Println(descBastion,descTarget)

			flag.Args()

			go tunnel.Recoverer(300, "tunnel listen", func() {

				tunnel.ListenOverMiddle(userBastion, userTarget, bastionKey, targetKey, *bastionAddress[0], *TargetAddress[0], localListnerAddr, remoteListnerAddr)

			})

			select {}

		} else {

			userBastion := flag.Arg(0)
			userTarget := flag.Arg(1)
			bastionKey := flag.Arg(2)
			targetKey := flag.Arg(3)
			ec2BastionName := flag.Arg(4)
			ec2TargetName := flag.Arg(5)
			localListnerAddr := flag.Arg(6)
			remoteListnerAddr := flag.Arg(7)

			flag.Args()

			descBastion := mEC2.DescribeByGeneralTag("Name", ec2BastionName)
			descTarget := mEC2.DescribeByGeneralTag("Name", ec2TargetName)

			TargetAddress := descTarget.PrivateIpAddress
			bastionAddress := descBastion.PublicIp

			fmt.Println(TargetAddress, bastionAddress)

			go tunnel.Recoverer(300, "tunnel listen", func() {

				tunnel.ListenOverMiddle(userBastion, userTarget, bastionKey, targetKey, *bastionAddress[0], *TargetAddress[0], localListnerAddr, remoteListnerAddr)

			})
			select {}

		}
	} else {

		userTarget := flag.Arg(0)
		targetKey := flag.Arg(1)
		ec2TargetName := flag.Arg(2)
		localListnerAddr := flag.Arg(3)
		remoteListnerAddr := flag.Arg(4)

		flag.Args()

		descTarget := mEC2.DescribeByGeneralTag("Name", ec2TargetName)
		TargetAddress := descTarget.PrivateIpAddress

		go tunnel.Recoverer(300, "tunnel listen", func() {

			tunnel.ListenDirect(userTarget, targetKey, *TargetAddress[0], localListnerAddr, remoteListnerAddr)

		})
		select {}

	}

}
