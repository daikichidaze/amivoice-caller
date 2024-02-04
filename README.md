# Audio File Converter and API Uploader

This Go module converts audio files to MP3 format and uploads them to a specified API endpoint, demonstrating how to work with files, external commands, and HTTP requests in Go.

## Features

- Converts `.m4a` or any compatible audio file format to `.mp3` using `ffmpeg`.
- Uploads the converted `.mp3` file to a specified API endpoint using multipart/form-data.
- Dynamically loads API keys from a local file.
- Saves API response to a local text file for further processing or review.

## Requirements

- Go 1.15 or later.
- `ffmpeg` must be installed on your system and accessible from the command line.

## Amivoice

https://acp.amivoice.com/

## Installation

Clone this repository or download the source code:

```bash
git clone git@github.com:daikichidaze/amivoice-caller.git
```

Ensure `ffmpeg` is installed on your system. On most Unix systems, you can install `ffmpeg` via your package manager. For example, on Ubuntu:

```bash
sudo apt-get install ffmpeg
```

## Usage

To use this module, run the Go script with the path to the audio file you wish to convert and upload as the first argument:

```bash
go run main.go <path to your audio file>
```

Make sure you have an `APIKEY` file in the same directory as your script, containing your API key.

### Example

```bash
go run main.go /path/to/your/audiofile.m4a
```

## How It Works

1. **Conversion to MP3**: The script uses `ffmpeg` to convert the input audio file to MP3 format.
2. **API Key Loading**: Reads the API key from a local file named `APIKEY`.
3. **Creating a Multi-Part Request**: Prepares a multi-part/form-data request with the converted MP3 file and the loaded API key.
4. **Sending the Request**: Uploads the file to the API endpoint and handles the response.
5. **Saving the Response**: The API response is saved to `response.txt` in the same directory for further inspection.
