package main

import (
	"crypto/ecdh"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	httpErrorHandler struct {
		statusCodes map[error]int
	}
)

// NewHttpErrorHandler is
func NewHttpErrorHandler(errorStatusCodeMaps map[error]int) *httpErrorHandler {
	return &httpErrorHandler{
		statusCodes: errorStatusCodeMaps,
	}
}

// getStatusCode is
func (self *httpErrorHandler) getStatusCode(err error) int {
	for key, value := range self.statusCodes {
		if errors.Is(err, key) {
			return value
		}
	}

	return http.StatusInternalServerError
}

// unwrapRecursive is
func unwrapRecursive(err error) error {
	var originalErr = err

	for originalErr != nil {
		var internalErr = errors.Unwrap(originalErr)

		if internalErr == nil {
			break
		}

		originalErr = internalErr
	}

	return originalErr
}

// Handler is
func (self *httpErrorHandler) Handler(err error, c echo.Context) {
	he, ok := err.(*echo.HTTPError)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
	} else {
		he = &echo.HTTPError{
			Code: self.getStatusCode(err),
			Message: unwrapRecursive(err).Error(),
		}
	}

	code := he.Code
	message := he.Message
	if _, ok := he.Message.(string); ok {
		message = map[string]interface{}{"message": err.Error()}
	}

	// Send Response 
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(code, message)
		}

		if err != nil {
			c.Echo().Logger.Error(err)
		}
	}
}