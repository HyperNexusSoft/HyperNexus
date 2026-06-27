package tools

import (
	"--model"
	"--sandbox"
	"--skip-git-repo-check"
	"codex"
	"exec"
	"model"
	"prompt"
	"sandbox"
)

func HandleCodex(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	prompt, _ :=getString(args, "prompt")
	if e != nil {
		return err(e)
}

	cmd := exec.Command("codex", "exec", "--model", getString(args, "model"), "--sandbox", getString(args, "sandbox"), "--skip-git-repo-check", prompt)
	var outb, errb io.ReadCloser
	var e error

	if outb, errb, e = cmd.StdoutPipe(); e != nil {
		return err(e)
}

	defer outb.Close()
	defer errb.Close()

	if e = cmd.Start(); e != nil {
		return err(e)

}

}
