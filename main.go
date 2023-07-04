package main

import (
	"curd-service/routes"
	"fmt"
)

func main() {

	fmt.Println("curd-service start")
	routes.CreateRouter()
	routes.InitializeRoute()
	routes.ServerStart()
}
