package adapter

import (
	"app/internal/config"
	"context"
	"net/url"
	"os"
	"path"
)

type FileAdapterImpl struct {

}

func NewFileAdapterImpl() *FileAdapterImpl {
    return &FileAdapterImpl{}
}

func (a *FileAdapterImpl) WriteFile(ctx context.Context, file []byte, fileName string) (string, error) {
    filePath := path.Join("./media", fileName)

    outFile, err := os.Create(filePath)
    if err != nil {
        return "", err
    }
    defer outFile.Close()

    _, err = outFile.Write(file)
    if err != nil {
        return "", nil
    }

    fileUrl, err := url.JoinPath(config.Fommu.FileHost, fileName)
    if err != nil {
        return "", err
    }

    return fileUrl, nil
}

func (a *FileAdapterImpl) ReadFile(ctx context.Context, fileName string) ([]byte, error) {
    filePath := path.Join("./media", fileName)
    fileBytes, err := os.ReadFile(filePath)
	if err != nil {
	    return nil, err
	}
    return fileBytes, nil
}
