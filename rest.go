package failure

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// FieldError is used to represent a form field, sturct field or any
// key whose value is invalid
type FormField struct {
	Key string
	Msg string
}

// NewFormField creates a new FieldError with a key and msessage
// the logError is optional and cdefaults to BadRequest. This is
// wahat would be used in logs
func NewFormField(key, msg string) FormField {
	return FormField{Key: key, Msg: msg}
}

func IsFieldError(e error) bool {
	var f FormField
	return errors.As(e, &f)
}

func (f FormField) Error() string {
	return fmt.Sprintf("%s: %s", f.Key, f.Msg)
}

func (f FormField) Empty() bool {
	return f.Key == "" && f.Msg == ""
}

// FieldErrorGroup is used to represent a group of FieldErrors like
// a form or a struct
type FormFieldGroup struct {
	Name   string
	Fields []FormField
}

func NewFormFieldGroup(name string) *FormFieldGroup {
	return &FormFieldGroup{
		Name:   name,
		Fields: make([]FormField, 0),
	}
}

func (f *FormFieldGroup) ErrorCount() int {
	return len(f.Fields)
}

func (f *FormFieldGroup) Error() string {
	if f.ErrorCount() == 0 {
		return ""
	}

	errors := []string{}
	for _, e := range f.Fields {
		errors = append(errors, e.Error())
	}
	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(errors, ", "))
}

func (f *FormFieldGroup) Add(items ...FormField) {
	if f.Fields == nil {
		f.Fields = make([]FormField, 0)
	}

	f.Fields = append(f.Fields, items...)
}

func (f *FormFieldGroup) AddField(name, msg string) {
	if f.Fields == nil {
		f.Fields = make([]FormField, 0)
	}

	f.Fields = append(f.Fields, NewFormField(name, msg))
}

func (f *FormFieldGroup) HasError(key string) bool {
	for _, field := range f.Fields {
		if field.Key == key {
			return true
		}
	}

	return false
}

func (f *FormFieldGroup) Message(key string) string {
	for _, field := range f.Fields {
		if field.Key == key {
			return field.Msg
		}
	}

	return ""
}

func (f *FormFieldGroup) Field(key string) (FormField, bool) {
	var field FormField
	for _, field := range f.Fields {
		if field.Key == key {
			return field, true
		}
	}
	return field, false
}

func (f *FormFieldGroup) HasErrors() bool {
	return f.ErrorCount() > 0
}

type Form struct {
	Status int
	Key    string
	Groups map[string]*FormFieldGroup
}

func NewForm(key string, opts ...int) *Form {
	var status int

	if len(opts) > 0 && opts[0] != 0 {
		status = opts[0]
	} else {
		status = http.StatusUnprocessableEntity
	}

	return &Form{
		Status: status,
		Key:    key,
		Groups: make(map[string]*FormFieldGroup),
	}
}

func (fc *Form) FormKey() string {
	return fc.Key
}

func (fc *Form) ErrorCount() int {
	count := 0
	for _, g := range fc.Groups {
		count += g.ErrorCount()
	}
	return count
}

func (fc *Form) HttpStatus() int {
	return fc.Status
}

func (fc *Form) MarkAsBadRequest() {
	fc.Status = http.StatusBadRequest
}

func (fc *Form) MarkAsUnprocessableEntity() {
	fc.Status = http.StatusUnprocessableEntity
}

func (fc *Form) SetStatus(status int) {
	fc.Status = status
}

func (fc *Form) Add(items ...*FormFieldGroup) {
	for _, group := range items {
		if group == nil {
			continue
		}
		fc.Groups[group.Name] = group
	}
}

func (fc *Form) AddNewGroup(name string) *FormFieldGroup {
	group := NewFormFieldGroup(name)
	fc.Groups[name] = group
	return group
}

func (fc *Form) Field(group, key string) (FormField, bool) {
	var field FormField
	g, ok := fc.Groups[group]
	if !ok {
		return field, false
	}

	return g.Field(key)
}

func (fc *Form) AddField(grp, name, msg string) {
	group, ok := fc.Groups[grp]
	if !ok {
		group = NewFormFieldGroup(grp)
		fc.Groups[grp] = group
	}
	group.AddField(name, msg)
}

func (fc Form) HasErrors() bool {
	for _, g := range fc.Groups {
		if g.HasErrors() {
			return true
		}
	}
	return false
}

func (fc Form) AllFailures() map[string]map[string]string {
	fails := make(map[string]map[string]string)
	for k, g := range fc.Groups {
		if !g.HasErrors() {
			continue
		}

		fails[k] = make(map[string]string)
		for _, field := range g.Fields {
			fails[k][field.Key] = field.Msg
		}
	}
	return fails
}

func (fc Form) Error() string {
	errors := []string{}
	if fc.ErrorCount() == 0 {
		return ""
	}

	for _, g := range fc.Groups {
		errors = append(errors, g.Error())
	}

	line := fmt.Sprintf("%s: %s", fc.Key, strings.Join(errors, ", "))
	return line
}

type RestAPI struct {
	StatusCode int
	Msg        string
	Fields     map[string]string
	Err        error
}

func (r *RestAPI) Error() string {
	return r.Err.Error()
}

func NewInvalidFields(f map[string]string, msg string, a ...any) *RestAPI {
	r := RestAPI{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        fmt.Sprintf(msg, a...),
		Fields:     f,
		Err:        InvalidAPIFieldsErr,
	}

	return &r
}

func InvalidFields(f map[string]string, msg string, a ...any) error {
	return NewInvalidFields(f, msg, a...)
}

func GetInvalidFields(e error) (map[string]string, bool) {
	var r *RestAPI
	if !errors.As(e, &r) {
		return nil, false
	}

	return r.Fields, true
}

func IsInvalidFields(e error) bool {
	var r *RestAPI

	if errors.As(e, &r) {
		return r.StatusCode == http.StatusUnprocessableEntity
	}

	return false
}

func NewBadRequest(msg string, a ...any) *RestAPI {
	r := RestAPI{
		StatusCode: http.StatusBadRequest,
		Msg:        fmt.Sprintf(msg, a...),
		Err:        BadRequestErr,
	}
	return &r
}

func BadRequest(msg string, a ...any) error {
	return NewBadRequest(msg, a...)
}

func ToBadRequest(e error, msg string, a ...any) error {
	r := RestAPI{
		StatusCode: http.StatusBadRequest,
		Msg:        fmt.Sprintf(msg, a...),
		Err:        e,
	}
	return &r
}

func IsBadRequest(e error) bool {
	var r *RestAPI

	if errors.As(e, &r) {
		return r.StatusCode == http.StatusBadRequest
	}

	return false
}

func RestStatusCode(e error) (int, bool) {
	var r *RestAPI

	if errors.As(e, &r) {
		return r.StatusCode, true
	}

	return 0, false
}

func RestMessage(e error) (string, bool) {
	var r *RestAPI

	if errors.As(e, &r) {
		return r.Msg, true
	}

	return "", false
}

func RestError(e error) (error, bool) {
	var r *RestAPI

	if errors.As(e, &r) {
		return r.Err, true
	}

	return nil, false
}

func IsRestAPI(e error) bool {
	var r *RestAPI

	return errors.As(e, &r)
}
