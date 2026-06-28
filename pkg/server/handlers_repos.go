package server

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/xevonlive-dev/xevon/pkg/agent/source"
	"go.uber.org/zap"
)

// RepoUploadResponse is returned by POST /api/repos/upload. It mirrors the
// frontend's RepoUploadResponse — `source` is the extracted local path that the
// agent --source flow consumes.
type RepoUploadResponse struct {
	RepoID  string `json:"repo_id"`
	Source  string `json:"source"`
	Message string `json:"message"`
}

// repoUploadExts are the archive types the dashboard "Upload source code"
// control accepts.
var repoUploadExts = []string{".zip", ".tar.gz", ".tgz", ".tar", ".tar.bz2", ".tbz2", ".tar.xz", ".txz"}

// HandleRepoUpload handles POST /api/repos/upload — the self-hosted backend for
// the dashboard "Upload source code" control. It accepts a multipart archive,
// extracts it to a per-upload directory under the agent sessions dir, and
// returns the local path so an agentic scan can use it as --source.
//
// Extraction reuses source.ResolveSource, which includes zip-slip protection.
func (h *Handlers) HandleRepoUpload(c fiber.Ctx) error {
	if h.settings == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "settings not available",
			Code:  fiber.StatusInternalServerError,
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "file field is required",
			Code:  fiber.StatusBadRequest,
		})
	}

	name := filepath.Base(file.Filename)
	lower := strings.ToLower(name)
	supported := false
	for _, ext := range repoUploadExts {
		if strings.HasSuffix(lower, ext) {
			supported = true
			break
		}
	}
	if !supported {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "unsupported archive type; upload a .zip, .tar.gz, .tgz, or .tar",
			Code:  fiber.StatusBadRequest,
		})
	}

	// Per-upload directory under the agent sessions dir so each upload is
	// isolated and the extracted tree survives until the scan consumes it.
	repoID := uuid.NewString()
	uploadDir := filepath.Join(h.settings.Agent.EffectiveSessionsDir(), "uploads", repoID)
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		zap.L().Error("repo upload: mkdir failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "failed to prepare upload directory",
			Code:  fiber.StatusInternalServerError,
		})
	}

	// Persist the uploaded archive next to where it will be extracted.
	archivePath := filepath.Join(uploadDir, name)
	if err := saveMultipartFile(file, archivePath); err != nil {
		zap.L().Error("repo upload: save failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "failed to save uploaded file",
			Code:  fiber.StatusInternalServerError,
		})
	}

	// ResolveSource extracts the archive (with zip-slip protection) and returns
	// the local root directory.
	resolved, err := source.ResolveSource(archivePath, uploadDir)
	if err != nil {
		zap.L().Error("repo upload: extraction failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "failed to extract archive: " + err.Error(),
			Code:  fiber.StatusBadRequest,
		})
	}
	if resolved == nil || resolved.LocalPath == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "archive extracted but produced no source directory",
			Code:  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(RepoUploadResponse{
		RepoID:  repoID,
		Source:  resolved.LocalPath,
		Message: "source uploaded and extracted",
	})
}

// saveMultipartFile streams an uploaded multipart file to dst on disk.
func saveMultipartFile(fh *multipart.FileHeader, dst string) error {
	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, src); err != nil {
		return err
	}
	return out.Close()
}
