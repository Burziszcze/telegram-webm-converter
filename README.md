# telegram-webm-converter is WebM to MP4 Converter Bot

This repository contains a Telegram bot that converts WebM files to MP4 format. The bot listens for messages containing WebM files or URLs to WebM files, downloads the files, converts them to MP4 using FFmpeg, and sends the converted files back to the user.

## Features

- **Download WebM files**: The bot can download WebM files either from Telegram messages or from URLs provided in the messages.
- **Convert to MP4**: Uses FFmpeg to convert WebM files to MP4 format.
- **Send converted files**: Sends the converted MP4 files back to the user via Telegram.
- **File management**: Moves the converted files to a specified directory and cleans up temporary files.

## Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or later)
- [FFmpeg](https://ffmpeg.org/download.html) (ensure it is installed and available in your system's PATH)
- [Telegram Bot API](https://github.com/go-telegram-bot-api/telegram-bot-api)

## Installation

1. **Clone the repository**:

   ```sh
   git clone https://github.com/Burziszcze/telegram-webm-converter.git
   cd telegram-webm-converter
   ```

2. **Install dependencies**:

   ```sh
   go mod tidy
   ```

3. **Set up your Telegram bot**:
   - Create a bot using [BotFather](https://core.telegram.org/bots#botfather) on Telegram and get the API token.

4. **Configure the bot**:
   - Create a `config.toml` file in the root directory and add your bot token:

     ```toml
     [Telegram]
     APIKey = "your_telegram_bot_token"
     ```

5. **Run the bot**:

   ```sh
   go run main.go
   ```

## Usage

- **Send a WebM file**: Send a WebM file directly to the bot via Telegram.
- **Send a URL**: Send a message containing a URL to a WebM file.

The bot will download the file, convert it to MP4, and send the MP4 file back to you.

## Code Overview

### Structure

- `utils/converter.go`: Contains the main logic for handling messages, downloading files, converting them, and sending the results back.
- `main.go`: Entry point for the bot.

### Functions

- `HandleMessage(msg *tgbotapi.Message)`: Handles incoming messages, determines if they contain a WebM file or URL, and processes them.
- `DownloadFile(url string, dest *os.File)`: Downloads a file from a given URL.
- `ConvertWebmToMp4(inputFile string)`: Converts a WebM file to MP4 using FFmpeg.
- `MoveFile(src, dstDir string)`: Moves a file to a specified directory.
- `EnsureTempDir()`: Ensures that the temporary directory exists.
- `CleanupTempFiles()`: Cleans up files in the temporary directory.

## Contributing

Feel free to submit issues, fork the repository, and send pull requests. Contributions are always welcome.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) for the Telegram Bot API implementation in Go.
- [FFmpeg](https://ffmpeg.org/) for the multimedia framework.

## Contact

If you have any questions or feedback, feel free to contact me at [buurzuum@gmail.com].

```

Make sure to update your code to read from the `config.toml` file. Here is an example of how you can modify your `main.go` to read the token from `config.toml`:

```