package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/KosyanMedia/burlesque/client"
	"github.com/codegangsta/cli"
)

func main() {
	var bsq *client.Client

	app := cli.NewApp()
	app.Name = "Burlesque API client"
	app.Usage = "Usage details here"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "127.0.0.1",
			Usage: "Burlesque server host",
		},
		cli.IntFlag{
			Name:  "port",
			Value: 4401,
			Usage: "Burlesque server port",
		},
		cli.IntFlag{
			Name:  "timeout",
			Value: 60,
			Usage: "Subscription timeout in seconds",
		},
	}
	app.Before = func(c *cli.Context) (err error) {
		cfg := &client.Config{
			Host:    c.String("host"),
			Port:    c.Int("port"),
			Timeout: time.Duration(c.Int("timeout")) * time.Second,
		}
		bsq = client.NewClient(cfg)
		return
	}
	app.Commands = []cli.Command{
		{
			Name:  "pub",
			Usage: "Publish a message to a queue",
			Action: func(c *cli.Context) {
				var msg *client.Message

				switch len(c.Args()) {
				case 2:
					msg = &client.Message{
						Queue: c.Args()[0],
						Body:  []byte(c.Args()[1]),
					}
				case 1:
					msg = &client.Message{
						Queue: c.Args()[0],
					}
					var err error
					msg.Body, err = ioutil.ReadAll(os.Stdin)
					if err != nil {
						panic(err)
					}
				case 0:
					fmt.Println(app.Usage)
					return
				}

				ok := bsq.Publish(msg)
				if ok {
					fmt.Printf("Message successfully published to queue %q\n", msg.Queue)
				} else {
					fmt.Printf("Failed to publish message to queue %q\n", msg.Queue)
				}
			},
		},
		{
			Name:  "sub",
			Usage: "Subscribe for message from queue",
			Action: func(c *cli.Context) {
				msg := bsq.Subscribe(c.Args()...)

				if msg != nil {
					fmt.Println(string(msg.Body))
				} else {
					fmt.Printf("Failed to recieve message from queues %s\n", strings.Join(c.Args(), ", "))
				}
			},
		},
		{
			Name:  "status",
			Usage: "Show server status",
			Action: func(c *cli.Context) {
				stat := bsq.Status()

				for _, queue := range stat {
					fmt.Println(queue.Name)
					fmt.Println("    Messages:", queue.Messages)
					fmt.Println("    Subscribers:", queue.Subscribers)
				}
			},
		},
		{
			Name:  "debug",
			Usage: "Show server debug info",
			Action: func(c *cli.Context) {
				info := bsq.Status()
				fmt.Println(info)
			},
		},
	}
	app.CommandNotFound = func(c *cli.Context, cmd string) {

	}

	app.Run(os.Args)
}
