package exceptions

import "fmt"

type ValidationException struct {
	Errors map[string][]string `json:"errors"`
}

func (v *ValidationException) Get(k string) []string {
	return v.Errors[k]
}

func (v *ValidationException) Put(k string, m ...string) {
	if v.Errors == nil {
		v.Errors = make(map[string][]string)
	}

	_, ok := v.Errors[k]
	if !ok {
		v.Errors[k] = make([]string, 0)
	}

	v.Errors[k] = append(v.Errors[k], m...)
}

func (v *ValidationException) Size() int {
	return len(v.Errors)
}

func (v *ValidationException) Error() string {

	s := ""
	for k, val := range v.Errors {
		for _, vv := range val {
			s += fmt.Sprintf("%s %s", k, vv)
		}
	}
	return s
}
