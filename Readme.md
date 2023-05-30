# EXIF to CSV/HTML Converter

This is a command-line utility written in Go that extracts EXIF data from images and writes the attributes to a CSV file. It utilizes the [go-exif](https://github.com/dsoprea/go-exif) library and its supporting following libraries

- [go-jpeg-image-structure](https://github.com/dsoprea/go-jpeg-image-structure)
- [go-png-image-structure](github.com/dsoprea/go-png-image-structure)

for parsing image metadata.

## Usage

To use the EXIF to CSV Converter, follow the instructions below:

1. Download and install Go(1.16 or higher) if you haven't already: https://golang.org/dl/
2. Install the required dependencies by running the following command:
    - go install github.com/dsoprea/go-exif github.com/dsoprea/go-jpeg-image-structure/v2 github.com/dsoprea/go-png-image-structure/v2
3. Build the executable by running the following command:
    - go build -o bin/read-exif-cli.exe .\cmd\main.go
4. Run the executable with the desired flags:
    - ./read-exif-cli --path /path/to/images --html
        - Use the `--path` flag to specify the path to the image file or the directory containing multiple images.
        - Use the `--html` flag to store the results in HTML format instead of the default CSV format.
        - Use the `--help` flag to display the usage instructions.

The utility will process the images, extract the relevant EXIF data (such as GPS position), and store the attributes in a CSV file (or HTML file if the `--html` flag is used).

## Examples

Here are a few examples of using the EXIF to CSV Converter:

1. Convert a single image to CSV:
    - ./read-exif-cli --path /path/to/image.jpg
2. Convert multiple images in a directory to CSV:
    - ./read-exif-cli --path /path/to/images
3. Convert images to HTML format:
    - ./read-exif-cli --path /path/to/images --html


