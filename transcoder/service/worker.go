package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/subrotokumar/glitchr/libs/core"
	"gitlab.com/subrotokumar/glitchr/libs/storage"
	"gitlab.com/subrotokumar/glitchr/transcoder/config"
	"gitlab.com/subrotokumar/glitchr/transcoder/ffmpeg"
)

type Service struct {
	cfg     config.Config
	log     *core.Logger
	storage *storage.Storage
	bucket  string
	path    string
}

func New() *Service {
	cfg := config.Config{}
	if err := core.ConfigFromEnv(&cfg); err != nil {
		panic(err)
	}
	log := core.NewLogger(cfg.App.Env, cfg.App.Name, cfg.Log.Level)
	storage := storage.NewStorageProvider(cfg.Aws.Region)
	return &Service{
		cfg:     cfg,
		log:     log,
		storage: storage,
	}
}

func (s *Service) Download(ctx context.Context, destPath string) error {
	s.log.Info("Downloading file", "path", destPath)

	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	out, err := s.storage.Client().GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.cfg.S3.Bucket),
		Key:    aws.String(s.cfg.S3.Key),
	})
	if err != nil {
		return fmt.Errorf("get object failed: %w", err)
	}
	defer out.Body.Close()

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, out.Body)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (s *Service) Transcode(ctx context.Context, inputPath, outputDir string) error {

	s.log.Info("Transcoding media", "input", inputPath, "output", outputDir)

	cmdArgs := ffmpeg.DASH_CMD(inputPath, outputDir)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}
		s.log.Debug("Transcoding media", "output", line)
	}
	return nil
}

func (s *Service) Upload(ctx context.Context, sourceDir string) error {
	s.log.Info("Uploading files from", "dir", sourceDir)

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = s.storage.Client().PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s.cfg.S3.Bucket),
			Key:    aws.String(relPath),
			Body:   file,
		})
		if err != nil {
			return fmt.Errorf("upload file %s: %w", relPath, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walk source dir: %w", err)
	}

	return nil
}

func (s *Service) UpdateVideo(ctx context.Context) error {
	s.log.Info("Updating video metadata in database")
	// Implement the logic to update video metadata in the database
	return nil
}

func (s *Service) Process(ctx context.Context) error {
	workDir := "./tmp/workspace"
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return fmt.Errorf("create work dir: %w", err)
	}

	inputPath := filepath.Join(workDir, "input.mp4")
	outputPath := filepath.Join(workDir, "output")
	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	defer func() {
		if _, err := os.Stat(inputPath); err == nil {
			os.Remove(inputPath)
		}
		if _, err := os.Stat(outputPath); err == nil {
			os.RemoveAll(outputPath)
		}
	}()

	if err := s.Download(ctx, inputPath); err != nil {
		return fmt.Errorf("download video: %w", err)
	}

	if err := s.Transcode(ctx, inputPath, outputPath); err != nil {
		return fmt.Errorf("transcode video: %w", err)
	}
	if err := s.Upload(ctx, outputPath); err != nil {
		return fmt.Errorf("upload files: %w", err)
	}
	// if err := s.UpdateVideo(ctx); err != nil {
	// 	return fmt.Errorf("update video: %w", err)
	// }
	return nil
}

func (s *Service) Run(ctx context.Context) {
	s.log.Info("Transcorder worker started processing")
	if err := s.Process(ctx); err != nil {
		s.log.Error("Error processing video", "error", err)
	} else {
		s.log.Info("Video processing completed successfully")
	}
}
