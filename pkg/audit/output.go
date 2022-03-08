package audit

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"text/tabwriter"

	"github.com/kubesphere/kubeeye/apis/kubeeye/v1alpha1"
	"github.com/pkg/errors"
)

func defaultOutput(receiver <-chan []v1alpha1.AuditResults) error {
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 3, ' ', 0)
	_, err := fmt.Fprintln(w, "\nNAMESPACE\tKIND\tNAME\tLEVEL\tMESSAGE\tREASON")
	if err != nil {
		return err
	}
	for r := range receiver {
		for _, results := range r {
			for _, resultInfo := range results.ResultInfos {
				for _, items := range resultInfo.ResultItems {
					s := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%-8v", results.NameSpace, resultInfo.ResourceType,
						resultInfo.Name, items.Level, items.Message, items.Reason)
					_, err := fmt.Fprintln(w, s)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

func JSONOutput(receiver <-chan []v1alpha1.AuditResults) error {
	var output []v1alpha1.AuditResults
	filename := "audit_result.json"
	// create csv file
	newFile, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "create file audit_result.json failed.")
	}
	defer newFile.Close()

	for r := range receiver {
		for _, results := range r {
			output = append(output, results)
		}
	}

	// output json
	jsonOutput, err := json.MarshalIndent(output, "", "    ")
	
	if err != nil {
		return err
	}
	fmt.Println(string(jsonOutput))
	writeErr := ioutil.WriteFile(filename, jsonOutput, 0644)
	if writeErr != nil {
		return writeErr
	}
	fmt.Printf("\033[1;36;49m请查阅审计报告 audit_result.json .\033[0m\n")
	return nil
}

func CSVOutput(receiver <-chan []v1alpha1.AuditResults) error {
	filename := "audit_result.csv"
	// create csv file
	newFile, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "create file audit_result.csv failed.")
	}

	defer newFile.Close()

	// write UTF-8 BOM to prevent print gibberish.
	if _, err = newFile.WriteString("\xEF\xBB\xBF"); err != nil {
		return err
	}

	// NewWriter returns a new Writer that writes to w.
	w := csv.NewWriter(newFile)
	header := []string{"namespace", "kind", "name", "level", "message", "reason"}
	contents := [][]string{
		header,
	}
	for r := range receiver {
		for _, results := range r {
			var resourceName string
			for _, resultInfo := range results.ResultInfos {
				for _, items := range resultInfo.ResultItems {
					if resourceName == "" {
						content := []string{
							results.NameSpace,
							resultInfo.ResourceType,
							resultInfo.Name,
							items.Level,
							items.Message,
							items.Reason,
						}
						contents = append(contents, content)
						resourceName = resultInfo.Name
					} else {
						content := []string{
							"",
							"",
							"",
							items.Level,
							items.Message,
							items.Reason,
						}
						contents = append(contents, content)
					}
				}
			}
		}
	}
	// WriteAll writes multiple CSV records to w using Write and then calls Flush,
	if err := w.WriteAll(contents); err != nil {
		return err
	}
	fmt.Printf("\033[1;36;49m请查阅审计报告audit_result.CSV.\033[0m\n")
	return nil
}
