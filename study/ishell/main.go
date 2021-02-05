package main

import (
	"github.com/abiosoft/ishell"
)

func main() {
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "login",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)

			c.Println("Let's simulate login")

			// prompt for input
			c.Print("Username: ")
			username := c.ReadLine()
			c.Print("Password: ")
			password := c.ReadPassword()

			// do something with username and password
			c.Println("Your inputs were", username, "and", password+".")

		},
		Help: "simulate a login",
	})

	// run shell
	shell.Run()
}
