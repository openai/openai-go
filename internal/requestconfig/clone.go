package requestconfig

import (
	"context"
	"errors"
	"net/http"
)

func (cfg *RequestConfig) Clone(ctx context.Context) *RequestConfig {
	clone, _ := cfg.CloneWithError(ctx)
	return clone
}

// CloneWithError copies request configuration while reporting body-replay
// failures without disclosing request content or caller-supplied errors.
// This function is internal API and may change without notice.
func (cfg *RequestConfig) CloneWithError(ctx context.Context) (*RequestConfig, error) {
	if cfg == nil {
		return nil, errors.New("requestconfig: cannot clone a nil request configuration")
	}
	if cfg.Request == nil {
		return nil, errors.New("requestconfig: cannot clone a nil request")
	}
	if ctx == nil {
		return nil, errors.New("requestconfig: cannot clone a request without a context")
	}
	req := cfg.Request.Clone(ctx)
	if req.Body != nil && req.Body != http.NoBody {
		if req.GetBody == nil {
			return nil, errors.New("requestconfig: cannot clone a non-replayable request body")
		}
		body, err := req.GetBody()
		if err != nil {
			if body != nil {
				_ = body.Close()
			}
			return nil, errors.New("requestconfig: could not recreate request body")
		}
		if body == nil {
			return nil, errors.New("requestconfig: could not recreate request body")
		}
		req.Body = body
	}
	clone := *cfg
	clone.Context = ctx
	clone.Request = req
	clone.Middlewares = append([]middleware(nil), cfg.Middlewares...)
	clone.finalizers = append([]requestFinalizer(nil), cfg.finalizers...)
	clone.authentication = cfg.authentication.cloneAsInherited(&clone)

	return &clone, nil
}
