package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

const (
	characterizationsEndpoint = "blueprint/script_characterizations?type=boot"
	conclusionsEndpoint       = "blueprint/script_conclusions"
	configPath                = "./tapp/client.xml"
	timestampLayout           = "2006-01-02T15:04:05.000000-07:00"
)

func main() {
	// Create an http client
	config := openTappConfiguration(configPath)
	client := createClient(config)

	// Scripts retrieval
	scriptChars := retrieveScripts(config, client)

	// Sort by execution order
	log.Println("Sorting scripts")
	sort.Sort(ByOrder(scriptChars))

	// Execute them sequentially and put conclusions in a channel
	conclusions := make(chan ScriptConclusion, len(scriptChars))
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
		conclusions <- scriptConclusion
	}
	close(conclusions)

	// Send conclusions back to the server and put responses on a channel
	// (I guess this can be done concurrently)
	log.Println("Sending conclusions back to the server")
	responses := make(chan []string)
	for conclusion := range conclusions {
		go sendConclusion(config, client, conclusion, responses)
	}

	// Log responses
	for i := 0; i < len(scriptChars); i++ {
		resp := <-responses
		log.Printf("Got response to %v : %v", resp[0], resp[1])
	}
}

func retrieveScripts(config TappConfig, client *http.Client) (scriptChars []ScriptCharacterization) {
	// Get scripts
	log.Println("Retrieving scripts")
	response, err := client.Get(config.ApiEndpoint + characterizationsEndpoint)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	// Parse them
	log.Println("Parsing scripts")
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Received : ", string(body))
	json.Unmarshal(body, &scriptChars)

	return
}

func sendConclusion(config TappConfig, client *http.Client, conclusion ScriptConclusion, responses chan []string) {
	var url = config.ApiEndpoint + conclusionsEndpoint
	wrapper := ConclusionWrapper{Conclusion: conclusion}
	j, err := json.Marshal(wrapper)
	if err != nil {
		log.Fatalln("Marshalling error: ", err)
		responses <- []string{conclusion.UUID, err.Error()}
		return
	}

	log.Println("Posting ", string(j))
	resp, err := client.Post(url, "application/json", strings.NewReader(string(j)))
	if err != nil {
		log.Fatalln("Error sending conclusion: ", err)
		responses <- []string{conclusion.UUID, err.Error()}
		return
	}
	defer resp.Body.Close()

	responses <- []string{conclusion.UUID, resp.Status}
}
