package sulat

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"

	"golang.org/x/exp/maps"
)

type ValidationErrors []*ValidationError

func (err ValidationErrors) Error() string {
	return fmt.Sprintf("%d validation errors", len(err))
}

func valErrOrNil(errs ValidationErrors) error {
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (err ValidationErrors) MarshalJSON() ([]byte, error) {
	errors := map[string][]string{}
	for _, validationError := range err {
		_, ok := errors[validationError.Field]
		if !ok {
			errors[validationError.Field] = []string{}
		}
		errors[validationError.Field] = append(errors[validationError.Field], validationError.Message)
	}
	return json.Marshal(errors)
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Message)
}

type Schema []SchemaField

func (s Schema) Scan(src interface{}) error {
	return scanJson(src, &s, "Schema")
}

func (s Schema) Value() (driver.Value, error) {
	return driverValueJson(s)
}

func (s Schema) MarshalJSON() ([]byte, error) {
	fields := []map[string]any{}
	for _, field := range s {
		fields = append(fields, ConvertSchemaFieldToMap(field))
	}
	return json.Marshal(fields)
}

func (s Schema) Validate(input map[string]any) error {
	var validationErrors ValidationErrors
	for _, field := range s {
		if valid, err := field.Validate(input[field.Name()]); !valid && err != nil {
			validationErrors = append(validationErrors, &ValidationError{
				Field:   field.Name(),
				Message: err.Error(),
			})
		}
	}
	return valErrOrNil(validationErrors)
}

func (s Schema) FindField(fieldName string) SchemaField {
	for _, field := range s {
		if field.Name() == fieldName {
			return field
		}
	}
	return nil
}

// CastValue cast the given v based on the field type
func (s Schema) CastValue(fieldName string, v any) any {
	field := s.FindField(fieldName)
	if field == nil {
		return v
	}
	return field.CastValue(v)
}

type SchemaField interface {
	Name() string
	Label() string
	Type() string
	Properties() map[string]any
	CastValue(input any) any
	Validate(input any) (bool, error)
}

func ConvertSchemaFieldToMap(field SchemaField) map[string]any {
	js := map[string]any{
		"name":       field.Name(),
		"title":      field.Label(),
		"type":       field.Type(),
		"properties": field.Properties(),
	}

	if nestedField, ok := field.(NestableSchemaField); ok {
		childFields := []map[string]any{}

		for _, childField := range nestedField.ChildSchema() {
			childFields = append(childFields, ConvertSchemaFieldToMap(childField))
		}

		if len(childFields) > 0 {
			js["children"] = childFields
		}
	}

	return js
}

func MarshalSchemaFieldJSON(field SchemaField) ([]byte, error) {
	return json.Marshal(ConvertSchemaFieldToMap(field))
}

type BaseField struct {
	FieldName  string
	FieldLabel string
	Required   bool
}

func (f BaseField) Name() string {
	return f.FieldName
}

func (f BaseField) Label() string {
	if len(f.FieldLabel) == 0 {
		return f.FieldName
	}
	return f.FieldLabel
}

func (f BaseField) Properties() map[string]any {
	return map[string]any{
		"required": f.Required,
	}
}

func (f BaseField) mergeProperties(mp map[string]any) map[string]any {
	maps.Copy(mp, f.Properties())
	return mp
}

func (f BaseField) Validate(input any) (bool, error) {
	if f.Required {
		if input == nil {
			return false, fmt.Errorf("value is required")
		}
		if v, ok := input.(string); ok && len(v) == 0 {
			return false, fmt.Errorf("value is required")
		}
	}
	return true, nil
}

type StringSchemaField struct {
	BaseField
	MaxLength int
	MinLength int
}

func (f StringSchemaField) Type() string {
	return "string"
}

func (f StringSchemaField) Properties() map[string]any {
	return f.mergeProperties(map[string]any{
		"max": f.MaxLength,
		"min": f.MinLength,
	})
}

func (f StringSchemaField) CastValue(input any) any {
	switch v := input.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return ""
	}
}

func (f StringSchemaField) Validate(input any) (bool, error) {
	v := f.CastValue(input).(string)
	if f.MinLength != 0 && f.MaxLength != 0 {
		if len(v) > f.MaxLength {
			return false, fmt.Errorf("value is too long")
		}
		if len(v) < f.MinLength {
			return false, fmt.Errorf("value is too short")
		}
	}
	return true, nil
}

type NumberSchemaField struct {
	BaseField
	Min       int
	Max       int
	IsDecimal bool
}

func (f NumberSchemaField) Name() string {
	return f.FieldName
}

func (f NumberSchemaField) Label() string {
	if len(f.FieldLabel) == 0 {
		return f.FieldName
	}
	return f.FieldLabel
}

func (f NumberSchemaField) Type() string {
	return "number"
}

func (f NumberSchemaField) Properties() map[string]any {
	return f.mergeProperties(map[string]any{
		"min":        f.Min,
		"max":        f.Max,
		"is_decimal": f.IsDecimal,
	})
}

