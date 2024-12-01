// internal/service/file_manager.go

package service

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"github.com/bodgit/sevenzip"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FileManager defines methods for file operations
type FileManager interface {
	RunTreeAtPath(path string) (string, error)
	CreateDirectory(path string, perm os.FileMode) error
	CreateProject(userID, projectName string) error
	FormulatePath(userID, name string, isProjectPath bool) string
	SaveUploadedFile(file *multipart.FileHeader, dst string) error
	RemovePath(path string) error
	RunTree(userID string, isProject bool, projectName string) (string, error)
	ProcessFilesInDirectory(path string, processFile func(relPath string, content []byte) error) error
	ExtractTarGz(src, dst string) error
	ExtractZip(src, dst string) error
	Extract7z(src, dst string) error
}

// OSFileManager is a concrete implementation of FileManager
type OSFileManager struct {
	basePath string
}

const basePath = "/mnt/hdd/projects/evraz_hack"

func NewOSFileManager() *OSFileManager {
	return &OSFileManager{basePath: basePath}
}

func (o *OSFileManager) RunTreeAtPath(path string) (string, error) {
	return o.RunCommandInDirectory(path, "tree")
}

// CreateDirectory creates a new directory with the given permissions
func (o *OSFileManager) CreateDirectory(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// CreateProject creates a new project directory
func (o *OSFileManager) CreateProject(userID, projectName string) error {
	isProject := false
	path := o.FormulatePath(userID, projectName, isProject)
	fmt.Printf("Creating directory at path: %s\n", path)
	return o.CreateDirectory(path, os.ModePerm)
}

// FormulatePath formulates a path based on userID and name
func (o *OSFileManager) FormulatePath(userID, name string, isProjectPath bool) string {
	if isProjectPath {
		return fmt.Sprintf("%s/projects/%s/%s", o.basePath, userID, name)
	} else {
		return fmt.Sprintf("%s/projects/%s/%s", o.basePath, userID, name)
	}
}

// SaveUploadedFile saves the uploaded file to the specified destination path
func (o *OSFileManager) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// RemovePath removes a file or directory and its contents
func (o *OSFileManager) RemovePath(path string) error {
	return os.RemoveAll(path)
}

// RunTree runs the "tree" command for the given project
func (o *OSFileManager) RunTree(userID string, isProject bool, projectName string) (string, error) {
	projectPath := ""
	if isProject {
		projectPath = filepath.Join(o.basePath, "projects", userID, projectName)
	} else {
		projectPath = filepath.Join(o.basePath, "projects", userID, projectName)
	}
	return o.RunCommandInDirectory(projectPath, "tree")
}

// RunCommandInDirectory runs a specified command in the given directory and returns the output
func (o *OSFileManager) RunCommandInDirectory(dir, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	return string(output), err
}

// ProcessFilesInDirectory walks through a directory and processes each file using the provided callback function
func (o *OSFileManager) ProcessFilesInDirectory(path string, processFile func(relPath string, content []byte) error) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath := strings.TrimPrefix(filePath, path+"/")
			content, err := o.ReadFile(filePath)
			if err != nil {
				return err
			}
			return processFile(relPath, content)
		}
		return nil
	})
}

// ExtractTarGz extracts a tar.gz file to the specified destination directory
func (o *OSFileManager) ExtractTarGz(src, dst string) error {
	tarFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	gzr, err := gzip.NewReader(tarFile)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tarReader := tar.NewReader(gzr)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := o.CreateDirectory(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}

func (o *OSFileManager) ExtractZip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dst, f.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			// Create directories as needed
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create directories for the file
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Open the destination file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		// Copy the file contents
		_, err = io.Copy(outFile, rc)

		// Close the file handles
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func (o *OSFileManager) Extract7z(src, dst string) error {
	archive, err := sevenzip.OpenReader(src)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, file := range archive.File {
		fpath := filepath.Join(dst, file.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if file.FileInfo().IsDir() {
			// Create directories as needed
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create directories for the file
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Open the destination file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := file.Open()
		if err != nil {
			return err
		}

		// Copy the file contents
		_, err = io.Copy(outFile, rc)

		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// ReadFile reads data from a file
func (o *OSFileManager) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
