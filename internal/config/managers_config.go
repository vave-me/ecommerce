package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
)

// ManagersConfig holds the configuration for the Media service (MinIO-based).
type ManagersConfig struct {
	OpenAIAPIKey         string `envconfig:"OPENAI_API_KEY"`
	OpenAIBaseURL        string `envconfig:"OPENAI_BASE_URL"`
	OpenAIBaseModel      string `envconfig:"OPENAI_BASE_MODEL"`
	AnthropicAPIKey      string `envconfig:"ANTHROPIC_API_KEY"`
	AnthropicBaseURL     string `envconfig:"ANTHROPIC_BASE_URL"`
	AnthropicModel       string `envconfig:"ANTHROPIC_MODEL"`
	DeepSeekAPIKey       string `envconfig:"DEEPSEEK_API_KEY"`
	DeepSeekBaseURL      string `envconfig:"DEEPSEEK_BASE_URL"`
	DeepSeekModel        string `envconfig:"DEEPSEEK_DEFAULT_MODEL"`
	GoogleAIAAPIKey      string `envconfig:"GOOGLE_AI_API_KEY"`
	GoogleAIABaseURL     string `envconfig:"GOOGLE_AI_BASE_URL"`
	GoogleAIDefaultModel string `envconfig:"GOOGLE_AI_DEFAULT_MODEL"`
}

// LoadMediaConfig loads configuration for the Media service from environment variables
func LoadManagersConfig() *ManagersConfig {
	fmt.Println("Loading environment variables for MediaConfig...")

	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	if openAIAPIKey == "" {
		openAIAPIKey = "" // default if not set
		fmt.Println("OPENAI_API_KEY not set, using default:", openAIAPIKey)
	} else {
		fmt.Println("OPENAI_API_KEY:", openAIAPIKey)
	}

	openAIendpoint := os.Getenv("OPENAI_BASE_URL")
	if openAIendpoint == "" {
		openAIendpoint = "https://api.openai.com/v1" // default if not set
		fmt.Println("OPENAI_BASE_URL not set, using default:", openAIendpoint)
	} else {
		fmt.Println("OPENAI_BASE_URL:", openAIendpoint)
	}
	openAIBaseModel := os.Getenv("OPENAI_BASE_MODEL")
	if openAIBaseModel == "" {
		openAIBaseModel = "gpt-4o-mini" // Use string constant instead of ai2.ModelGPT4oMini
		fmt.Println("OPENAI_BASE_MODEL not set, using default:", openAIendpoint)
	} else {
		fmt.Println("OPENAI_BASE_MODEL:", openAIendpoint)
	}

	anthropicAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	if anthropicAPIKey == "" {
		anthropicAPIKey = "" // default if not set
		fmt.Println("anthropicAPIKey not set, using default:", anthropicAPIKey)
	} else {
		fmt.Println("OPENAI_API_KEY:", anthropicAPIKey)
	}

	anthropicEndpoint := os.Getenv("ANTHROPIC_BASE_URL")
	if anthropicEndpoint == "" {
		anthropicEndpoint = "https://api.anthropic.com" // default if not set
		fmt.Println("ANTHROPIC_BASE_URL not set, using default:", anthropicEndpoint)
	} else {
		fmt.Println("ANTHROPIC_BASE_URL:", anthropicEndpoint)
	}
	anthropicBaseModel := os.Getenv("ANTHROPIC_MODEL")
	if anthropicBaseModel == "" {
		anthropicBaseModel = "claude-3-5-haiku-20241022" // Use string constant instead of ai2.ModelClaude35Haiku20241022
		fmt.Println("ANTHROPIC_MODEL not set, using default:", anthropicBaseModel)
	} else {
		fmt.Println("ANTHROPIC_MODEL:", anthropicBaseModel)
	}
	deepSeekAPIKey := os.Getenv("DEEPSEEK_API_KEY")
	if deepSeekAPIKey == "" {
		deepSeekAPIKey = "" // default if not set
		fmt.Println("DEEPSEEK_API_KEY not set, using default:", deepSeekAPIKey)
	} else {
		fmt.Println("DEEPSEEK_API_KEY:", deepSeekAPIKey)
	}

	deepSeekEndpoint := os.Getenv("DEEPSEEK_BASE_URL")
	if deepSeekEndpoint == "" {
		deepSeekEndpoint = "https://api.deepseek.com/v1" // default if not set
		fmt.Println("DEEPSEEK_BASE_URL not set, using default:", deepSeekEndpoint)
	} else {
		fmt.Println("DEEPSEEK_BASE_URL:", anthropicEndpoint)
	}
	deepSeekModel := os.Getenv("DEEPSEEK_DEFAULT_MODEL")
	if deepSeekModel == "" {
		deepSeekModel = "deepseek-r1-0528" // Use string constant instead of ai2.ModelDeepSeekR1_0528
		fmt.Println("DEEPSEEK_DEFAULT_MODEL not set, using default:", deepSeekModel)
	} else {
		fmt.Println("DEEPSEEK_DEFAULT_MODEL:", deepSeekModel)
	}

	googleAPIKey := os.Getenv("GOOGLE_AI_API_KEY")
	if googleAPIKey == "" {
		googleAPIKey = "" // default if not set
		fmt.Println("GOOGLE_AI_API_KEY not set, using default:", googleAPIKey)
	} else {
		fmt.Println("GOOGLE_AI_API_KEY:", googleAPIKey)
	}

	googleAIBaseUrl := os.Getenv("GOOGLE_AI_BASE_URL")
	if googleAIBaseUrl == "" {
		googleAIBaseUrl = "https://api.deepseek.com/v1" // default if not set
		fmt.Println("GOOGLE_AI_BASE_URL not set, using default:", googleAIBaseUrl)
	} else {
		fmt.Println("GOOGLE_AI_BASE_URL:", googleAIBaseUrl)
	}
	googleDefaultModel := os.Getenv("ANTHROPIC_MODEL")
	if googleDefaultModel == "" {
		googleDefaultModel = ""
		fmt.Println("GOOGLE_AI_DEFAULT_MODEL not set, using default:", googleDefaultModel)
	} else {
		fmt.Println("GOOGLE_AI_DEFAULT_MODEL:", googleDefaultModel)
	}

	return &ManagersConfig{
		OpenAIAPIKey:         openAIAPIKey,
		OpenAIBaseURL:        openAIBaseModel,
		OpenAIBaseModel:      openAIBaseModel,
		AnthropicAPIKey:      anthropicAPIKey,
		AnthropicBaseURL:     anthropicEndpoint,
		AnthropicModel:       anthropicBaseModel,
		DeepSeekAPIKey:       deepSeekAPIKey,
		DeepSeekBaseURL:      deepSeekEndpoint,
		DeepSeekModel:        deepSeekModel,
		GoogleAIAAPIKey:      googleAPIKey,
		GoogleAIABaseURL:     googleAIBaseUrl,
		GoogleAIDefaultModel: googleDefaultModel,
	}
}

func InitManagersConfig() (managersConfig ManagersConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &managersConfig)
	return
}
