package github

import "net/http"

// Logger ...
type Logger struct {
	preRequestCallback  func(*http.Request)
	postRequestCallback func(*http.Response, error)
}

// LoggerOption is an option of Logger.
type LoggerOption func(*Logger)

// NewLogger creates a new Logger.
func NewLogger(opts ...LoggerOption) *Logger {
	l := &Logger{}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *Logger) preRequest(r *http.Request) {
	if l.preRequestCallback != nil {
		l.preRequestCallback(r)
	}
}

// LoggerPreRequest ...
func LoggerPreRequest(callback func(*http.Request)) LoggerOption {
	return func(l *Logger) {
		l.preRequestCallback = callback
	}
}

func (l *Logger) postRequest(r *http.Response, err error) {
	if l.postRequestCallback != nil {
		l.postRequestCallback(r, err)
	}
}

// LoggerPostRequest ...
func LoggerPostRequest(callback func(*http.Response, error)) LoggerOption {
	return func(l *Logger) {
		l.postRequestCallback = callback
	}
}
