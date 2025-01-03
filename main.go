package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"main.go/models"
)

func main() {
	var config models.Config

	log.Println("reading file..")
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}
	log.Println(config.CharacterName)

	items, err := parseGBankClassicDB(config.InputFilePath, config.OutputDirectory)
	if err != nil {
		log.Fatal(err)
	}

	getEvents(config)
	postToDiscord(config, items)

}

// postToDiscord sends a message to the Discord channel
func postToDiscord(config models.Config, items []models.Item) {
	url := fmt.Sprintf("https://raid-helper.dev/api/v2/servers/%s/channels/%s/embed", config.ServerID, config.DiscordChannel)
	payload := buildRaidHelperEmbedMessage(items)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON payload: %v", err)
		return
	}

	//nicely format the json payload
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, jsonPayload, "", "  "); err != nil {
		log.Printf("Error indenting JSON payload: %v", err)
		return
	}

	log.Printf("Sending HTTP request to %s with payload:\n%s", url, prettyJSON.String())

	//test it out with dummy data for now
	dummyTest := buildRaidHelperEmbedMessageString(items)

	cmd := exec.Command("curl", "--request", "POST", "--url", url,
		"--header", "Authorization: "+config.RaidHelperAPIKey,
		"--header", "Content-Type: application/json; charset=utf-8",
		"--data", dummyTest)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error sending HTTP request: %v", err)
		return
	}

	if strings.Contains(string(output), "error") {
		log.Printf("Request failed with output: %s", output)
		return
	}

	log.Println("Message sent successfully")
}

func getEvents(config models.Config) {
	url := fmt.Sprintf("https://raid-helper.dev/api/v3/servers/%s/events", config.ServerID)

	cmd := exec.Command("curl", "--request", "GET", "--url", url,
		"--header", "Authorization: "+config.RaidHelperAPIKey)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error sending HTTP request: %v", err)
		return
	}

	if strings.Contains(string(output), "error") {
		log.Printf("Request failed with output: %s", output)
		return
	}

	//parse events
	json.Marshal(output)

	log.Println("Message sent successfully")
}

// buildRaidHelperEmbedMessage builds the embed message
func buildRaidHelperEmbedMessage(items []models.Item) models.RaidHelperEmbedMessage {

	embedMessage := models.RaidHelperEmbedMessage{}
	embedMessage.Title.Text = "Test"
	embedMessage.Description = "Test"
	embedMessage.Fields = []models.RaidHelperEmbedField{}

	for _, item := range items {
		embedMessage.Fields = append(embedMessage.Fields, models.RaidHelperEmbedField{
			Name:   item.Info.Name,
			Value:  fmt.Sprintf("%d", item.Count),
			Inline: true,
		})
	}

	return embedMessage
}

func buildRaidHelperEmbedMessageString(items []models.Item) string {
	//time in mountain time
	time := time.Now().Format("2006-01-02-15-04-05 MST")

	itemses := ""
	for _, item := range items {
		itemses += `\n` + item.Info.Name + `: ` + strconv.Itoa(item.Count)
	}
	output := `{ "title":{ "text": "Wizard Vault"}, "description": "` + itemses + `","footer":{ "text": "` + time + `"}}`

	//print items in a single line

	fmt.Println(output)

	return output
}

func parseItem(scanner *bufio.Scanner) (models.Item, error) {
	item := models.Item{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[\"ID\"]") {
			match := regexp.MustCompile(`\[\"ID\"\] = (\d+),`).FindStringSubmatch(line)
			if match != nil {
				item.Id = match[1]
			}
		} else if strings.HasPrefix(line, "[\"Info\"]") {
			match := regexp.MustCompile(`\[\"Info\"\] = {`).FindStringSubmatch(line)
			if match != nil {
				itemInfo, err := parseInfo(scanner)
				if err != nil {
					return item, err
				}

				item.Info = itemInfo
			}

		} else if strings.HasPrefix(line, "[\"Count\"]") {
			match := regexp.MustCompile(`\[\"Count\"\] = (\d+),`).FindStringSubmatch(line)
			if match != nil {
				count, err := strconv.Atoi(match[1])
				if err != nil {
					return item, err
				}
				item.Count = count
			}

		} else if strings.HasPrefix(line, "[\"Link\"]") { //this is the singular item without info path
			match := regexp.MustCompile(`\[\"Link\"\] = (.*?),`).FindStringSubmatch(line)
			if match != nil {
				link := match[1]
				parts := strings.Split(link, "[")
				if len(parts) > 1 {
					link = parts[1]
					parts = strings.Split(link, "]")
					if len(parts) > 1 {
						item.Info.Name = parts[0]
					}
				}
			}
		} else if strings.HasPrefix(line, "}") {
			break
		}
	}
	return item, nil
}

