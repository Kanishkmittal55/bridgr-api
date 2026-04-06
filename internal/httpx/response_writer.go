package httpx

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/go-chi/render"

	"github.com/Kanishkmittal55/bridgr-api/internal/ctxlog"
	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/errors"
	"github.com/Kanishkmittal55/bridgr-api/pkg/types"
)

const (
	HeaderKeyContentType                  = "Content-Type"
	HeaderValueContentTypeApplicationJson = "application/json"
)

type stubFilePathKey struct{}

// SetStubbedFilePath sets the output path for CaptureResponse (optional test hook).
func SetStubbedFilePath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, stubFilePathKey{}, path)
}

func stubFilePath(ctx context.Context) string {
	p, _ := ctx.Value(stubFilePathKey{}).(string)
	return p
}

type CapturedRequestResponse struct {
	Path        string      `json:"Path"`
	QueryParams url.Values  `json:"QueryParams"`
	Status      int         `json:"Status"`
	Response    interface{} `json:"Response"`
}

func WriteResponse(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	render.Status(r, status)
	render.Respond(w, r, v)
}

type FailureMetricsRecorder func(ctx context.Context, status int, err error)

type ResponseWriter struct {
	logger                 *ctxlog.ContextualLogger
	failureMetricsRecorder FailureMetricsRecorder
	shouldCaptureResponse  bool
}

func NewResponseWriter(logger *ctxlog.ContextualLogger) *ResponseWriter {
	return NewResponseWriterWithCapture(logger, false)
}

func NewResponseWriterWithCapture(logger *ctxlog.ContextualLogger, captureResponse bool) *ResponseWriter {
	return &ResponseWriter{
		logger:                logger,
		shouldCaptureResponse: captureResponse,
	}
}

func (rw *ResponseWriter) SetFailureMetricsRecorder(recorder FailureMetricsRecorder) {
	rw.failureMetricsRecorder = recorder
}

func (rw *ResponseWriter) WriteResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, i interface{}, err error) {
	rw.writeResponse(ctx, w, r, i, err, http.StatusOK)
}

func (rw *ResponseWriter) WriteOkResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, i interface{}, err error) {
	rw.writeResponse(ctx, w, r, i, err, http.StatusOK)
}

func (rw *ResponseWriter) WriteAcceptedResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, i interface{}, err error) {
	rw.writeResponse(ctx, w, r, i, err, http.StatusAccepted)
}

func (rw *ResponseWriter) WriteNoContentResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	rw.writeResponse(ctx, w, r, nil, err, http.StatusNoContent)
}

func (rw *ResponseWriter) writeResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, i interface{}, err error, successStatus int) {
	if err == nil {
		WriteResponse(w, r, successStatus, i)
		return
	}

	var httpStatus int
	switch {
	case stderrors.Is(err, apierrors.ErrBadRequest):
		httpStatus = http.StatusBadRequest
		WriteResponse(w, r, httpStatus, types.NewErrorResponse(err.Error()))
	case stderrors.Is(err, apierrors.ErrNotFound):
		httpStatus = http.StatusNotFound
		WriteResponse(w, r, httpStatus, types.NewNotFoundErrorResponse(err.Error()))
	case stderrors.Is(err, apierrors.ErrForbidden):
		httpStatus = http.StatusForbidden
		WriteResponse(w, r, httpStatus, types.NewErrorResponse(err.Error()))
	case stderrors.Is(err, apierrors.ErrTooManyRequests):
		httpStatus = http.StatusTooManyRequests
		WriteResponse(w, r, httpStatus, types.NewErrorResponse(err.Error()))
	case stderrors.Is(err, apierrors.ErrInternal):
		httpStatus = http.StatusInternalServerError
		rw.getLogger(ctx).Errorw("error handling response", "error", err)
		WriteResponse(w, r, httpStatus, types.NewErrorResponse(err.Error()))
	case stderrors.Is(err, apierrors.ErrNotImplemented):
		httpStatus = http.StatusNotImplemented
		WriteResponse(w, r, httpStatus, types.NewErrorResponse(err.Error()))
	case stderrors.Is(err, apierrors.ErrServiceUnavailable):
		httpStatus = http.StatusServiceUnavailable
		WriteResponse(w, r, httpStatus, types.NewErrorResponse(err.Error()))
	default:
		httpStatus = http.StatusInternalServerError
		rw.getLogger(ctx).Errorw("error handling response", "error", err)
		WriteResponse(w, r, httpStatus, types.NewErrorResponse(err.Error()))
	}

	if rw.failureMetricsRecorder != nil {
		rw.failureMetricsRecorder(ctx, httpStatus, err)
	}
}

func (rw *ResponseWriter) CaptureResponse(ctx context.Context, path string, queryParams url.Values, status int, i interface{}) {
	fp := stubFilePath(ctx)
	if fp == "" {
		return
	}
	dirname := filepath.Dir(fp)
	if err := os.MkdirAll(dirname, 0o755); err != nil {
		rw.logger.Errorw("failed to make directories", "dir", dirname, "err", err)
		return
	}
	rsp := CapturedRequestResponse{
		Path:        path,
		QueryParams: queryParams,
		Status:      status,
		Response:    i,
	}
	rspBytes, err := json.MarshalIndent(rsp, "", "\t")
	if err != nil {
		rw.logger.Errorw("failed to marshal captured request response", "reqResp", rsp, "err", err)
		return
	}
	if err := os.WriteFile(fp, rspBytes, 0o644); err != nil {
		rw.logger.Errorw("failed to write capture file", "err", err)
	}
}

func (rw *ResponseWriter) getLogger(ctx context.Context) *ctxlog.ContextualLogger {
	if ctx == nil {
		return rw.logger
	}
	return rw.logger.AddFromCtx(ctx)
}