func (f NumberSchemaField) CastValue(input any) any {
	switch v := input.(type) {
	case json.Number:
		if f.IsDecimal {
			fv, _ := v.Float64()
			return fv
		} else {
			iv, _ := v.Int64()
			return iv
		}
	case int64:
		return v
	case float64:
		return v
	case int:
		return int64(v)
	case float32:
		return float64(v)
	default:
		return 0
	}
}

func (f NumberSchemaField) Validate(input any) (bool, error) {
	v := f.CastValue(input)
	if iv, ok := v.(int64); ok {
		if f.IsDecimal {
			return false, fmt.Errorf("value is not decimal")
		}

		if iv > int64(f.Max) {
			return false, fmt.Errorf("value is too large")
		}
		if iv < int64(f.Min) {
			return false, fmt.Errorf("value is too small")
		}
	} else if fv, ok := v.(float64); ok {
		if !f.IsDecimal {
			return false, fmt.Errorf("value is not integer")
		}

		if fv > float64(f.Max) {
			return false, fmt.Errorf("value is too large")
		}
		if fv < float64(f.Min) {
			return false, fmt.Errorf("value is too small")
		}
	}
	return true, nil
}

type BooleanSchemaField struct {
	BaseField
}

func (f BooleanSchemaField) Type() string {
	return "boolean"
}

func (f BooleanSchemaField) CastValue(input any) any {
	switch v := input.(type) {
	case bool:
		return v
	default:
		return false
	}
}

func (f BooleanSchemaField) Validate(input any) (bool, error) {
	v, okBool := f.CastValue(input).(bool)
	if !okBool {
		return false, fmt.Errorf("value is not boolean")
	}

	if f.Required && (!okBool || !v) {
		return false, fmt.Errorf("value is required")
	}

	return true, nil
}

type SelectSchemaField struct {
	BaseField
	Options []string
	Min     int
	Max     int
}

func (f SelectSchemaField) Type() string {
	return "select"
}

func (f SelectSchemaField) Properties() map[string]any {
	return f.mergeProperties(map[string]any{
		"options":  f.Options,
		"required": f.Required,
		"min":      f.Min,
		"max":      f.Max,
	})
}

func (f SelectSchemaField) CastValue(input any) any {
	switch v := input.(type) {
	case []string:
		return v
	case string:
		return []string{v}
	case []byte:
		return []string{string(v)}
	default:
		return []string{}
	}
}

func (f SelectSchemaField) Validate(input any) (bool, error) {
	v := f.CastValue(input).([]string)
	if f.Required && len(v) == 0 {
		return false, fmt.Errorf("value is required")
	}
	if len(v) > f.Max {
		return false, fmt.Errorf("value is too long")
	}
	if len(v) < f.Min {
		return false, fmt.Errorf("value is too short")
	}
	for _, item := range v {
		if !slices.Contains(f.Options, item) {
			return false, fmt.Errorf("invalid option")
		}
	}
	return true, nil
}

type RepeaterSchemaField struct {
	BaseField
	BaseSchemaField SchemaField
	MinLength       int
	MaxLength       int
}

func (f RepeaterSchemaField) Type() string {
	return "repeater"
}

func (f RepeaterSchemaField) Properties() map[string]any {
	return f.mergeProperties(map[string]any{
		"base_schema_field": ConvertSchemaFieldToMap(f.BaseSchemaField),
		"min":               f.MinLength,
		"max":               f.MaxLength,
	})
}

func (f RepeaterSchemaField) CastValue(input any) any {
	// todo: add checks for content of input type
	rt := reflect.TypeOf(input)
	if rt.Kind() != reflect.Slice && rt.Kind() != reflect.Array {
		return []any{}
	}
	return input
}

func (f RepeaterSchemaField) Validate(input any) (bool, error) {
	v := f.CastValue(input).([]any)
	if len(v) > f.MaxLength {
		return false, fmt.Errorf("value is too long")
	}
	if len(v) < f.MinLength {
		return false, fmt.Errorf("value is too short")
	}

	var validationErrors ValidationErrors
	for idx, item := range v {
		if valid, err := f.BaseSchemaField.Validate(item); !valid && err != nil {
			validationErrors = append(validationErrors, &ValidationError{
				Field: fmt.Sprintf("%d.%s", idx, f.Name()),
			})
		}
	}
	return len(validationErrors) == 0, valErrOrNil(validationErrors)
}

type RelationSchemaField struct {
	FieldName  string
	FieldLabel string
	Collection *Collection
}

type NestableSchemaField interface {
	SchemaField
	ChildSchema() Schema
}

type NestedSchemaField struct {
	FieldName  string
	FieldLabel string
	Fields     Schema
}

func (f NestedSchemaField) Type() string {
	return "object"
}

func (f NestedSchemaField) CastValue(input any) any {
	mp, ok := input.(map[string]any)
	if !ok {
		return map[string]any{}
	}

	result := map[string]any{}
	for _, field := range f.Fields {
		result[field.Name()] = field.CastValue(mp[field.Name()])
	}

	return result
}

