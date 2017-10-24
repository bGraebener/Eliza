package elizaHelper

import (
	"encoding/json"
	"io/ioutil"
)

// loads the resources from the json file
func LoadResources() map[string][]string {

	dataMap := make(map[string][]string)

	// read the json file
	raw, err := ioutil.ReadFile("./res/elizaData.json")
	if err != nil {
		panic("Couldn't read resource file!")
	}

	// parse the json data
	if err := json.Unmarshal(raw, &dataMap); err != nil {
		panic("Couldn't parse json file")
	}
	return dataMap

}

// converts a string slice into a map, convience function for faster, easier lookup of keywords and responses
func SliceToMap(slice []string) map[string]int {

	tmpMap := make(map[string]int)

	for _, i := range slice {
		tmpMap[i]++
	}
	return tmpMap
}
