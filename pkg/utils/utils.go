package utils

import (
	"bytes"
	"encoding/json"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
)

func IsEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func ToPointer[T any](value T) *T {
	return &value
}

func WriteKubeconfigToFile(config *clientcmdapi.Config, format, output string) error {
	if format == "json" {
		yamlBytes, err := clientcmd.Write(*config)
		if err != nil {
			return err
		}
		jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
		if err != nil {
			return err
		}

		dir := filepath.Dir(output)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
		if err := os.WriteFile(output, PrettyJsonBytes(jsonBytes), 0600); err != nil {
			return err
		}
	} else {
		if err := clientcmd.WriteToFile(*config, output); err != nil {
			return err
		}
	}
	return nil
}

func PrettyJsonBytes(b []byte) []byte {
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, b, "", "    ")
	return prettyJSON.Bytes()
}
