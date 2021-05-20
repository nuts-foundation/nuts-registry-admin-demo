package api

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type LogoHandler struct {
	FilePath    string
	contentType string
	data        []byte
}

func (h *LogoHandler) Handle(context echo.Context) error {
	if len(h.FilePath) == 0 {
		// No logo configured
		context.Response().WriteHeader(http.StatusNotFound)
		return nil
	}
	// We cache the logo file. We don't need to worry about race conditions because if those happen,
	// cached contentType and/or data is just overwritten with the same data. So no need for locking here.
	var err error
	if h.contentType == "" || len(h.data) == 0 {
		h.data, err = os.ReadFile(h.FilePath)
		if err != nil {
			// Unable to read logo file
			logrus.Errorf("Unable to read logo file (path=%s): %v", h.FilePath, err)
			context.Response().WriteHeader(http.StatusNotFound)
			return nil
		}
		h.contentType = http.DetectContentType(h.data)
	}
	context.Response().Header().Set("Content-Type", h.contentType)
	context.Response().WriteHeader(http.StatusOK)
	_, err = context.Response().Write(h.data)
	return err
}
