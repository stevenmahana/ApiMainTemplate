package main

import (
	"github.com/stevenmahana/ApiMainTemplate/src/controllers"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"fmt"
)

/*
	The MUX uses generic routes that allow us to add new services without changing the interface layer.
	It also allows us to server static content as needed.

	URL: /<service>/<object>/<method>
	Version: Version is handled by the service. The service will create it's own internal version. Default = V1
	Object: Connects to corresponding micro service which is mapped to database object
	Method: This tells the service which function to run
	Params: <method?key=value> URL params can be added to the method to provide additional context to query

 */
func main() {

	router := httprouter.New()
	ctlr := controllers.NewController()

	// public routes
	router.GET("/", ctlr.Index)

	// secure routes
	router.GET("/service/:object/:method", ctlr.GetController)
	router.POST("/service/:object/:method", ctlr.CreateController)
	router.PUT("/service/:object/:method", ctlr.UpdateController)
	router.DELETE("/service/:object/:method", ctlr.RemoveController)

	// file or binary upload. requires POST method and object, object uuid
	router.POST("/upload/:object/:uuid", ctlr.UploadController)

	log.Print("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))



}