package helpers

import (
    "cloud.google.com/go/storage"
    "context"
    "fmt"
    "io"
    "mime/multipart"
    "google.golang.org/api/option"
)

const (
    googleProjectID   = "blogflex-images"
    googleBucketName  = "images-blogs"
    googleCredentials = "./credentials/blogflex-images-d00c068cf344.json" // Update this path
)

func UploadFileToGCS(file multipart.File, fileName string) (string, error) {
    ctx := context.Background()

    client, err := storage.NewClient(ctx, option.WithCredentialsFile(googleCredentials))
    if err != nil {
        return "", fmt.Errorf("storage.NewClient: %v", err)
    }
    defer client.Close()

    wc := client.Bucket(googleBucketName).Object(fileName).NewWriter(ctx)
    if _, err = io.Copy(wc, file); err != nil {
        return "", fmt.Errorf("io.Copy: %v", err)
    }
    if err := wc.Close(); err != nil {
        return "", fmt.Errorf("Writer.Close: %v", err)
    }

    return fmt.Sprintf("https://storage.googleapis.com/%s/%s", googleBucketName, fileName), nil
}
