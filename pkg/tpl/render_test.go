package tpl

import "testing"

func TestRender(t *testing.T) {
	type args struct {
		s   string
		ctx map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple template",
			args: args{
				s:   "Hello, {{.name}}!",
				ctx: map[string]interface{}{"name": "World"},
			},
			want: "Hello, World!",
		},
		{
			name: "template with if statement",
			args: args{
				s:   "{{if .show}}Hello, {{.name}}!{{end}}",
				ctx: map[string]interface{}{"name": "World", "show": true},
			},
			want: "Hello, World!",
		},
		{
			name: "template with range statement",
			args: args{
				s:   "{{range .items}}{{.}}\n{{end}}",
				ctx: map[string]interface{}{"items": []string{"foo", "bar", "baz"}},
			},
			want: "foo\nbar\nbaz\n",
		},
		{
			name: "invalid template",
			args: args{
				s:   "{{.name}",
				ctx: map[string]interface{}{"name": "World"},
			},
			want: "{{.name}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Render(tt.args.s, tt.args.ctx); got != tt.want {
				t.Errorf("Render() = %v, want %v", got, tt.want)
			}
		})
	}
}
