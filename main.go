package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	cmd "gitlab.eng.vmware.com/nidhig1/goassignment-postandcomments/cmd"
)

func main() {
	//calling the cmd package
	fmt.Println("Starting the application")
	if err := cmd.Run(os.Args, os.Stdout); err != nil {
		switch err {
		case context.Canceled:
			fmt.Println("Parent Context Shutting")
		case http.ErrServerClosed:
			fmt.Println("Closing the server")
		default:
			log.Fatalf("could not run application: %v", err)
		}
	}
}

//main should be as small as possible-folder structure---------done
// function and type to seperate files-folder structure--------done
//periodic reconciler-ticker based fn which for every 10 sec reads post and comments.---done
//add loggers: meaningfull---done
//follow server template given by dilli---done
//update gracefull shutdown http server.-context
// do not use panic-- write error to http writer if big in loop then log the rrror and  continue

/////// improve function for fetching one post and comments--pass id while reading post and comments seperates
