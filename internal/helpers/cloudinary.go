package helpers

import (
    "context"
    "fmt"
    "mime/multipart"
    "os"

    "github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadFileToCloudinary(file multipart.File, fileName string) (string, error) {
    ctx := context.Background()

    cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
    apiKey := os.Getenv("CLOUDINARY_API_KEY")
    apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

    if cloudName == "" || apiKey == "" || apiSecret == "" {
        return "", fmt.Errorf("Cloudinary credentials are not set")
    }

    cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
    if err != nil {
        return "", fmt.Errorf("cloudinary.NewFromParams: %v", err)
    }

    uploadParams := uploader.UploadParams{
        PublicID: fileName,
    }

    uploadResult, err := cld.Upload.Upload(ctx, file, uploadParams)
    if err != nil {
        return "", fmt.Errorf("cld.Upload.Upload: %v", err)
    }

    return uploadResult.SecureURL, nil
}
