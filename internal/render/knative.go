package render

import (
	"fmt"
	"strings"

	"github.com/<your-module>/golden-path-platform/internal/config"
	"github.com/<your-module>/golden-path-platform/internal/templates"
)

func RenderKnativeService(app *config.App) ([]byte, error) {
	min := int32(0)
	max := int32(10)
	if app.Spec.Scale.Min != nil {
		min = *app.Spec.Scale.Min
	}
	if app.Spec.Scale.Max != nil {
		max = *app.Spec.Scale.Max
	}

	envBlock := renderEnvBlock(app.Spec.Env)
	out := string(templates.KnativeServiceYAML)

	repls := map[string]string{
		"{{NAME}}":       app.Metadata.Name,
		"{{NAMESPACE}}":  app.Metadata.Namespace,
		"{{IMAGE}}":      app.Spec.Image,
		"{{PORT}}":       fmt.Sprintf("%d", app.Spec.Port),
		"{{MIN_SCALE}}":  fmt.Sprintf("%d", min),
		"{{MAX_SCALE}}":  fmt.Sprintf("%d", max),
		"{{ENV_BLOCK}}":  envBlock,
	}

	for k, v := range repls {
		out = strings.ReplaceAll(out, k, v)
	}

	return []byte(out), nil
}

func renderEnvBlock(envs []config.Env) string {
	// envが空なら、空配列にしてYAMLを壊さない
	if len(envs) == 0 {
		// env: [] にしたいのでインデントを合わせる
		return "            []"
	}

	var b strings.Builder
	for _, e := range envs {
		// service.yaml の env: の下は 12スペース + "- ..."
		b.WriteString("            - name: " + escapeYAML(e.Name) + "\n")
		b.WriteString("              value: " + quoteYAML(e.Value) + "\n")
	}
	// 最後の改行を消して見栄え調整（無くても動く）
	return strings.TrimRight(b.String(), "\n")
}

func escapeYAML(s string) string {
	// name は基本プレーンでOK（必要なら厳密化）
	return s
}

func quoteYAML(s string) string {
	// value はクォートして安全側
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}
