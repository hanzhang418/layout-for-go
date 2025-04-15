package env

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantActive  Profile
		wantWarning bool
	}{
		{
			name:        "dev environment",
			args:        []string{"-env", "dev"},
			wantActive:  Dev,
			wantWarning: false,
		},
		{
			name:        "fat environment",
			args:        []string{"-env", "fat"},
			wantActive:  Fat,
			wantWarning: false,
		},
		{
			name:        "uat environment",
			args:        []string{"-env", "uat"},
			wantActive:  Uat,
			wantWarning: false,
		},
		{
			name:        "pro environment",
			args:        []string{"-env", "pro"},
			wantActive:  Pro,
			wantWarning: false,
		},
		{
			name:        "pro with whitespace",
			args:        []string{"-env", " pro "},
			wantActive:  Pro,
			wantWarning: false,
		},
		{
			name:        "invalid environment",
			args:        []string{"-env", "invalid"},
			wantActive:  Dev,
			wantWarning: true,
		},
		{
			name:        "empty environment",
			args:        []string{"-env", ""},
			wantActive:  Dev,
			wantWarning: true,
		},
		{
			name:        "no environment flag",
			args:        []string{},
			wantActive:  Dev,
			wantWarning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 备份原始命令行参数和标准输出
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = append([]string{"cmd"}, tt.args...)

			// 重置flag解析器
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// 捕获标准输出
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() {
				os.Stdout = oldStdout
				w.Close()
			}()

			// 执行初始化函数
			Init()

			// 读取捕获的输出
			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// 验证active的值
			if Current() != tt.wantActive {
				t.Errorf("current = %v, want %v", current, tt.wantActive)
			}

			// 验证警告信息
			warningMsg := "Warning: '-env' cannot be found, or it is illegal. The default 'dev' will be used.\n"
			if tt.wantWarning {
				if !strings.Contains(output, warningMsg) {
					t.Errorf("expected warning message not found in output: %q", output)
				}
			} else {
				if strings.Contains(output, warningMsg) {
					t.Errorf("unexpected warning message in output: %q", output)
				}
			}
		})
	}
}
