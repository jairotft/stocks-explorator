package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type OpenAIMessagePayload struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIPayload struct {
	MaxTokens        int                    `json:"max_tokens"`
	Temperature      float32                `json:"temperature"`
	FrequencyPenalty int                    `json:"frequency_penalty"`
	PresencePenalty  int                    `json:"presence_penalty"`
	TopP             float32                `json:"top_p"`
	Stop             *int                   `json:"stop"`
	Messages         []OpenAIMessagePayload `json:"messages"`
	Model            *string                `json:"model"`
}

type OpenAIContentFilterItemResponse struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity"`
}

type OpenAIContentFilterResultResponse struct {
	Hate     OpenAIContentFilterItemResponse `json:"hate"`
	SelfHarm OpenAIContentFilterItemResponse `json:"self_harm"`
	Sexual   OpenAIContentFilterItemResponse `json:"sexual"`
	Violence OpenAIContentFilterItemResponse `json:"violence"`
}

type OpenAIPromptAnotation struct {
	PromptIndex          int                               `json:"prompt_index"`
	ContentFilterResults OpenAIContentFilterResultResponse `json:"content_filter_results"`
}

type OpenAIChoice struct {
	Index                int                               `json:"index"`
	FinishReason         string                            `json:"finish_reason"`
	Message              OpenAIMessagePayload              `json:"message"`
	ContentFilterResults OpenAIContentFilterResultResponse `json:"content_filter_results"`
}

type OpenAIUsageResponse struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenAIResponse struct {
	ID                string                  `json:"id"`
	Object            string                  `json:"object"`
	Created           int64                   `json:"created"`
	Model             string                  `json:"model"`
	PromptAnnotations []OpenAIPromptAnotation `json:"prompt_annotations"`
	Choices           []OpenAIChoice          `json:"choices"`
	Usage             OpenAIUsageResponse     `json:"usage"`
}

type OpenAICredentialChannel struct {
	BASE     string
	ENGINE   string
	VERSION  string
	MODEL    string
	PROVIDER string
	KEY      string
}

func CreateChat(data OpenAIPayload) (OpenAIResponse, error) {
	var url string
	method := "POST"

	getEnv := func(name string) string {
		return os.Getenv(fmt.Sprintf(`OPENAI_API_%s`, name))
	}

	credential := OpenAICredentialChannel{
		BASE:     getEnv("BASE"),
		ENGINE:   getEnv("ENGINE"),
		VERSION:  getEnv("VERSION"),
		MODEL:    getEnv("MODEL"),
		PROVIDER: getEnv("PROVIDER"),
		KEY:      getEnv("KEY"),
	}

	url = fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s", credential.BASE, credential.ENGINE, credential.VERSION)

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return OpenAIResponse{}, err
	}

	payload := strings.NewReader(string(jsonBytes))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return OpenAIResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", credential.KEY)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return OpenAIResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return OpenAIResponse{}, err
	}

	if res.StatusCode == 200 {
		openAIResponse := OpenAIResponse{}
		err_unmarshal := json.Unmarshal(body, &openAIResponse)
		if err_unmarshal != nil {
			return OpenAIResponse{}, err_unmarshal
		}

		return openAIResponse, nil
	}

	return OpenAIResponse{}, nil
}
