package main

import ( 
	"flag"
	"fmt"
	"os"
	"passkeeper/internal/writer"
	"passkeeper/internal/crypto"
	"github.com/atotto/clipboard"
)

var GeneratedKey = crypto.GenerateKey()
var KeyFileName = "secret.key"

func main() {
	service := flag.String("service", "", "Name of the service (e.g. github)")
	key := flag.String("key", "", "Key for the service (e.g. username or password)")
	value := flag.String("value", "", "Value for the service (e.g. 123 or pass)")
	action := flag.String("action", "", "Action to perform (e.g. show, add, copy)")

	flag.Parse()

	if *service == "" {
		fmt.Println("Error: --service, --action required.")
		flag.Usage()
		return
	}

	if *action == "add" && (*key == "" || *value == "") {
		fmt.Println("Error: --value, --key required.")
		flag.Usage()
		return
	}

	writer.InitServiceStorage(*service)

	if _, err := os.Stat(KeyFileName); err != nil {
		err := os.WriteFile(KeyFileName, GeneratedKey, 0600) 
		if err != nil {
			fmt.Printf("failed to save key: %v", err)
		}
	}
	loadedKey, err := os.ReadFile("secret.key")
	if err != nil {
		fmt.Printf("failed to load key: %v", err)
	}

	switch *action {
		case "show":
			result, err := writer.ShowValue(*key, *service, *action, loadedKey)
			if err != nil {
				fmt.Println("error copying to clipboard:", err)
				return
			}
			fmt.Println(result)
		case "copy":
			result, err := writer.ShowValue(*key, *service, *action, loadedKey)
			if err != nil {
				fmt.Println("error copying to clipboard:", err)
				return
			}
			err = clipboard.WriteAll(result)
			if err != nil {
				fmt.Println("error copying to clipboard:", err)
				return
			}
		case "decode":
			err = writer.DecryptFile(*service, *action, loadedKey)
			if err != nil {
				fmt.Println("error copying to clipboard:", err)
				return
			}
		case "encode":
			err = writer.EncryptFile(*service, *action, loadedKey)
			if err != nil {
				fmt.Println("error copying to clipboard:", err)
				return
			}
		case "add":
			err = writer.WriteValue(*key, *value, *service, loadedKey)
			if err != nil {
				fmt.Println("error copying to clipboard:", err)
				return
			}
		case "remove":
			err = writer.RemoveValue(*key, *service)
			if err != nil {
				fmt.Println("error copying to clipboard:", err)
				return
			}
		default:
			err = writer.ShowList(*service)
			if err != nil {
				fmt.Println("error copying to clipboard:", err)
				return
			}
	}
}
