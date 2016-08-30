package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/cldmnky/gosa/api"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

func main() {
	var hostname string
	app := cli.NewApp()
	app.Name = "run a command against a salt-api endpoint"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host, H",
			Value:       "https://localhost",
			Usage:       "Salt-Api endpoint",
			Destination: &hostname,
		},
	}
	eauth := "ldap"
	port := "443"
	app.Action = func(c *cli.Context) error {
		if c.NArg() != 3 {
			return errors.New("Please supply target, function and args")
		}
		fmt.Print("Username: ")
		var username string
		fmt.Scanln(&username)
		fmt.Print("Password: ")
		password, err := gopass.GetPasswd()
		client := api.NewSaltClient(hostname, port)
		token, err := client.Login(string(username), string(password), eauth)
		if err != nil {
			log.Fatal(err)
		}
		client.SetToken(token)
		fmt.Printf("Running %v against %v\n", c.Args()[1], c.Args()[1])
		res, err := client.Run(c.Args()[0], c.Args()[1], c.Args()[2])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
		return nil
	}
	app.Run(os.Args)
}
