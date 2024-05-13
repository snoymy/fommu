package adapter

import "context"

type FileAdapter interface {
    WriteFile(ctx context.Context, file []byte, fileName string) (string, error)
    ReadFile(ctx context.Context, fileUrl string) ([]byte, error)
}
