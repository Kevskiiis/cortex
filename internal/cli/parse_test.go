package cli

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    []string
		wantErr bool
	}{
		{name: "empty", line: "   ", want: []string{}},
		{name: "command and flags", line: "doctor --verbose", want: []string{"doctor", "--verbose"}},
		{name: "quoted argument", line: `group create --name "Help Desk Admins"`, want: []string{"group", "create", "--name", "Help Desk Admins"}},
		{name: "does not expand environment", line: `echo $CORTEX_TEST_VALUE`, want: []string{"echo", "$CORTEX_TEST_VALUE"}},
		{name: "unclosed quote", line: `group create --name "Help Desk`, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseLine(test.line)
			if test.wantErr {
				if err == nil {
					t.Fatal("ParseLine() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseLine() error = %v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("ParseLine() = %#v, want %#v", got, test.want)
			}
		})
	}
}
