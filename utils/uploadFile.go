package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/guatom999/BBBot/config"
)

func UploadFile(cfg *config.Config, client *storage.Client, pctx context.Context, data []byte) (string, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*50)
	defer cancel()

	buff := bytes.NewBuffer(data)

	destination := fmt.Sprintf(cfg.Gcp.UploadPath + fmt.Sprintf("donate_%v", time.Now().UnixMilli()))

	wc := client.Bucket(cfg.Gcp.BucketName).Object(destination).NewWriter(ctx)
	if _, err := io.Copy(wc, buff); err != nil {
		fmt.Printf("Error:Failed to Upload File io.Copy: %s", err.Error())
		return "", errors.New("faile to used io.copy")
	}
	if err := wc.Close(); err != nil {
		fmt.Printf("Error:Failed to Upload File wc.Close: %s", err.Error())
		return "", errors.New("failed to closed writer")
	}

	if err := makePublic(cfg, ctx, client, destination); err != nil {
		fmt.Printf("Error:Faile to Make File Public: %s", err.Error())
		return "", errors.New("error: failed to make file public")
	}

	urlFile := fmt.Sprintf("https://storage.googleapis.com/%s/%s", cfg.Gcp.BucketName, destination)

	return urlFile, nil

}

func makePublic(cfg *config.Config, ctx context.Context, client *storage.Client, destination string) error {

	acl := client.Bucket(cfg.Gcp.BucketName).Object(destination).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %w", err)
	}
	// fmt.Printf("Blob %v is now publicly accessible.\n", destination)
	return nil
}
