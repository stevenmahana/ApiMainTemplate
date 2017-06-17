package controllers

import (
	"github.com/stevenmahana/ApiMainTemplate/src/models"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/go-nats"
	"encoding/json"
	"net/http"
	"time"
	"log"
	"fmt"
	"io"
)


type (
	// MainController represents the controller for operating on the Service Object
	MainController struct{}
	test_struct struct {}
)

// NewController exposes all of the controller methods
func NewController() *MainController {
	return &MainController{}
}

/*
	** PUBLIC ROUTES **

	This is the index route. Serves static page or message
*/
func (uc MainController) Index(resp http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Fprint(resp, "Hello World!\n")
}


/*
	** SECURE ROUTES **

	This is the GET Controller of all micro services.
	This is a generic controller. It allows dynamic addition of new micro services without changing interface layer.

	URL: /<service>/<object>/<method>
	Version: <method?v=V1.0> The service will create it's own internal version. Default = V1
	Object: Connects to corresponding micro service which is mapped to database object
	Method: This tells the service which function to run
	Params: <method?key=value> URL params can be added to the method to provide additional context to query

 */
func (uc MainController) GetController(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	log.Print("Request to: ", r.URL)

	auth := models.Access()
	// verify header was set correctly and check for required header elements
	if auth.VerifyHeader(r.Header) == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// verify token exists, matches token issued by auth server and is valid
	if auth.VerifyToken(r.Header) == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// verify key exists, matches key in cache
	user, valid := auth.VerifyKey(r.Header)
	if valid == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Get URL Params ?key=value; Route Params are in "p"
	q := r.URL.Query()

	// Build Message Payload
	payload := models.MessagePayload{
		Auid: user.Auid,
		Uuid: q.Get("uuid"),
		Key: q.Get("key"),
		Keyword: q.Get("keyword"),
		Perspective: q.Get("perspective"),
		Body: "",
		Object: p.ByName("object"),
		Method: p.ByName("method"),
		Version: q.Get("v"),
		Results: q.Get("results"),
		Page: q.Get("page"),
		Http_method: "GET",
	}

	// Subject is mapped to micro service
	subject := p.ByName("object")

	// Marshal payload into JSON structure
	mp, _ := json.Marshal(payload)

	// Publish message on subject
	message := string(mp)

	// Connect to NATS server; defer close
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	defer natsConnection.Close()

	log.Println("Connected to " + nats.DefaultURL)

	// Set Response Header
	w.Header().Set("Content-Type", "application/json")

	// Send Message
	msg, err := natsConnection.Request(subject, []byte(message), 3000*time.Millisecond)
	if err != nil {
		log.Println(">>> ERROR: Service Connect Error - ", err)

		// Set HTTP Response Method
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Set HTTP Response Method
	w.WriteHeader(http.StatusOK)

	// Response Object is a JSON String; Response Object is created by the service
	fmt.Fprintf(w, "%s", msg.Data)

}


/*
	This is the UPLOAD (POST) Controller for uploading files or streaming files to S3

	URL: /<upload>/<object>/<uuid>
	Upload: Fixed
	Object: Database object
	Uuid: uuid returned from the server when the object was created.
	Params: <method?key=value> URL params can be added to the method to provide additional context to query
 */
func (uc MainController) UploadController(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	auth := models.Access()
	// verify header was set correctly and check for required header elements
	if auth.VerifyHeader(r.Header) == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// verify token exists, matches token issued by auth server and is valid
	if auth.VerifyToken(r.Header) == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// verify key exists, matches key in cache
	user, valid := auth.VerifyKey(r.Header)
	if valid == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Get URL Params ?key=value; Route Params are in "p"
	q := r.URL.Query()

	// Read body, check valid JSON, process errors and ensure body size isn't larger than 1M
	var jbody interface{}
	err := json.NewDecoder(io.LimitReader(r.Body, 1000000)).Decode(&jbody)
	if err != nil {
		log.Println(">>> ERROR: JSON Decoder error - ", err)

		// Set HTTP Response Method
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	defer r.Body.Close() // close body, can cause memory leaks

	// create json string for message body
	body, err := json.Marshal(jbody)
	if err != nil {
		log.Println(">>> ERROR: JSON Marshal error - ", err)

		// Set HTTP Response Method
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Build Message Payload
	payload := models.MessagePayload{
		Auid: user.Auid,
		Uuid: q.Get("uuid"),
		Key: q.Get("key"),
		Keyword: q.Get("keyword"),
		Perspective: q.Get("perspective"),
		Body: string(body),
		Object: p.ByName("object"),
		Method: p.ByName("method"),
		Version: q.Get("v"),
		Results: q.Get("results"),
		Page: q.Get("page"),
		Http_method: "POST",
	}

	// Subject is mapped to micro service
	subject := p.ByName("object")

	// Marshal payload into JSON structure
	mp, _ := json.Marshal(payload)

	// Publish message on subject
	message := string(mp)

	// Connect to NATS server; defer close
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	defer natsConnection.Close()

	log.Println("Connected to " + nats.DefaultURL)

	// Set Response Header
	w.Header().Set("Content-Type", "application/json")

	// Send Message
	msg, err := natsConnection.Request(subject, []byte(message), 3000*time.Millisecond)
	if err != nil {
		log.Println(">>> ERROR: Service Connect Error - ", err)

		// Set HTTP Response Method
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//fmt.Println(string(msg.Data))

	// Set HTTP Response Method
	w.WriteHeader(http.StatusOK)

	// Response Object is a JSON String; Response Object is created by the service
	fmt.Fprintf(w, "%s", msg.Data)

}


/*
	This is the POST Controller of all micro services.
	This is a generic controller that creates new objects.

	URL: /<service>/<object>/<method>
	Version: <method?v=V1.0> The service will create it's own internal version. Default = V1
	Object: Connects to corresponding micro service which is mapped to database object
	Method: This tells the service which function to run
	Params: <method?key=value> URL params can be added to the method to provide additional context to query

 */
func (uc MainController) CreateController(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	auth := models.Access()
	// verify header was set correctly and check for required header elements
	if auth.VerifyHeader(r.Header) == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// verify token exists, matches token issued by auth server and is valid
	if auth.VerifyToken(r.Header) == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// verify key exists, matches key in cache
	user, valid := auth.VerifyKey(r.Header)
	if valid == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Get URL Params ?key=value; Route Params are in "p"
	q := r.URL.Query()

	// Read body, check valid JSON, process errors and ensure body size isn't larger than 1M
	var jbody interface{}
	err := json.NewDecoder(io.LimitReader(r.Body, 1000000)).Decode(&jbody)
	if err != nil {
		log.Println(">>> ERROR: JSON Decoder error - ", err)

		// Set HTTP Response Method
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	defer r.Body.Close() // close body, can cause memory leaks

	// create json string for message body
	body, err := json.Marshal(jbody)
	if err != nil {
		log.Println(">>> ERROR: JSON Marshal error - ", err)

		// Set HTTP Response Method
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Build Message Payload
	payload := models.MessagePayload{
		Auid: user.Auid,
		Uuid: q.Get("uuid"),
		Key: q.Get("key"),
		Keyword: q.Get("keyword"),
		Perspective: q.Get("perspective"),
		Body: string(body),
		Object: p.ByName("object"),
		Method: p.ByName("method"),
		Version: q.Get("v"),
		Results: q.Get("results"),
		Page: q.Get("page"),
		Http_method: "POST",
	}

	// Subject is mapped to micro service
	subject := p.ByName("object")

	// Marshal payload into JSON structure
	mp, _ := json.Marshal(payload)

	// Publish message on subject
	message := string(mp)

	// Connect to NATS server; defer close
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	defer natsConnection.Close()

	log.Println("Connected to " + nats.DefaultURL)

	// Set Response Header
	w.Header().Set("Content-Type", "application/json")

	// Send Message
	msg, err := natsConnection.Request(subject, []byte(message), 3000*time.Millisecond)
	if err != nil {
		log.Println(">>> ERROR: Service Connect Error - ", err)

		// Set HTTP Response Method
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//fmt.Println(string(msg.Data))

	// Set HTTP Response Method
	w.WriteHeader(http.StatusOK)

	// Response Object is a JSON String; Response Object is created by the service
	fmt.Fprintf(w, "%s", msg.Data)

}


/*
	This is the PUT Controller of all micro services.
	This is a generic controller that updates objects.

	URL: /<service>/<object>/<method>
	Version: <method?v=V1.0> The service will create it's own internal version. Default = V1
	Object: Connects to corresponding micro service which is mapped to database object
	Method: This tells the service which function to run
	Params: <method?key=value> URL params can be added to the method to provide additional context to query

 */
func (uc MainController) UpdateController(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	auth := models.Access()
	// verify header was set correctly and check for required header elements
	if auth.VerifyHeader(r.Header) == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// verify token exists, matches token issued by auth server and is valid
	if auth.VerifyToken(r.Header) == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// verify key exists, matches key in cache
	user, valid := auth.VerifyKey(r.Header)
	if valid == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Get URL Params ?key=value; Route Params are in "p"
	q := r.URL.Query()

	// Read body, check valid JSON, process errors and ensure body size isn't larger than 1M
	var jbody interface{}
	err := json.NewDecoder(io.LimitReader(r.Body, 1000000)).Decode(&jbody)
	if err != nil {
		log.Println(">>> ERROR: JSON Decoder error - ", err)

		// Set HTTP Response Method
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	defer r.Body.Close() // close body, can cause memory leaks

	// create json string for message body
	body, err := json.Marshal(jbody)
	if err != nil {
		log.Println(">>> ERROR: JSON Marshal error - ", err)

		// Set HTTP Response Method
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Build Message Payload
	payload := models.MessagePayload{
		Auid: user.Auid,
		Uuid: q.Get("uuid"),
		Key: q.Get("key"),
		Keyword: q.Get("keyword"),
		Perspective: q.Get("perspective"),
		Body: string(body),
		Object: p.ByName("object"),
		Method: p.ByName("method"),
		Version: q.Get("v"),
		Results: q.Get("results"),
		Page: q.Get("page"),
		Http_method: "PUT",
	}

	// Subject is mapped to micro service
	subject := p.ByName("object")

	// Marshal payload into JSON structure
	mp, _ := json.Marshal(payload)

	// Publish message on subject
	message := string(mp)

	// Connect to NATS server; defer close
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	defer natsConnection.Close()

	log.Println("Connected to " + nats.DefaultURL)

	// Set Response Header
	w.Header().Set("Content-Type", "application/json")

	// Send Message
	msg, err := natsConnection.Request(subject, []byte(message), 3000*time.Millisecond)
	if err != nil {
		log.Println(">>> ERROR: Service Connect Error - ", err)

		// Set HTTP Response Method
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//fmt.Println(string(msg.Data))

	// Set HTTP Response Method
	w.WriteHeader(http.StatusOK)

	// Response Object is a JSON String; Response Object is created by the service
	fmt.Fprintf(w, "%s", msg.Data)
}


/*
	This is the DELETE Controller of all micro services.
	This is a generic controller that removes objects.

	URL: /<service>/<object>/<method>
	Version: <method?v=V1.0> The service will create it's own internal version. Default = V1
	Object: Connects to corresponding micro service which is mapped to database object
	Method: This tells the service which function to run
	Params: <method?key=value> URL params can be added to the method to provide additional context to query

 */
func (uc MainController) RemoveController(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	auth := models.Access()
	// verify header was set correctly and check for required header elements
	if auth.VerifyHeader(r.Header) == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// verify token exists, matches token issued by auth server and is valid
	if auth.VerifyToken(r.Header) == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// verify key exists, matches key in cache
	user, valid := auth.VerifyKey(r.Header)
	if valid == false {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Get URL Params ?key=value; Route Params are in "p"
	q := r.URL.Query()

	// Build Message Payload
	payload := models.MessagePayload{
		Auid: user.Auid,
		Uuid: q.Get("uuid"),
		Key: q.Get("key"),
		Keyword: q.Get("keyword"),
		Perspective: q.Get("perspective"),
		Body: "",
		Object: p.ByName("object"),
		Method: p.ByName("method"),
		Version: q.Get("v"),
		Results: q.Get("results"),
		Page: q.Get("page"),
		Http_method: "DELETE",
	}

	// Subject is mapped to micro service
	subject := p.ByName("object")

	// Marshal payload into JSON structure
	mp, _ := json.Marshal(payload)

	// Publish message on subject
	message := string(mp)

	// Connect to NATS server; defer close
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	defer natsConnection.Close()

	log.Println("Connected to " + nats.DefaultURL)

	// Set Response Header
	w.Header().Set("Content-Type", "application/json")

	// Send Message
	msg, err := natsConnection.Request(subject, []byte(message), 3000*time.Millisecond)
	if err != nil {
		log.Println(">>> ERROR: Service Connect Error - ", err)

		// Set HTTP Response Method
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//fmt.Println(string(msg.Data))

	// Set HTTP Response Method
	w.WriteHeader(http.StatusOK)

	// Response Object is a JSON String; Response Object is created by the service
	fmt.Fprintf(w, "%s", msg.Data)
}