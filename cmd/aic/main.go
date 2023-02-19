package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"crg.eti.br/go/ai-companion/luaengine"
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
	"github.com/PullRequestInc/go-gpt3"
)

type Config struct {
	Prompt       string `json:"prompt" ini:"prompt" cfg:"prompt" cfgDefault:"-" cfgDescription:"Prompt to use for GPT-3. Use - to read from stdin."`
	OpenAIAPIKey string `json:"openai_api_key" ini:"openai_api_key" cfg:"openai_api_key" cfgDescription:"OpenAI API key."`
}

func readStdin() string {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func main() {
	luaengine := luaengine.New()
	err := luaengine.Compile("aic.lua")
	if err != nil {
		panic(err)
	}
	luaengine.InitState()
	defer luaengine.Close()

	var cfg Config
	config.File = "config.json"
	config.Parse(&cfg)

	prompt := cfg.Prompt
	if cfg.Prompt == "-" {
		prompt = readStdin()
	}

	client := gpt3.NewClient(cfg.OpenAIAPIKey)
	ctx := context.Background()

	//buf := strings.Builder{}
	err = client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			prompt,
		},
		MaxTokens:   gpt3.IntPtr(100),
		Temperature: gpt3.Float32Ptr(0.7), // TODO: make this configurable,
	}, func(resp *gpt3.CompletionResponse) {
		//buf.WriteString(resp.Choices[0].Text)
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		log.Printf("GPT-3 error: %v", err)
		return
	}

	/*
		response := buf.String()
		if len(response) > 0 {
			fmt.Println(response)
		}
	*/

	fmt.Println()
}
