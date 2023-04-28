package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ResponseBody struct {
	FlagUrl string
	Name    string
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := http.Get("https://flagcdn.com/en/codes.json")
	if err != nil {
		fmt.Print(err.Error())
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result map[string]string
	err = json.Unmarshal([]byte(responseData), &result)
	if err != nil {
		log.Fatal(err)
	}
	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	randomIndex := rand.Intn(len(keys))
	pick := keys[randomIndex]

	responseBody := &ResponseBody{
		FlagUrl: "https://flagcdn.com/h240/" + pick + ".jpg",
		Name:    result[pick],
	}
	responseJson, err := json.Marshal(responseBody)
	if err != nil {
		log.Fatal(err)
	}
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(responseJson),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
