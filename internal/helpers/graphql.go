package helpers

import (
    "bytes"
    "encoding/json"
    "net/http"
	"blogflex/internal/database"
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
