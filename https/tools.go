package https

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func LogAsJson(objectToLog interface{}) {
	byteObjectToLog, _ := json.Marshal(objectToLog)
	log.Println(string(byteObjectToLog))
}

type requestHandler func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func HandleCors(request events.APIGatewayProxyRequest, h requestHandler) (response events.APIGatewayProxyResponse, err error) {
	// traitement de la requete
	response, err = h(request)

	// décoration de la réponse
	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}
	response.Headers[`Access-Control-Allow-Origin`] = getOrigin(request)

	return response, err
}

func getOrigin(e events.APIGatewayProxyRequest) string {
	if origin, ok := e.Headers["origin"]; ok {
		return origin
	}
	if origin, ok := e.Headers["Origin"]; ok {
		return origin
	}
	return ""
}
