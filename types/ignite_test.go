package types

import (
	"embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/json/*
var testFiles embed.FS

func getFile(t *testing.T, input bool) []byte {
	part := "/output.json"
	if input {
		part = "/input.json"
	}
	name := "testdata/json/" + t.Name() + part
	b, err := testFiles.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

type errType interface {
	error() string
}

type unmarshalErr string

func (e unmarshalErr) error() string { return string(e) }

type marshalErr string

func (e marshalErr) error() string { return string(e) }

func TestPresetForm_JSONMarshalling(t *testing.T) {
	tests := []struct {
		name string

		wantsErr  errType
		mutations func(form *PresetForm)
	}{
		// Unmarshal errors.

		{
			name:     "root fields missing",
			wantsErr: unmarshalErr("(root): v is required, (root): fields is required"),
		},
		{
			name:     "version wrong",
			wantsErr: unmarshalErr("(root): v is required"),
		},
		{
			name: "invalid version unmarshal",
			wantsErr: unmarshalErr("v: Must validate one and only one schema " +
				"(oneOf), v: v does not match: 1"),
		},
		{
			name: "fields empty object",
			wantsErr: unmarshalErr("fields.0: input is required, fields.0:" +
				" title is required, fields.0: required is required, fields.0: map_to is required"),
		},
		{
			name: "blank input object",
			wantsErr: unmarshalErr("fields.0.input: Must validate one and only one schema (oneOf)," +
				" fields.0.input: type is required"),
		},
		{
			name: "invalid autogen input object",
			wantsErr: unmarshalErr("fields.0.input: Must validate one and only one schema (oneOf), " +
				"fields.0.input.autogen: Must validate one and only one schema (oneOf), " +
				"fields.0.input.autogen: fields.0.input.autogen does not match: \"PROJECT_NAMESPACE\""),
		},
		{
			name: "invalid field map to",
			wantsErr: unmarshalErr("fields.0.map_to.0: Must validate one and only one schema (oneOf), " +
				"fields.0.map_to.0: type is required"),
		},

		// Marshal error.

		{
			name: "check marshal also uses schema",
			mutations: func(form *PresetForm) {
				form.Version = 100
			},
			wantsErr: marshalErr("json: error calling MarshalJSON for type *types.PresetForm: v: Must validate " +
				"one and only one schema (oneOf), v: v does not match: 1"),
		},

		// Success

		{
			name: "version defaults",
			mutations: func(form *PresetForm) {
				form.Version = 0
			},
		},
		{
			name: "valid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := getFile(t, true)

			p := &PresetForm{}
			err := json.Unmarshal(input, &p)
			if tt.wantsErr == nil {
				assert.NoError(t, err)
			} else if u, ok := tt.wantsErr.(unmarshalErr); ok {
				assert.EqualError(t, err, string(u))
				return
			}

			if tt.mutations != nil {
				tt.mutations(p)
			}

			b, err := json.Marshal(p)
			if tt.wantsErr == nil {
				assert.NoError(t, err)
			} else if m, ok := tt.wantsErr.(marshalErr); ok {
				assert.EqualError(t, err, string(m))
				return
			}

			output := getFile(t, false)
			assert.JSONEq(t, string(output), string(b))
		})
	}
}
