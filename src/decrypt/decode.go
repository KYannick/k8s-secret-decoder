package main

import (
	"bufio"
	b64 "encoding/base64"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type secret struct {
	APIVersion string
	Kind       string
	Type       string
	Metadata   map[string]string
	Data       map[string]string
}

var cmdlineArgs = os.Args[1:]

func init() {
	if len(cmdlineArgs) != 0 {
		fmt.Println("no params are allowed")
		os.Exit(1)
	}
}

func main() {
	//read stdin
	scanner := bufio.NewScanner(os.Stdin)
	var inputyaml string
	var mysecret secret
	for scanner.Scan() {
		inputyaml = inputyaml + scanner.Text()
		inputyaml = inputyaml + "\n"
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	// convert to struct
	inputbytearray := []byte(inputyaml)
	err := yaml.Unmarshal(inputbytearray, &mysecret)
	if err != nil {
		panic(err)
	}

	//decode base64 strings if yaml has kind secret
	if mysecret.Kind == "Secret" {
		//fmt.Println("secret found")
		for k, v := range mysecret.Data {
			//fmt.Printf("key: %s, value: %s \n", key, value)
			decoded, err := b64.StdEncoding.DecodeString(v)
			if err != nil {
				panic(err)
			}
			mysecret.Data[k] = fmt.Sprintf("%s", decoded)
		}
	} else {
		fmt.Println("no secret yaml")
		os.Exit(1)
	}

	//convert back to string
	out, err := yaml.Marshal(mysecret)
	if err != nil {
		panic(err)
	}

	//print result
	fmt.Printf("%s", out)
}
