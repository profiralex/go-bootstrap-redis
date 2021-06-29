package server

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

type response struct {
	Data   interface{} `json:"data"`
	Errors []apiError  `json:"errors"`
	Status int         `json:"status"`
}

func (response *response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, response.Status)
	if response.Status == 500 {
		for _, e := range response.Errors {
			log.Warnf("API ERROR: %s %s %s", e.Reference, e.Field, e.Message)
		}
	}
	return nil
}

type apiError struct {
	Message   string `json:"message"`
	Field     string `json:"field"`
	Reference string `json:"ref"`
}

func createSuccessResponse(data interface{}, status ...int) response {
	finalStatus := http.StatusOK
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return response{
		Data:   data,
		Errors: nil,
		Status: finalStatus,
	}
}

func createAPIErrorsResponse(errors []apiError, status ...int) response {
	finalStatus := http.StatusInternalServerError
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return response{
		Data:   nil,
		Errors: errors,
		Status: finalStatus,
	}
}

func createAPIErrorResponse(err apiError, status ...int) response {
	finalStatus := http.StatusInternalServerError
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return response{
		Data:   nil,
		Errors: []apiError{err},
		Status: finalStatus,
	}
}

func createErrorResponse(err error, status ...int) response {
	finalStatus := http.StatusInternalServerError
	if len(status) > 0 {
		finalStatus = status[0]
	}

	return response{
		Data:   nil,
		Errors: []apiError{{Message: err.Error()}},
		Status: finalStatus,
	}
}

func respondSuccess(w http.ResponseWriter, r *http.Request, data interface{}, status ...int) {
	rsp := createSuccessResponse(data, status...)
	_ = render.Render(w, r, &rsp)
}

func respondAPIError(w http.ResponseWriter, r *http.Request, err apiError, status ...int) {
	rsp := createAPIErrorResponse(err, status...)
	_ = render.Render(w, r, &rsp)
}

func respondError(w http.ResponseWriter, r *http.Request, err error, status ...int) {
	rsp := createErrorResponse(err, status...)
	_ = render.Render(w, r, &rsp)
}

func respondValidationErrors(w http.ResponseWriter, r *http.Request, err error, status ...int) {
	validationErrors, ok := err.(govalidator.Errors)
	if !ok {
		respondError(w, r, err, status...)
		return
	}

	var apiErrors []apiError
	for _, fieldErrors := range validationErrors {
		fieldValidationErrors, ok := fieldErrors.(govalidator.Errors)
		if !ok {
			fieldValidationError, ok := fieldErrors.(govalidator.Error)
			if ok {
				apiError := apiError{
					Field:   fieldValidationError.Name,
					Message: fieldValidationError.Err.Error(),
				}
				apiErrors = append(apiErrors, apiError)
				continue
			}

			apiErrors = append(apiErrors, apiError{Message: fieldErrors.Error()})
			continue
		}

		for _, fieldError := range fieldValidationErrors {
			fieldValidationError, ok := fieldError.(govalidator.Error)
			if !ok {
				apiErrors = append(apiErrors, apiError{Message: fieldError.Error()})
				continue
			}

			apiError := apiError{
				Field:   fieldValidationError.Name,
				Message: fieldValidationError.Err.Error(),
			}
			apiErrors = append(apiErrors, apiError)
		}
	}

	rsp := createAPIErrorsResponse(apiErrors, status...)
	_ = render.Render(w, r, &rsp)
}
