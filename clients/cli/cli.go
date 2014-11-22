// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	tool "github.com/Alienero/quick-know/restful_tool"

	"github.com/codegangsta/cli"
)

func init() {
	tool.Url = "http://127.0.0.1:9901"
	tool.ID = "1234"
	tool.Psw = "10086"
}

func main() {
	app := cli.NewApp()
	app.Name = "cli"
	app.Usage = "control the quick-know server in command line"
	app.Commands = []cli.Command{
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "add a thins to the server",
			Subcommands: []cli.Command{
				{
					Name:      "user",
					ShortName: "u",
					Usage:     "add a user",
					Action: func(c *cli.Context) {
						id, err := tool.AddUser(c.Args().First())
						check_err(err)
						log.Println("user id is :", id)
					},
				},
				{
					Name:      "sub",
					ShortName: "s",
					Usage:     "add a sub",
					Action: func(c *cli.Context) {
						id, err := tool.AddSub(0, 0)
						check_err(err)
						log.Println("sub id is:", id)
					},
				},
				{
					Name:      "user2sub",
					ShortName: "us",
					Usage:     "add a user to sub(sub_id,user_id)",
					Action: func(c *cli.Context) {
						check_err(tool.User2Sub(c.Args().Get(0), c.Args().Get(1)))
					},
				},
				{
					Name:      "pmsg",
					ShortName: "pm",
					Usage:     "add a private message",
					Action: func(c *cli.Context) {
						check_err(tool.AddPrivateMsg(c.Args().Get(0), 0, []byte(c.Args().Get(1))))
					},
				},
				{
					Name:      "gmsg",
					ShortName: "gm",
					Usage:     "add a group message",
					Action: func(c *cli.Context) {
						check_err(tool.GroupMsg(c.Args().Get(0), 0, []byte(c.Args().Get(1))))
					},
				},
				{
					Name:      "wmsg",
					ShortName: "wm",
					Usage:     "add a message to whole users",
					Action: func(c *cli.Context) {
						check_err(tool.Broadcast(0, []byte(c.Args().First())))
					},
				},
			},
		},
		// TODO :DELETE CMD!!
		// {
		// 	Name:      "complete",
		// 	ShortName: "c",
		// 	Usage:     "complete a task on the list",
		// 	Action: func(c *cli.Context) {
		// 		println("completed task: ", c.Args().First())
		// 	},
		// },
	}
	app.Run(os.Args)
}

func check_err(err error) {
	if err != nil {
		log.Printf("oops! error:%v\n", err)
	}
}
