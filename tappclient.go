package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
)

func main() {

	const configPath = "./tapp/client.xml"
	const endPoint = "blueprint/script_characterizations?type=boot"
	const timestampLayout = "2006-01-02T15:04:05.000000-07:00"

	// Create an http client
	config := openTappConfiguration(configPath)
	client := createClient(config)

	// Get scripts
	response, err := client.Get(config.ApiEndpoint + endPoint)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	// Parse them
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var scriptChars []ScriptCharacterization
	json.Unmarshal(body, &scriptChars)

	// Sort by execution order
	sort.Sort(ByOrder(scriptChars))

	// Execute them sequentially
	for _, ex := range scriptChars {
		log.Println("Executing :\n", ex.Script.Code)
		output, exitCode, startedAt, finishedAt := ExecCode(ex.Script.Code)
		scriptConclusion := ScriptConclusion{
			UUID:       ex.UUID,
			Output:     output,
			ExitCode:   exitCode,
			StartedAt:  startedAt.Format(timestampLayout),
			FinishedAt: finishedAt.Format(timestampLayout),
		}
		log.Println("Conclusion :\n", scriptConclusion)
	}
}
