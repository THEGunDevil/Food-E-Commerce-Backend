package services

import (
	"log"
	"mime/multipart"
	"net/http"

	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/models"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

var cld *cloudinary.Cloudinary

// InitCloudinary initializes Cloudinary once from the env URL.
func InitCloudinary(cloudURL string) {
	var err error
	cld, err = cloudinary.NewFromURL(cloudURL)
	if err != nil {
		log.Fatalf("Cloudinary init error: %v", err)
	}
}

// UploadToCloudinary is a generic function to upload files to a specific folder.
func UploadToCloudinary(file multipart.File, filename, folder string) (string, string, error) {
	uploadResp, err := cld.Upload.Upload(db.Ctx, file, uploader.UploadParams{
		Folder:   folder,
		PublicID: filename,
	})
	if err != nil {
		return "", "", err
	}

	return uploadResp.SecureURL, uploadResp.PublicID, nil
}

// DeleteImageFromCloudinary removes an asset using its PublicID.
func DeleteImageFromCloudinary(publicID string) error {
	_, err := cld.Upload.Destroy(db.Ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

// HandleImageUpload processes the Gin file header and uploads it to the specified folder.
// This replaces both UploadImage and UploadProfileImage.
func HandleImageUpload(c *gin.Context, fileHeader *multipart.FileHeader, folder string) (string, string, error) {
	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to open image",
		})
		return "", "", err
	}
	defer f.Close()

	imageURL, publicID, err := UploadToCloudinary(f, fileHeader.Filename, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Image upload failed",
		})
		return "", "", err
	}

	return imageURL, publicID, nil
}