func (f NestedSchemaField) Validate(input any) (bool, error) {
	v := f.CastValue(input).(map[string]any)
	err := f.Fields.Validate(v)
	return err == nil, err
}

func (f NestedSchemaField) ChildSchema() Schema {
	return f.Fields
}

type GroupSchemaField struct {
	BaseField
	Fields Schema
}

func (f GroupSchemaField) Type() string {
	return "group"
}

func (f GroupSchemaField) CastValue(input any) any {
	mp, ok := input.(map[string]any)
	if !ok {
		return map[string]any{}
	}

	result := map[string]any{}
	for _, field := range f.Fields {
		result[field.Name()] = field.CastValue(mp[field.Name()])
	}

	return result
}

func (f GroupSchemaField) Validate(input any) (bool, error) {
	v := f.CastValue(input).(map[string]any)
	err := f.Fields.Validate(v)
	return err == nil, err
}

func (f GroupSchemaField) ChildSchema() Schema {
	return f.Fields
}

type KVGroupSchemaField struct {
	BaseField
	KeySchema   SchemaField
	ValueSchema SchemaField
}

func (f KVGroupSchemaField) Type() string {
	return "kv_group"
}

func (f KVGroupSchemaField) CastValue(input any) any {
	mp, ok := input.(map[string]any)
	if !ok {
		return map[string]any{}
	}
	return mp
}

func (f KVGroupSchemaField) Validate(input any) (bool, error) {
	v := f.CastValue(input).(map[string]any)
	var validationErrors ValidationErrors
	for key, value := range v {
		if valid, err := f.KeySchema.Validate(key); !valid && err != nil {
			validationErrors = append(validationErrors, &ValidationError{
				Field: fmt.Sprintf("%s.%s", key, f.KeySchema.Name()),
			})
		}
		if valid, err := f.ValueSchema.Validate(value); !valid && err != nil {
			validationErrors = append(validationErrors, &ValidationError{
				Field: fmt.Sprintf("%s.%s", key, f.ValueSchema.Name()),
			})
		}
	}
	return len(validationErrors) == 0, valErrOrNil(validationErrors)
}

type CustomSchemaFieldFactory struct {
	FieldType       string                                                  `json:"type"`
	FieldProperties map[string]any                                          `json:"properties"`
	Children        Schema                                                  `json:"children"`
	Validator       func(field *CustomSchemaField, input any) (bool, error) `json:"-"`
	Caster          func(field *CustomSchemaField, input any) any           `json:"-"`
}

func (f *CustomSchemaFieldFactory) GetProperty(name string, defaultValue any) any {
	value, exists := f.FieldProperties[name]
	if !exists {
		return defaultValue
	}
	return value
}

func (f *CustomSchemaFieldFactory) Create(name string, label string) SchemaField {
	return &CustomSchemaField{
		CustomSchemaFieldFactory: f,
		BaseField: BaseField{
			FieldName:  name,
			FieldLabel: label,
		},
	}
}

type CustomSchemaField struct {
	*CustomSchemaFieldFactory
	BaseField
	FieldProperties map[string]any `json:"properties"`
}

func (f *CustomSchemaField) Type() string {
	return f.FieldType
}

func (f *CustomSchemaField) Properties() map[string]any {
	props := map[string]any{}
	maps.Copy(props, f.CustomSchemaFieldFactory.FieldProperties)
	maps.Copy(props, f.FieldProperties)
	return f.mergeProperties(props)
}

func (f *CustomSchemaField) CastValue(input any) any {
	if f.Caster != nil {
		return f.Caster(f, input)
	} else if len(f.Children) == 1 {
		return f.Children[0].CastValue(input)
	}
	return input
}

func (f *CustomSchemaField) Validate(input any) (bool, error) {
	if f.Validator != nil {
		return f.Validator(f, input)
	} else if len(f.Children) == 1 {
		return f.Children[0].Validate(input)
	}
	return true, nil
}

func (f CustomSchemaField) ChildSchema() Schema {
	return f.Children
}

func stringFieldValueCaster(field *CustomSchemaField, input any) any {
	switch v := input.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return ""
	}
}

var RichTextSchemaField = &CustomSchemaFieldFactory{
	FieldType: "rich_text",
	Caster:    stringFieldValueCaster,
	Validator: func(field *CustomSchemaField, input any) (bool, error) {
		return true, nil
	},
}

// var EmailTextSchemaField = &CustomSchemaFieldFactory{
// 	FieldType: "email",
// 	FieldProperties: map[string]any{
// 		"allowed_domains": []string{},
// 	},
// 	Validator: func(field *CustomSchemaField, input any) (bool, error) {
// 		// allowedDomains := field.GetProperty("allowed_domains", []string{}).([]string)
// 		// if len(allowedDomains) != 0 {

// 		// }
// 	},
// }
