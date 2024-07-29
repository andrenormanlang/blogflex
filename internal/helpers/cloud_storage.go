package helpers

import (
    "cloud.google.com/go/storage"
    "context"
    "fmt"
    "io"
    "mime/multipart"
)

const (
    googleProjectID  = "blogflex-images"
    googleBucketName = "images-blogs"
)

// UploadFileToGCS uploads a file to Google Cloud Storage and returns the public URL
func UploadFileToGCS(file multipart.File, fileName string) (string, error) {
    ctx := context.Background()

    client, err := storage.NewClient(ctx)
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
