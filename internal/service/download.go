package service

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// downloadFile загружает URL в папку downloads/{taskID}/ и возвращает финальный путь.
// Файл сначала пишется в .tmp, затем переименовывается на успешное завершение.
func downloadFile(ctx context.Context, rawurl, taskID string) (string, error) {
	// создаём запрос с контекстом (чтобы поддержать отмену)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawurl, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{"task": taskID, "url": rawurl}).WithError(err).Error("failed to create http request")
		return "", err
	}

	client := &http.Client{
		Timeout: 0, // rely on ctx for cancellation; can set timeout as configuration
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{"task": taskID, "url": rawurl}).WithError(err).Error("http request failed")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{"task": taskID, "url": rawurl, "status": resp.Status}).Error("non-200 response")
		return "", err
	}

	// корректный выбор имени файла из URL
	u, err := url.Parse(rawurl)
	var fname string
	if err == nil {
		fname = path.Base(u.Path) // use path (URL path)
	}
	if fname == "" || fname == "." || fname == "/" {
		fname = uuid.NewString()
	}

	dir := filepath.Join("downloads", taskID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logrus.WithError(err).WithField("dir", dir).Error("mkdir failed")
		return "", err
	}

	tmpPath := filepath.Join(dir, fname+".tmp")
	out, err := os.Create(tmpPath)
	if err != nil {
		logrus.WithError(err).WithField("tmp", tmpPath).Error("create tmp file failed")
		return "", err
	}

	// копируем тело в файл (при отмене ctx операция прервётся с ошибкой)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		out.Close()
		os.Remove(tmpPath)
		logrus.WithError(err).WithFields(logrus.Fields{"tmp": tmpPath, "task": taskID}).Error("copy failed")
		return "", err
	}

	if err := out.Close(); err != nil {
		os.Remove(tmpPath)
		logrus.WithError(err).WithFields(logrus.Fields{"tmp": tmpPath}).Error("close file failed")
		return "", err
	}

	finalPath := filepath.Join(dir, fname)
	if err := os.Rename(tmpPath, finalPath); err != nil {
		os.Remove(tmpPath)
		logrus.WithError(err).WithFields(logrus.Fields{"tmp": tmpPath, "final": finalPath}).Error("rename tmp to final failed")
		return "", err
	}

	return finalPath, nil
}
