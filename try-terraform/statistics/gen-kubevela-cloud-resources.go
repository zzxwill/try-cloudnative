package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type TFDownload struct {
	Data     []DataItem     `json:"data"`
	Included []IncludedItem `json:"included"`
}

type IncludedItem struct {
	Id         string     `json:"id"`
	Attributes Attributes `json:"attributes"`
}

type DataItem struct {
	Attributes    Attributes    `json:"attributes"`
	Relationships Relationships `json:"relationships"`
}

type Relationships struct {
	LatestVersion RelationshipLatestVersion `json:"latest-version"`
}

type RelationshipLatestVersion struct {
	Data RelationshipData `json:"data"`
}

type RelationshipData struct {
	Id string `json:"id"`
}

type Attributes struct {
	Name        string `json:"name"`
	Downloads   int    `json:"downloads"`
	Source      string `json:"source"`
	Description string `json:"description"`
	Verified    bool   `json:"verified"`
}

func main() {
	// var terraformModulesUrl = "https://registry.terraform.io/v2/modules?filter%5Bprovider%5D=aws&include=latest-version&page%5Bsize%5D=50&page%5Bnumber%5D=1"
	if len(os.Args) < 2 {
		fmt.Println("Please provide the cloud provider name and an official Terraform modules URL")
		os.Exit(1)
	}
	providerName := os.Args[1]
	terraformModulesUrl := os.Args[2]
	resp, err := http.Get(terraformModulesUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var modules TFDownload
	if err := json.Unmarshal(body, &modules); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if _, err = os.Stat(providerName); err == nil {
		if err := os.RemoveAll(providerName); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Successfully deleted existed directory %s\n", providerName)
	}
	if _, err = os.Stat(providerName); os.IsNotExist(err) {
		if err := os.Mkdir(providerName, 0755); err != nil {
			if !os.IsExist(err) {
				log.Fatal(err)
			}
			fmt.Printf("Successfully created directory %s\n", providerName)
		}
	}

	for _, module := range modules.Data {
		var description string
		for _, attr := range modules.Included {
			if module.Relationships.LatestVersion.Data.Id == attr.Id {
				description = attr.Attributes.Description
			}
		}
		if description == "" {
			description = strings.ToUpper(providerName) + " " + strings.Title(module.Attributes.Name)
		}

		outputFile := fmt.Sprintf("%s/terraform-%s-%s.yaml", providerName, providerName, module.Attributes.Name)
		if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
			continue
		}
		if providerName == "aws" && (module.Attributes.Name == "rds" || module.Attributes.Name == "s3-bucket" ||
			module.Attributes.Name == "subnet" || module.Attributes.Name == "vpc") {
			continue
		}
		cmd := fmt.Sprintf("vela def init %s --type component --provider %s --git %s.git --desc \"%s\" -o %s",
			module.Attributes.Name, providerName, module.Attributes.Source, description, outputFile)
		fmt.Println(cmd)
		stdout, err := exec.Command("bash", "-c", cmd).CombinedOutput()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(string(stdout))
	}
}
