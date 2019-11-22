package handlers

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func getTestHash(raw string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(raw), 14)
	return string(bytes)
}

func Test_ResponsePayloadGenerator_Generate(t *testing.T) {
	type args struct {
		request  map[string]interface{}
		response map[string]interface{}
	}

	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Simple payload substitution",
			args: args{
				request: map[string]interface{}{
					"inputName": "something",
				},
				response: map[string]interface{}{
					"name": "{{ .Request.inputName }}",
				},
			},
			want: map[string]interface{}{
				"name": "something",
			},
		},
		{
			name: "Multiple field substitution",
			args: args{
				request: map[string]interface{}{
					"firstInput":  "hello",
					"secondInput": "goodbye",
				},
				response: map[string]interface{}{
					"firstOutput":  "{{ .Request.firstInput }}",
					"secondOutput": "{{ .Request.secondInput }}",
				},
			},
			want: map[string]interface{}{
				"firstOutput":  "hello",
				"secondOutput": "goodbye",
			},
		},
		{
			name: "Multiple nested fields substitution",
			args: args{
				request: map[string]interface{}{
					"firstInput":  "hello",
					"secondInput": "goodbye",
				},
				response: map[string]interface{}{
					"firstOutput": "{{ .Request.firstInput }}",
					"nested": map[string]interface{}{
						"secondOutput": "{{ .Request.secondInput }}",
					},
				},
			},
			want: map[string]interface{}{
				"firstOutput": "hello",
				"nested": map[string]interface{}{
					"secondOutput": "goodbye",
				},
			},
		},
		{
			name: "Multiple nested input fields substitution",
			args: args{
				request: map[string]interface{}{
					"firstInput": "hello",
					"nested": map[string]interface{}{
						"secondInput": "goodbye",
					},
				},
				response: map[string]interface{}{
					"firstOutput":  "{{ .Request.firstInput }}",
					"secondOutput": "{{ .Request.nested.secondInput }}",
				},
			},
			want: map[string]interface{}{
				"firstOutput":  "hello",
				"secondOutput": "goodbye",
			},
		},
		{
			name: "Can use additional templating features",
			args: args{
				request: map[string]interface{}{
					"firstInput": "hello",
					"beNice":     true,
				},
				response: map[string]interface{}{
					"firstOutput":  "{{ .Request.firstInput | upper | repeat 2}}",
					"secondOutput": "{{if .Request.beNice -}} Nice to see you {{- else -}} ya chump {{- end }}",
				},
			},
			want: map[string]interface{}{
				"firstOutput":  "HELLOHELLO",
				"secondOutput": "Nice to see you",
			},
		},
		{
			name: "Can support defaults",
			args: args{
				request: map[string]interface{}{
					"firstInput": "hello",
				},
				response: map[string]interface{}{
					"firstOutput":  "{{ default \"missing\" .Request.firstInput }}",
					"secondOutput": "{{ default \"missing\" .Request.secondInput }}",
				},
			},
			want: map[string]interface{}{
				"firstOutput":  "hello",
				"secondOutput": "missing",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := ResponsePayloadGenerator{tt.args.request, nil}
			got := generator.Generate(tt.args.response)

			fmt.Println(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResponsePayloadGenerator.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