func parseInfo(scanner *bufio.Scanner) (models.Info, error) {
	info := models.Info{}
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "[\"icon\"]"):
			if match := regexp.MustCompile(`\[\"icon\"\] = (\d+)`).FindStringSubmatch(line); match != nil {
				info.Icon, _ = strconv.Atoi(match[1])
			}
		case strings.HasPrefix(line, "[\"level\"]"):
			if match := regexp.MustCompile(`\[\"level\"\] = (\d+)`).FindStringSubmatch(line); match != nil {
				info.Level, _ = strconv.Atoi(match[1])
			}
		case strings.HasPrefix(line, "[\"rarity\"]"):
			if match := regexp.MustCompile(`\[\"rarity\"\] = (\d+)`).FindStringSubmatch(line); match != nil {
				info.Rarity, _ = strconv.Atoi(match[1])
			}
		case strings.HasPrefix(line, "[\"equipId\"]"):
			if match := regexp.MustCompile(`\[\"equipId\"\] = (\d+)`).FindStringSubmatch(line); match != nil {
				info.EquipId, _ = strconv.Atoi(match[1])
			}
		case strings.HasPrefix(line, "[\"price\"]"):
			if match := regexp.MustCompile(`\[\"price\"\] = (\d+)`).FindStringSubmatch(line); match != nil {
				info.Price, _ = strconv.Atoi(match[1])
			}
		case strings.HasPrefix(line, "[\"class\"]"):
			if match := regexp.MustCompile(`\[\"class\"\] = (\d+)`).FindStringSubmatch(line); match != nil {
				info.Class, _ = strconv.Atoi(match[1])
			}
		case strings.HasPrefix(line, "[\"subClass\"]"):
			if match := regexp.MustCompile(`\[\"subClass\"\] = (\d+)`).FindStringSubmatch(line); match != nil {
				info.SubClass, _ = strconv.Atoi(match[1])
			}
		case strings.HasPrefix(line, "[\"name\"]"):
			if match := regexp.MustCompile(`\[\"name\"\] = \"([^\"]+)\"`).FindStringSubmatch(line); match != nil {
				info.Name = match[1]
			}
		case strings.HasPrefix(line, "}"):
			return info, nil
		}
	}
	return info, fmt.Errorf("reached end of info block without closing bracket")
}
func parseItems(scanner *bufio.Scanner) ([]models.Item, error) {
	items := []models.Item{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "items") {
			continue
		}
		if strings.HasPrefix(line, "{") {
			item, err := parseItem(scanner)
			if err != nil {
				return items, err
			}
			items = append(items, item)
		}
	}

	//some items have the same name so we need to add the count together
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].Info.Name == items[j].Info.Name {
				items[i].Count += items[j].Count
				items = append(items[:j], items[j+1:]...)
			}
		}
	}
	//alphabetize
	sort.Slice(items, func(i, j int) bool {
		return items[i].Info.Name < items[j].Info.Name
	})

	return items, nil
}

func writeItemsToCSV(items []models.Item, fileOut string) error {
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	fileOut = fileOut + "-" + timestamp + ".csv"
	f, err := os.Create(fileOut)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	//make the seperator a tab
	w.Comma = ','
	defer w.Flush()

	if err := w.Write([]string{"ITEM_NAME", "COUNT"}); err != nil {
		return err
	}
	for _, item := range items {
		if err := w.Write([]string{item.Info.Name, strconv.Itoa(item.Count)}); err != nil {
			return err
		}
	}
	return nil
}

func parseGBankClassicDB(fileIn, fileOut string) (items []models.Item, err error) {
	items = []models.Item{}
	file, err := os.Open(fileIn)
	if err != nil {
		return items, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	parsedItems, err := parseItems(scanner)
	if err != nil {
		return items, err
	}
	items = append(items, parsedItems...)
	log.Printf("found %d items", len(items))

	return items, writeItemsToCSV(items, fileOut)
}
