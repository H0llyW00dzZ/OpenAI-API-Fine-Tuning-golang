# Fine-tuning OpenAI API GPT-3.5 Turbo model Uploader

![Language: Go](https://img.shields.io/badge/Language-Go-blue.svg)

Previously, I wrote in Python, but now I'm rewriting it in Golang. This repository contains a Go implementation for fine-tuning the OpenAI GPT-3.5 Turbo model, utilizing the OpenAI API.

## Installation

1. Make sure you have Go installed on your system. If not, you can download and install it from the official website: <https://go.dev>

2. Clone this repository to your local machine:

   ```shell
   git clone https://github.com/H0llyW00dzZ/openai-api-fine-tuning-golang.git
   ```

3. Navigate to the project directory:

   ```shell
   cd openai-api-fine-tuning-golang
   ```

4. Install the required dependencies:

   ```shell
   go get -u github.com/fatih/color
   go get -u github.com/tidwall/gjson
   ```

## Usage

To use the OpenAI GPT-3.5 Turbo Fine-Tuning tool, follow these steps:

1. Obtain an OpenAI API token. You can sign up for an API key at <https://platform.openai.com/account/api-keys> .

2. Prepare your training data file in a text format.

3. Run the fine-tuning tool with the following command:

   ```shell
   go run main.go -file /path/to/training/file.jsonl -token YOUR_API_TOKEN
   ```

   Replace `/path/to/training/file.jsonl` with the actual path to your training data file, and `YOUR_API_TOKEN` with your OpenAI API token.

4. The tool will upload the training data file, create a fine-tuning job, and wait for the job to complete. The progress will be displayed in the console, including the elapsed time.

5. Once the fine-tuning job is completed, the fine-tuned model ID will be displayed. You can use this ID to access the fine-tuned model for further use.

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
