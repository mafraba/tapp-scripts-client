package main

import "time"
import "io/ioutil"
import "os/exec"
import "syscall"
import "log"
import "os"

type ScriptCharacterization struct {
	Order  int    `json:"execution_order"`
	UUID   string `json:"uuid"`
	Script Script `json:"script"`
}

type Script struct {
	Code string `json:"code"`
	UUID string `json:"uuid"`
}

type ScriptConclusion struct {
	UUID       string `json:"script_characterization_id"`
	Output     string `json:"output"`
	ExitCode   int    `json:"exit_code"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
}

type ConclusionWrapper struct {
	Conclusion ScriptConclusion `json:"script_conclusion"`
}

// ByOrder implements sort.Interface for []ScriptCharacterization based on the Order field
type ByOrder []ScriptCharacterization

func (a ByOrder) Len() int           { return len(a) }
func (a ByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

/*
	Execute script and return output, exit code, and init/end timestamps
*/
func ExecCode(code string) (output string, exitCode int, startedAt time.Time, finishedAt time.Time) {
	// First we'll create a temp file
	tmpFile, err := ioutil.TempFile("", "tappscript-tmp")
	if err != nil {
		log.Fatalf("Error creating temp file : ", err)
		return
	}
	// close and remove the file when we dont need it anymore
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Then put the code in it
	_, err = tmpFile.WriteString(code)
	if err != nil {
		log.Fatalf("Error writing to file : ", err)
		return
	}

	// Then we execute that file with 'sh'
	cmd := exec.Command("sh", tmpFile.Name())
	startedAt = time.Now()
	bytes, err := cmd.CombinedOutput()
	finishedAt = time.Now()
	output = string(bytes)
	exitCode = extractExitCode(err)

	// Return the output, exit code, and timestamps
	return
}

/*
	Extracting the exit code is not trivial ...
*/
func extractExitCode(err error) int {
	if err != nil {
		exiterr := err.(*exec.ExitError)
		status := exiterr.Sys().(syscall.WaitStatus)
		return status.ExitStatus()
	} else {
		return 0
	}
}
