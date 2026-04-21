package envfile

import (
	"fmt"
	"strings"
)

// ConvertFormat describes a supported conversion target format.
type ConvertFormat string

const (
	FormatDockerCompose ConvertFormat = "docker-compose"
	FormatShell         ConvertFormat = "shell"
	FormatDotenv        ConvertFormat = "dotenv"
)

// ConvertResult holds the converted output string.
type ConvertResult struct {
	Format  ConvertFormat
	Content string
}

// Convert transforms a map of env entries into the specified output format.
func Convert(entries map[string]string, format ConvertFormat, maskSecrets bool) (*ConvertResult, error) {
	keys := sortedKeys(entries)

	switch format {
	case FormatDockerCompose:
		return convertDockerCompose(entries, keys, maskSecrets), nil
	case FormatShell:
		return convertShell(entries, keys, maskSecrets), nil
	case FormatDotenv:
		return convertDotenv(entries, keys, maskSecrets), nil
	default:
		return nil, fmt.Errorf("unsupported convert format: %q", format)
	}
}

func convertDockerCompose(entries map[string]string, keys []string, maskSecrets bool) *ConvertResult {
	var sb strings.Builder
	sb.WriteString("environment:\n")
	for _, k := range keys {
		v := valueOf(entries[k], k, maskSecrets)
		sb.WriteString(fmt.Sprintf("  - %s=%s\n", k, v))
	}
	return &ConvertResult{Format: FormatDockerCompose, Content: sb.String()}
}

func convertShell(entries map[string]string, keys []string, maskSecrets bool) *ConvertResult {
	var sb strings.Builder
	for _, k := range keys {
		v := valueOf(entries[k], k, maskSecrets)
		sb.WriteString(fmt.Sprintf("export %s=%s\n", k, quoteIfNeeded(v)))
	}
	return &ConvertResult{Format: FormatShell, Content: sb.String()}
}

func convertDotenv(entries map[string]string, keys []string, maskSecrets bool) *ConvertResult {
	var sb strings.Builder
	for _, k := range keys {
		v := valueOf(entries[k], k, maskSecrets)
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, quoteIfNeeded(v)))
	}
	return &ConvertResult{Format: FormatDotenv, Content: sb.String()}
}
