package fingerprint

import (
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	BaselinePath       = "/__nonexistent_fuzzer_baseline__"
	NumBaselineSamples = 3
)

type Fingerprint struct {
	StatusCode    int
	ContentLength int64
	Headers       map[string]string
}

// BaselineRange only compares status code, content type, and content length 
type BaselineRange struct {
	StatusCodeMin    int
	StatusCodeMax    int
	ContentLengthMin int64
	ContentLengthMax int64
	ContentType     string 
}

func CaptureRange(client *http.Client, baseURL, method string, delayMs int) (*BaselineRange, error) {
	url := strings.TrimSuffix(baseURL, "/") + BaselinePath
	samples := make([]*Fingerprint, 0, NumBaselineSamples)

	for i := 0; i < NumBaselineSamples; i++ {
		fp, err := captureOne(client, url, method)
		if err != nil {
			return nil, err
		}
		samples = append(samples, fp)
		if i < NumBaselineSamples-1 && delayMs > 0 {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
	}
	return buildRange(samples), nil
}

// FromResponse builds a Fingerprint from an http.Response (caller must close body if needed).
func FromResponse(resp *http.Response) *Fingerprint {
	fp := &Fingerprint{
		StatusCode:    resp.StatusCode,
		ContentLength: resp.ContentLength,
		Headers:       copyHeaders(resp.Header),
	}
	return fp
}

// InRange reports whether fp matches the baseline using only status code, content type, and content length.
func (r *BaselineRange) InRange(fp *Fingerprint) bool {
	if r == nil || fp == nil {
		return false
	}
	if fp.StatusCode < r.StatusCodeMin || fp.StatusCode > r.StatusCodeMax {
		return false
	}
	if fp.ContentLength < r.ContentLengthMin || fp.ContentLength > r.ContentLengthMax {
		return false
	}
	if r.ContentType != "" {
		if trimContentType(fp.Headers["Content-Type"]) != r.ContentType {
			return false
		}
	}
	return true
}

func trimContentType(s string) string {
	if idx := strings.Index(s, ";"); idx >= 0 {
		return strings.TrimSpace(s[:idx])
	}
	return strings.TrimSpace(s)
}

// Equal reports whether f and other have the same status, length, and headers.
func (f *Fingerprint) Equal(other *Fingerprint) bool {
	if f == nil || other == nil {
		return false
	}
	if f.StatusCode != other.StatusCode || f.ContentLength != other.ContentLength {
		return false
	}
	for k, v := range f.Headers {
		if other.Headers[k] != v {
			return false
		}
	}
	return true
}

// captureOne performs a single request and returns a Fingerprint with body length resolved.
func captureOne(client *http.Client, url, method string) (*Fingerprint, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fp := FromResponse(resp)
	if fp.ContentLength == -1 {
		fp.ContentLength = int64(len(body))
	}
	return fp, nil
}

func buildRange(samples []*Fingerprint) *BaselineRange {
	if len(samples) == 0 {
		return &BaselineRange{}
	}
	r := &BaselineRange{
		StatusCodeMin:    samples[0].StatusCode,
		StatusCodeMax:    samples[0].StatusCode,
		ContentLengthMin: samples[0].ContentLength,
		ContentLengthMax: samples[0].ContentLength,
		ContentType:      contentTypeSameInAll(samples),
	}
	for _, s := range samples[1:] {
		r.StatusCodeMin = min(r.StatusCodeMin, s.StatusCode)
		r.StatusCodeMax = max(r.StatusCodeMax, s.StatusCode)
		r.ContentLengthMin = min(r.ContentLengthMin, s.ContentLength)
		r.ContentLengthMax = max(r.ContentLengthMax, s.ContentLength)
	}
	return r
}

// contentTypeSameInAll returns the Content-Type (media type only) if identical in all samples; else "".
func contentTypeSameInAll(samples []*Fingerprint) string {
	if len(samples) == 0 {
		return ""
	}
	ct := trimContentType(samples[0].Headers["Content-Type"])
	for _, s := range samples[1:] {
		if trimContentType(s.Headers["Content-Type"]) != ct {
			return ""
		}
	}
	return ct
}

func copyHeaders(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) > 0 {
			out[http.CanonicalHeaderKey(k)] = strings.TrimSpace(v[0])
		}
	}
	return out
}
