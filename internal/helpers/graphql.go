package helpers

import (
    "bytes"
    "encoding/json"
    "net/http"
	"blogflex/internal/database"
    "fmt"
    "strings"
    

)

func GraphQLQuery(query string, variables map[string]interface{}, headers map[string]string) (map[string]interface{}, error) {
    body := map[string]interface{}{
        "query":     query,
        "variables": variables,
    }
    jsonBody, _ := json.Marshal(body)

    req, err := http.NewRequest("POST", database.HasuraEndpoint, bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    for key, value := range headers {
        req.Header.Set(key, value)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func GraphQLRequest(query string, variables map[string]interface{}) (map[string]interface{}, error) {
    requestBody, err := json.Marshal(database.GraphQLRequest{
        Query:     query,
        Variables: variables,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to marshal GraphQL request: %v", err)
    }

    req, err := http.NewRequest("POST", database.HasuraEndpoint, bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create new HTTP request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-hasura-admin-secret", database.HasuraAdminSecret)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to perform HTTP request: %v", err)
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode GraphQL response: %v", err)
    }

    if errors, ok := result["errors"].([]interface{}); ok {
        var errorMessages []string
        for _, err := range errors {
            errorMap := err.(map[string]interface{})
            errorMessages = append(errorMessages, errorMap["message"].(string))
        }
        return nil, fmt.Errorf("GraphQL errors: %s", strings.Join(errorMessages, "; "))
    }

    return result, nil
}

