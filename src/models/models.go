package models

// This struct standardizes / normalizes the message payload
type MessagePayload struct {
	Auid 		string `json:"auid"`	// UUID of person making request (authorized UUID)
	Uuid 		string `json:"uuid"`	// UUID of object were referring to
	Object 		string `json:"object"` 	// object type were referring to. Usually mapped to the service
	Key 		string `json:"key"`		// key used for search
	Keyword 	string `json:"keyword"`	// value used for search
	Body 		string `json:"body"`	// Message Body key / value map {key:value, key:value}
	Perspective 	string `json:"perspective"`	// perspective were making the query from. ex: admin, superadmin, etc
	Results 	string `json:"results"`	// qty of objects returned
	Page 		string `json:"page"`	// page were results start
	Http_method 	string `json:"http_method"`  // GET, POST, PUT, DELETE - tell service what request method was used
	Method 		string `json:"method"` 	// Methods are the function that will process the request
	Version 	string `json:"version"` // Version of service requested
}