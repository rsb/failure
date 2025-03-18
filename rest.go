package failure

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Field is used to represent a form field, sturct field or any
// key whose value is invalid
type Field struct {
	Key string
	Msg string
}

func NewField(key, msg string) Field {
	return Field{Key: key, Msg: msg}
}

func IsFieldError(e error) bool {
	var f Field
	return errors.As(e, &f)
}

func (f Field) Error() string {
	return fmt.Sprintf("%s: %s", f.Key, f.Msg)
}

func (f Field) Empty() bool {
	return f.Key == "" && f.Msg == ""
}

// FieldGroup is used to represent a group of Fields like
// a form or a struct
type FieldGroup struct {
	Name   string
	Fields []Field
}

func NewFieldGroup(name string) *FieldGroup {
	return &FieldGroup{
		Name:   name,
		Fields: make([]Field, 0),
	}
}

func (f *FieldGroup) ErrorCount() int {
	return len(f.Fields)
}

func (f *FieldGroup) Error() string {
	if f.ErrorCount() == 0 {
		return ""
	}

	errors := []string{}
	for _, e := range f.Fields {
		errors = append(errors, e.Error())
	}
	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(errors, ", "))
}

func (f *FieldGroup) Add(items ...Field) {
	if f.Fields == nil {
		f.Fields = make([]Field, 0)
	}

	f.Fields = append(f.Fields, items...)
}

func (f *FieldGroup) AddField(name, msg string) {
	if f.Fields == nil {
		f.Fields = make([]Field, 0)
	}

	f.Fields = append(f.Fields, NewField(name, msg))
}

func (f *FieldGroup) HasError(key string) bool {
	for _, field := range f.Fields {
		if field.Key == key {
			return true
		}
	}

	return false
}

func (f *FieldGroup) Message(key string) string {
	for _, field := range f.Fields {
		if field.Key == key {
			return field.Msg
		}
	}

	return ""
}

func (f *FieldGroup) Field(key string) (Field, bool) {
	var field Field
	for _, field := range f.Fields {
		if field.Key == key {
			return field, true
		}
	}
	return field, false
}

func (f *FieldGroup) HasErrors() bool {
	return f.ErrorCount() > 0
}

type Catalog struct {
	Status int
	Key    string
	Groups map[string]*FieldGroup
}

func NewCatalog(key string, opts ...int) *Catalog {
	var status int

	if len(opts) > 0 && opts[0] != 0 {
		status = opts[0]
	} else {
		status = http.StatusUnprocessableEntity
	}

	return &Catalog{
		Status: status,
		Key:    key,
		Groups: make(map[string]*FieldGroup),
	}
}

func (fc *Catalog) FormKey() string {
	return fc.Key
}

func (fc *Catalog) ErrorCount() int {
	count := 0
	for _, g := range fc.Groups {
		count += g.ErrorCount()
	}
	return count
}

func (fc *Catalog) HttpStatus() int {
	return fc.Status
}

func (fc *Catalog) MarkAsBadRequest() {
	fc.Status = http.StatusBadRequest
}

func (fc *Catalog) MarkAsUnprocessableEntity() {
	fc.Status = http.StatusUnprocessableEntity
}

func (fc *Catalog) SetStatus(status int) {
	fc.Status = status
}

func (fc *Catalog) Add(items ...*FieldGroup) {
	for _, group := range items {
		if group == nil {
			continue
		}
		fc.Groups[group.Name] = group
	}
}

func (fc *Catalog) AddNewGroup(name string) *FieldGroup {
	group := NewFieldGroup(name)
	fc.Groups[name] = group
	return group
}

func (fc *Catalog) Field(group, key string) (Field, bool) {
	var field Field
	g, ok := fc.Groups[group]
	if !ok {
		return field, false
	}

	return g.Field(key)
}

func (fc *Catalog) AddField(grp, name, msg string) {
	group, ok := fc.Groups[grp]
	if !ok {
		group = NewFieldGroup(grp)
		fc.Groups[grp] = group
	}
	group.AddField(name, msg)
}

func (fc Catalog) HasErrors() bool {
	for _, g := range fc.Groups {
		if g.HasErrors() {
			return true
		}
	}
	return false
}

func (fc Catalog) AllFailures() map[string]map[string]string {
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

func (fc Catalog) Error() string {
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
