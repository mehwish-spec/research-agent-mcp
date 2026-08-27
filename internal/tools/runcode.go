package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var supportedLanguages = map[string]langConfig{
	"python": {
		image: "python:3.11-slim",
		shell: "python3",
		flag:  "-c",
	},
	"javascript": {
		image: "node:20-slim",
		shell: "node",
		flag:  "-e",
	},
	"bash": {
		image: "alpine:3.20",
		shell: "sh",
		flag:  "-c",
	},
}

type langConfig struct {
	image string
	shell string
	flag  string
}

// RunCode executes code in an isolated Docker container and returns the output.
// The container is automatically removed after execution.
func RunCode(language, code string) (string, error) {
	lang := strings.ToLower(strings.TrimSpace(language))
	cfg, ok := supportedLanguages[lang]
	if !ok {
		supported := make([]string, 0, len(supportedLanguages))
		for k := range supportedLanguages {
			supported = append(supported, k)
		}
		return "", fmt.Errorf("unsupported language %q. Supported: %s", language, strings.Join(supported, ", "))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// docker run --rm          → auto-remove after exit
	//            --network none → no internet access inside container
	//            --memory 64m  → cap RAM at 64MB
	//            --cpus 0.5    → cap at half a CPU core
	args := []string{
		"run", "--rm",
		"--network", "none",
		"--memory", "64m",
		"--cpus", "0.5",
		cfg.image,
		cfg.shell, cfg.flag, code,
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := strings.TrimSpace(stdout.String())
	errOutput := strings.TrimSpace(stderr.String())

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("execution timed out after 30 seconds")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Language: %s\nImage:    %s\n\n", lang, cfg.image))

	if output != "" {
		sb.WriteString("Output:\n" + output + "\n")
	}
	if errOutput != "" {
		sb.WriteString("\nStderr:\n" + errOutput + "\n")
	}
	if err != nil && output == "" && errOutput == "" {
		sb.WriteString("Error: " + err.Error() + "\n")
	}
	if output == "" && errOutput == "" {
		sb.WriteString("(no output)\n")
	}

	return sb.String(), nil
}
