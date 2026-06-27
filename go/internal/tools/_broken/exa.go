package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
)

// ToolResponse, ok, e, getString, getInt, getBool, TextContent は parity.go で定義されていると仮定します。

func HandleX(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// 引数の取得
	key, _ :=getString(args, "key")
	if key == "" {
		return err("キーが必要です")
}

	// ここで exa コマンドを実行します
	cmd := exec.Command("exa", key)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	// コマンドの実行
	if e := cmd.Run(); e != nil {
		return err(fmt.Sprintf("exa コマンドの実行に失敗しました: %s", errb.String()))
}

	// 出力の取得
	output := outb.String()
	if output == "" {
		return err("exa コマンドの出力が空です")
}

	// 成功時の処理
	return ok(output)
}