package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
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

var errNoVariables = errors.New("failed to find main.tf or variables.tf in Terraform configurations")

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
		if err := generateDefinition(providerName, module.Attributes.Name, module.Attributes.Source, "", description); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if err := checkDocGen(providerName, module.Attributes.Name); err != nil {
			if err == errNoVariables {
				if err := regenerateDefinition(providerName, module.Attributes.Name, module.Attributes.Source, description); err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
		}
	}
}

func checkDocGen(provider, name string) error {
	defName := fmt.Sprintf("%s-%s", provider, name)
	defYaml := filepath.Join(provider, fmt.Sprintf("terraform-%s.yaml", defName))
	cmd := fmt.Sprintf("kubectl apply -f %s", defYaml)
	stdout, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(stdout))
	}
	fmt.Println(string(stdout))
	fmt.Printf("Generate doc for %s", defName)
	cmd = fmt.Sprintf("vela def doc-gen %s -n vela-system", defName)
	stdout, err = exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		if strings.Contains(string(stdout), errNoVariables.Error()) {
			return errNoVariables
		}
		return err
	}
	return nil
}

func generateDefinition(provider, name, gitURL, path, description string) error {
	defYaml := filepath.Join(provider, fmt.Sprintf("terraform-%s-%s.yaml", provider, name))

	cmd := fmt.Sprintf("vela def init %s --type component --provider %s --git %s.git --desc \"%s\" -o %s",
		name, provider, gitURL, description, defYaml)
	if path != "" {
		cmd = fmt.Sprintf("%s --path %s", cmd, path)
	}
	fmt.Println(cmd)
	stdout, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(stdout))
	return nil
}

func regenerateDefinition(provider, name, gitUrl, description string) error {
	defYaml := filepath.Join(provider, fmt.Sprintf("terraform-%s-%s.yaml", provider, name))
	if err := os.Remove(defYaml); err != nil {
		return err
	}
	tmpPath := "./tmp/terraform"
	// Check if the directory exists. If yes, remove it.
	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		err := os.RemoveAll(tmpPath)
		if err != nil {
			return errors.Wrap(err, "failed to remove the directory")
		}
	}
	_, err := git.PlainClone(tmpPath, false, &git.CloneOptions{
		URL:      gitUrl,
		Progress: nil,
	})
	if err != nil {
		return err
	}

	infos, err := ioutil.ReadDir(filepath.Join(tmpPath, "modules"))
	if err != nil {
		return err
	}
	for _, info := range infos {
		var newName = info.Name()
		if name == "cloudwatch" || name == "route53" {
			newName = name + "-" + info.Name()
		}
		path := filepath.Join("modules", info.Name())
		if err := generateDefinition(provider, newName, gitUrl, path, description); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if err := checkDocGen(provider, newName); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	return nil
}
