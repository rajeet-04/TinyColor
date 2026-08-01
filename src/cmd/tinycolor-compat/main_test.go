package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestRunHumanCommands(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		want   string
		verify func(*testing.T, string)
	}{
		{
			name: "parse JSON inspection",
			args: []string{"parse", "--json", "red"},
			verify: func(t *testing.T, output string) {
				t.Helper()
				var got map[string]any
				if err := json.Unmarshal([]byte(output), &got); err != nil {
					t.Fatal(err)
				}
				if got["hex"] != "ff0000" || got["valid"] != true {
					t.Fatalf("parse JSON = %#v", got)
				}
			},
		},
		{
			name: "parse invalid JSON inspection",
			args: []string{"parse", "--json", "not-a-color"},
			verify: func(t *testing.T, output string) {
				t.Helper()
				var got map[string]any
				if err := json.Unmarshal([]byte(output), &got); err != nil {
					t.Fatal(err)
				}
				if got["valid"] != false {
					t.Fatalf("parse JSON = %#v", got)
				}
			},
		},
		{name: "convert HSL", args: []string{"convert", "--to", "hsl", "red"}, want: "hsl(0, 100%, 50%)\n"},
		{name: "lighten", args: []string{"lighten", "--amount", "10", "#000"}, want: "#1a1a1a\n"},
		{
			name: "triad palette",
			args: []string{"palette", "--type", "triad", "red"},
			want: "#ff0000\n#00ff00\n#0000ff\n",
		},
		{
			name: "triad palette JSON inspections",
			args: []string{"palette", "--type", "triad", "--json", "red"},
			verify: func(t *testing.T, output string) {
				t.Helper()
				var got []map[string]any
				if err := json.Unmarshal([]byte(output), &got); err != nil {
					t.Fatal(err)
				}
				if len(got) != 3 || got[0]["valid"] != true || got[0]["value"] != "red" {
					t.Fatalf("palette JSON = %#v", got)
				}
			},
		},
		{
			name: "contrast JSON",
			args: []string{"contrast", "--json", "#000", "#fff"},
			verify: func(t *testing.T, output string) {
				t.Helper()
				var got map[string]any
				if err := json.Unmarshal([]byte(output), &got); err != nil {
					t.Fatal(err)
				}
				want := map[string]any{"ratio": float64(21), "aaSmall": true, "aaLarge": true, "aaaSmall": true, "aaaLarge": true}
				for key, value := range want {
					if got[key] != value {
						t.Fatalf("contrast %s = %#v, want %#v", key, got[key], value)
					}
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if status := run(test.args, strings.NewReader(""), &stdout, &stderr); status != 0 {
				t.Fatalf("status = %d, stderr = %q", status, stderr.String())
			}
			if test.verify != nil {
				test.verify(t, stdout.String())
				return
			}
			if stdout.String() != test.want {
				t.Fatalf("stdout = %q, want %q", stdout.String(), test.want)
			}
		})
	}
}

func TestRunJSONLAndUsageErrors(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		stdin      string
		wantStatus int
		wantStdout string
		verify     func(*testing.T, string)
	}{
		{
			name:       "JSONL inspect without arguments",
			stdin:      "{\"id\":\"red\",\"operation\":\"inspect\",\"input\":\"red\"}\n",
			wantStdout: "{\"id\":\"red\",\"result\":{\"alpha\":1,\"format\":\"name\",\"original\":\"red\",\"rgb\":{\"a\":1,\"b\":0,\"g\":0,\"r\":255},\"valid\":true,\"value\":\"red\"}}\n",
		},
		{
			name:       "JSONL string forwards explicit format",
			stdin:      "{\"id\":\"string-format\",\"operation\":\"string\",\"input\":\"red\",\"args\":{\"format\":\"hsl\"}}\n",
			wantStdout: "{\"id\":\"string-format\",\"result\":\"hsl(0, 100%, 50%)\"}\n",
		},
		{
			name:       "JSONL inspect preserves one-percent HSL channels",
			stdin:      "{\"id\":\"hsl-one-percent\",\"operation\":\"inspect\",\"input\":\"hsl(115, 1%, 1%)\"}\n",
			wantStdout: "{\"id\":\"hsl-one-percent\",\"result\":{\"alpha\":1,\"format\":\"hsl\",\"original\":\"hsl(115, 1%, 1%)\",\"rgb\":{\"a\":1,\"b\":3,\"g\":3,\"r\":3},\"valid\":true,\"value\":\"hsl(115, 1%, 1%)\"}}\n",
		},
		{
			name:  "malformed JSONL request does not stop the stream",
			stdin: "{bad json\n{\"id\":\"red\",\"operation\":\"inspect\",\"input\":\"red\"}\n",
			verify: func(t *testing.T, output string) {
				t.Helper()
				lines := strings.Split(strings.TrimSpace(output), "\n")
				if len(lines) != 2 || !strings.Contains(lines[0], "malformed JSON") || !strings.Contains(lines[1], "\"id\":\"red\"") {
					t.Fatalf("stdout = %q", output)
				}
			},
		},
		{name: "unknown command", args: []string{"unknown"}, wantStatus: 2},
		{name: "missing convert target", args: []string{"convert", "red"}, wantStatus: 2},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if status := run(test.args, strings.NewReader(test.stdin), &stdout, &stderr); status != test.wantStatus {
				t.Fatalf("status = %d, want %d", status, test.wantStatus)
			}
			if test.verify != nil {
				test.verify(t, stdout.String())
			} else if stdout.String() != test.wantStdout {
				t.Fatalf("stdout = %q, want %q", stdout.String(), test.wantStdout)
			}
			if test.wantStatus == 2 && !strings.Contains(stderr.String(), "Usage: tinycolor ") {
				t.Fatalf("stderr = %q, want tinycolor usage", stderr.String())
			}
		})
	}
}
