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
        log.Fatalf("Error loading .env file: %v", err)
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

func SendGraphQLRequest(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
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
        return &response, errors.New(response.Errors[0].Message)
    }

    return &response, nil
}

func ExecuteGraphQL(query string, variables map[string]interface{}) (map[string]interface{}, error) {
    requestBody, err := json.Marshal(map[string]interface{}{
        "query":     query,
        "variables": variables,
    })
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

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    if errs, ok := result["errors"].([]interface{}); ok {
        var errorMessages []string
        for _, err := range errs {
            errorMap := err.(map[string]interface{})
            errorMessages = append(errorMessages, errorMap["message"].(string))
        }
        return nil, errors.New(strings.Join(errorMessages, "; "))
    }

    return result["data"].(map[string]interface{}), nil
}

