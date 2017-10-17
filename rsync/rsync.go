package rsync

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/MarkLux/JudgeServer/config"
)

func SyncSingle(testCaseId string) (err error) {
	remote := os.Getenv("RSYNC_USER") + "@" + os.Getenv("RSYNC_HOST")
	rsyncPwd := filepath.Join(config.TEST_CASE_DIR, "rsync.passwd")
	rsyncCmd := `rsync -av --password-file ` + rsyncPwd + ` --include="*.in" --include="*.out"  --include="*/"  --exclude="*" --delete` + remote + `::testcases/` + testCaseId + ` ` + config.TEST_CASE_DIR
	in := bytes.NewBuffer(nil)
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = in
	in.WriteString(rsyncCmd)
	err = cmd.Run()
	return
}
