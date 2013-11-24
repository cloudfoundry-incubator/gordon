package main

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/vito/gordon"
	"github.com/vito/gordon/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "gordon"
	app.Usage = "manage warden containers"
	app.Flags = []cli.Flag{
		cli.StringFlag{"socket", "/tmp/warden.sock", "path to the warden command socket"},
	}

	ui := commands.BasicUI{
		Writer: os.Stdout,
	}

	app.Commands = []cli.Command{
		{
			Name:  "list",
			Usage: "list running containers",
			Action: func(c *cli.Context) {
				status := commands.NewList(client(c), ui)
				status.Run()
			},
		},
		{
			Name:  "create",
			Usage: "create a container",
			Action: func(c *cli.Context) {
				status := commands.NewCreate(client(c), ui)
				status.Run()
			},
		},
	}

	app.Run(os.Args)
}

func client(c *cli.Context) warden.Client {
	connectionInfo := &warden.ConnectionInfo{
		SocketPath: c.GlobalString("socket"),
	}
	client := warden.NewClient(connectionInfo)
	client.Connect()

	return client
}
