package github

import "fmt"

type apiError struct {
	Resource string `json:"resource"`
	Code     string `json:"code"`
	Field    string `json:"field"`
	Value    string `json:"value"`
}

func (e *apiError) Error() string {
	if e.Value == "" {
		return fmt.Sprintf("%s (%s.%s)", e.Code, e.Resource, e.Field)
	}
	return fmt.Sprintf("%s (%s.%s = %q)", e.Code, e.Resource, e.Field, e.Value)
}

type apiErrors []apiError

func (es apiErrors) Error() string {
	var s string
	for i, e := range es {
		if i > 0 {
			s += ", "
		}
		s += e.Error()
	}
	return s
}
