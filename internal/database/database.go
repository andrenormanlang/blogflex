package database

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "os"
    "log"
    "strings"

    "github.com/joho/godotenv"
)

var HasuraEndpoint string
var HasuraAdminSecret string

type GraphQLRequest struct {
    Query     string                 `json:"query"`
    Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
    Data   json.RawMessage `json:"data"`
    Errors []GraphQLError  `json:"errors"`
}

type GraphQLError struct {
    Message string `json:"message"`
}

// Initialize the Hasura endpoint from environment variables
func InitHasura() {
    // Load .env file
    err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}


    HasuraEndpoint = os.Getenv("HASURA_ENDPOINT")
    if HasuraEndpoint == "" {
        log.Fatalf("HASURA_ENDPOINT environment variable is not set")
    }

    HasuraAdminSecret = os.Getenv("HASURA_ADMIN_SECRET")
    if HasuraAdminSecret == "" {
        log.Fatalf("HASURA_ADMIN_SECRET environment variable is not set")
    }
}

func SendGraphQLRequest(query string, variables map[string]interface{}) (map[string]interface{}, error) {
    requestBody, err := json.Marshal(GraphQLRequest{Query: query, Variables: variables})
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", HasuraEndpoint, bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-hasura-admin-secret", HasuraAdminSecret)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var response GraphQLResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, err
    }

    if len(response.Errors) > 0 {
        return nil, errors.New(response.Errors[0].Message)
    }

    var data map[string]interface{}
    if err := json.Unmarshal(response.Data, &data); err != nil {
        return nil, err
    }

    return data, nil
}


func ExecuteGraphQL(query string, variables map[string]interface{}) (map[string]interface{}, error) {
    requestBody, err := json.Marshal(map[string]interface{}{
        "query":     query,
        "variables": variables,
    })
    if err != nil {
        log.Printf("Failed to marshal request body: %v", err)
        return nil, err
    }

    req, err := http.NewRequest("POST", HasuraEndpoint, bytes.NewBuffer(requestBody))
    if err != nil {
        log.Printf("Failed to create new HTTP request: %v", err)
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-hasura-admin-secret", HasuraAdminSecret)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("HTTP request failed: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        log.Printf("Failed to decode response body: %v", err)
        return nil, err
    }

    if errs, ok := result["errors"].([]interface{}); ok {
        var errorMessages []string
        for _, err := range errs {
            errorMap := err.(map[string]interface{})
            errorMessages = append(errorMessages, errorMap["message"].(string))
        }
        log.Printf("GraphQL errors: %v", strings.Join(errorMessages, "; "))
        return nil, errors.New(strings.Join(errorMessages, "; "))
    }

    data, ok := result["data"].(map[string]interface{})
    if !ok {
        log.Printf("Response does not contain 'data': %v", result)
        return nil, errors.New("response does not contain 'data'")
    }

    return data, nil
}
