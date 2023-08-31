package googledrive

import (
	"avito-internship/internal/apperror"
	"bytes"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GDriveWebAPI struct {
	driveService *drive.Service
	isAvailable  bool
}

func New(apiJSONFilePath string) *GDriveWebAPI {
	if apiJSONFilePath == "" {
		return &GDriveWebAPI{isAvailable: false}
	}

	driveService, err := drive.NewService(context.Background(), option.WithCredentialsFile(apiJSONFilePath))
	if err != nil {
		panic(err)
	}

	return &GDriveWebAPI{
		driveService: driveService,
		isAvailable:  true,
	}
}

func (w *GDriveWebAPI) IsAvailable() bool {
	return w.isAvailable
}

func (w *GDriveWebAPI) UploadCSVFile(ctx context.Context, name string, data []byte) (string, error) {
	fileId, err := w.getFileIdByName(ctx, name)
	if err != nil {
		if !errors.Is(err, apperror.ErrFileNotFound) {
			return "", err
		}

		id, err := w.createFile(ctx, name, data)
		if err != nil {
			return "", err
		}

		return w.getFileURL(id), nil
	}

	err = w.updateFile(ctx, fileId, data, name)
	if err != nil {
		return "", err
	}

	return w.getFileURL(fileId), nil
}

func (w *GDriveWebAPI) DeleteFile(ctx context.Context, name string) error {
	fileId, err := w.getFileIdByName(ctx, name)
	if err != nil {
		return err
	}

	err = w.driveService.Files.Delete(fileId).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("GDriveWebAPI.DeleteFile - w.driveService.Files.Delete: %w", err)
	}

	return nil
}

func (w *GDriveWebAPI) GetAllFilenames(ctx context.Context) ([]string, error) {
	files, err := w.getAllFiles(ctx)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, file.Name)
	}

	return names, nil
}

func (w *GDriveWebAPI) createFile(ctx context.Context, name string, content []byte) (string, error) {
	file := &drive.File{
		Name:     name,
		MimeType: "text/csv",
	}

	permissions := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}

	_, err := w.driveService.Files.Create(file).Context(ctx).Media(bytes.NewReader(content)).Do()
	if err != nil {
		return "", fmt.Errorf("GDriveWebAPI.createFile - w.driveService.Files.Create: %w", err)
	}

	fileId, err := w.getFileIdByName(ctx, name)
	if err != nil {
		return "", err
	}

	_, err = w.driveService.Permissions.Create(fileId, permissions).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("GDriveWebAPI.createFile - w.driveService.Permissions.Create: %w", err)
	}

	return fileId, nil
}

func (w *GDriveWebAPI) updateFile(ctx context.Context, id string, content []byte, name string) error {
	file := &drive.File{
		Name:     name,
		MimeType: "text/csv",
	}

	_, err := w.driveService.Files.Update(id, file).Context(ctx).Media(bytes.NewReader(content)).Do()
	if err != nil {
		return fmt.Errorf("GDriveWebAPI.updateFile - w.driveService.Files.Update: %w", err)
	}

	return nil
}

func (w *GDriveWebAPI) getFileURL(id string) string {
	return fmt.Sprintf("https://drive.google.com/file/d/%s/view?usp=sharing", id)
}

func (w *GDriveWebAPI) getAllFiles(ctx context.Context) ([]*drive.File, error) {
	r, err := w.driveService.Files.List().Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("GDriveWebAPI.getAllFiles - w.driveService.Files.List: %w", err)
	}

	return r.Files, nil
}

func (w *GDriveWebAPI) getFileIdByName(ctx context.Context, name string) (string, error) {
	files, err := w.getAllFiles(ctx)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.Name == name {
			return file.Id, nil
		}
	}

	return "", apperror.ErrFileNotFound
}
