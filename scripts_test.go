package main

import (
	"encoding/json"
	"sort"
	"testing"
)

/*
	Test that the unmarshalling is correct
*/
func TestUnmarshal(t *testing.T) {
	const jsonSample = `[
	  {
	    "execution_order": 1,
	    "parameter_values": {
	
	    },
	    "type": "boot",
	    "uuid": "53c3b86e63051f336b00036f",
	    "script": {
	      "code": "#!/bin/bash\nhostname\npwd",
	      "uuid": "53c3b79963051f7a8800036b",
	      "attachment_paths": [
	      ]
	    }
	  },
	  {
	    "execution_order": 4611686018427387903,
	    "parameter_values": {
	    },
	    "type": "boot",
	    "uuid": "53c3b73163051f6d46000367",
	    "script": {
	      "code": "echo .",
	      "uuid": "4ea19048d907b10559000003",
	      "attachment_paths": [
	      ]
	    }
	  }
	]`

	var scriptChars []ScriptCharacterization

	json.Unmarshal([]byte(jsonSample), &scriptChars)

	if len(scriptChars) != 2 {
		t.Errorf("Expected 2 execution items")
	}

	if scriptChars[0].Order != 1 {
		t.Errorf("Test failed. Incorrect Order %v at %v ", scriptChars[0].Order, scriptChars[0])
	}

	if scriptChars[0].UUID != "53c3b86e63051f336b00036f" {
		t.Errorf("Test failed. Incorrect UUID %v at %v ", scriptChars[0].UUID, scriptChars[0])
	}

	if scriptChars[0].Script.Code != "#!/bin/bash\nhostname\npwd" {
		t.Errorf("Test failed. Incorrect Script Code %v at %v ", scriptChars[0].Script.Code, scriptChars[0])
	}

	if scriptChars[0].Script.UUID != "53c3b79963051f7a8800036b" {
		t.Errorf("Test failed. Incorrect Script UUID %v at %v ", scriptChars[0].Script.UUID, scriptChars[0])
	}
}

/*
	Test the execution of correct code
*/
func TestExecCode(t *testing.T) {
	const code = "echo Hello!"
	const expectedOutput = "Hello!\n"

	output, exitCode, startedAt, finishedAt := ExecCode(code)

	if output != expectedOutput {
		t.Errorf("Output was %v but expected was %v", output, expectedOutput)
	}

	if exitCode != 0 {
		t.Errorf("Exit code was %v but expected was %v", exitCode, 0)
	}

	if &startedAt == nil {
		t.Errorf("Start timestamp was nil")
	}

	if &finishedAt == nil {
		t.Errorf("End timestamp was nil")
	}

	if startedAt.After(finishedAt) {
		t.Errorf("Inconsistent Timestamps")
	}
}

/*
	Test the execution of code with an exit error != 0
*/
func TestExecBadCode(t *testing.T) {
	const code = "lk"

	_, exitCode, startedAt, finishedAt := ExecCode(code)

	if exitCode != 127 {
		t.Errorf("Exit code was %v but expected was %v", exitCode, 127)
	}

	if &startedAt == nil {
		t.Errorf("Start timestamp was nil")
	}

	if &finishedAt == nil {
		t.Errorf("End timestamp was nil")
	}

	if startedAt.After(finishedAt) {
		t.Errorf("Inconsistent Timestamps")
	}
}

/*
	Test reordering Script Characterizations by order field
*/
func TestSortByOrder(t *testing.T) {
	var scriptCharacterizations = []ScriptCharacterization{{Order: 2}, {Order: 3}, {Order: 1}}
	sort.Sort(ByOrder(scriptCharacterizations))

	for i, ex := range scriptCharacterizations {
		if ex.Order != i+1 {
			t.Errorf("Sorting Script Characterizations fails!")
		}
	}
}
