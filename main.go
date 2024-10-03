package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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

	parseGBankClassicDB(config.InputFilePath, config.OutputDirectory)

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
	var info models.Info
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[\"icon\"]") {
			match := regexp.MustCompile(`\[\"icon\"\] = (\d+)`).FindStringSubmatch(line)
			if match != nil {
				info.Icon, _ = strconv.Atoi(match[1])
			}
		} else if strings.HasPrefix(line, "[\"level\"]") {
			match := regexp.MustCompile(`\[\"level\"\] = (\d+)`).FindStringSubmatch(line)
			if match != nil {
				info.Level, _ = strconv.Atoi(match[1])
			}
		} else if strings.HasPrefix(line, "[\"rarity\"]") {
			match := regexp.MustCompile(`\[\"rarity\"\] = (\d+)`).FindStringSubmatch(line)
			if match != nil {
				info.Rarity, _ = strconv.Atoi(match[1])
			}
		} else if strings.HasPrefix(line, "[\"equipId\"]") {
			match := regexp.MustCompile(`\[\"equipId\"\] = (\d+)`).FindStringSubmatch(line)
			if match != nil {
				info.EquipId, _ = strconv.Atoi(match[1])
			}
		} else if strings.HasPrefix(line, "[\"price\"]") {
			match := regexp.MustCompile(`\[\"price\"\] = (\d+)`).FindStringSubmatch(line)
			if match != nil {
				info.Price, _ = strconv.Atoi(match[1])
			}
		} else if strings.HasPrefix(line, "[\"class\"]") {
			match := regexp.MustCompile(`\[\"class\"\] = (\d+)`).FindStringSubmatch(line)
			if match != nil {
				info.Class, _ = strconv.Atoi(match[1])
			}
		} else if strings.HasPrefix(line, "[\"subClass\"]") {
			match := regexp.MustCompile(`\[\"subClass\"\] = (\d+)`).FindStringSubmatch(line)
			if match != nil {
				info.SubClass, _ = strconv.Atoi(match[1])
			}
		} else if strings.HasPrefix(line, "[\"name\"]") {
			match := regexp.MustCompile(`\[\"name\"\] = \"([^\"]+)\"`).FindStringSubmatch(line)
			if match != nil {
				info.Name = match[1]
			}
		}
		if strings.HasPrefix(line, "}") {
			break
		}
	}
	return info, nil
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

func parseGBankClassicDB(fileIn, fileOut string) error {
	file, err := os.Open(fileIn)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	items, err := parseItems(scanner)
	if err != nil {
		return err
	}
	log.Printf("found %d items", len(items))

	return writeItemsToCSV(items, fileOut)
}
