package main

import (
	"encoding/binary"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/abdulmoeid7112/read-exif-cli/utils"

	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	jpg "github.com/dsoprea/go-jpeg-image-structure/v2"
	png "github.com/dsoprea/go-png-image-structure/v2"
)

var (
	pathArg       = ""
	saveAsHTMLArg = false
	helpArg       = false
	csvFilePath   = "output.csv"
)

type Data struct {
	FilePath  string
	Latitude  string
	Longitude string
}

func main() {
	// Define command-line flags
	flag.StringVar(&pathArg, "path", "", "Path of image or directory")
	flag.BoolVar(&saveAsHTMLArg, "html", false, "Save result as HTML")
	flag.BoolVar(&helpArg, "help", false, "Display usage instructions")

	// Parse command-line flags
	flag.Parse()

	if pathArg == "" {
		fmt.Printf("Please provide a file-path for an image.\n")
		os.Exit(1)
	}

	// Display usage instructions if help flag is provided
	if helpArg {
		flag.Usage()
		return
	}

	// Check if input path is valid
	pathType, err := utils.IsPathExists(pathArg)
	if err != nil {
		fmt.Printf("Please provide a valid file-path for an image.\n")
		os.Exit(1)
	}

	data := make([]Data, 0)

	if pathType == utils.File {
		// Check if the file is an image
		if !utils.IsImage(pathArg) {
			os.Exit(1)
		}

		// Read EXIF data from the image
		latitude, longitude, err := readEXIF(pathArg)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		data = append(data, Data{
			FilePath:  pathArg,
			Latitude:  latitude,
			Longitude: longitude,
		})
	} else {
		// Walk through the directory and subdirectories to find images
		err = filepath.Walk(pathArg, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Check if the file is an image
			if utils.IsImage(path) {
				// Read EXIF data from the image
				latitude, longitude, err := readEXIF(path)
				if err != nil {
					log.Printf("Failed to read EXIF data from %s: %v", path, err)
					return nil // Skip to the next file
				}

				data = append(data, Data{
					FilePath:  path,
					Latitude:  latitude,
					Longitude: longitude,
				})
			}

			return nil
		})
		if err != nil {
			log.Fatalf("Error while walking through the directory: %v", err)
		}
	}

	if saveAsHTMLArg {
		if err = saveToHTML(&data); err != nil {
			log.Printf("Error writing HTML file: %v", err)
			os.Exit(1)
		}
	} else {
		if err = saveToCSV(data); err != nil {
			log.Printf("Error writing csv file: %v", err)
			os.Exit(1)
		}
	}
}

// Helper function to read EXIF data from an image file
func readEXIF(path string) (string, string, error) {
	extension := filepath.Ext(path)
	var rootIb *exif.IfdBuilder

	switch extension {
	case ".jpg", ".jpeg":
		// Create new jpeg image parser
		jpgFileParser, err := jpg.NewJpegMediaParser().ParseFile(path)
		if err != nil {
			return "", "", err
		}

		sl := jpgFileParser.(*jpg.SegmentList)
		// Construct exif builder preloaded with all existing tags
		rootIb, err = sl.ConstructExifBuilder()
		if err != nil {
			return "", "", err
		}
	case ".png":
		// Create new png image parser
		pngFileParser, err := png.NewPngMediaParser().ParseFile(path)
		if err != nil {
			return "", "", err
		}

		sl := pngFileParser.(*png.ChunkSlice)
		// Construct exif builder preloaded with all existing tags
		rootIb, err = sl.ConstructExifBuilder()
		if err != nil {
			return "", "", err
		}
	default:
		return "", "", errors.New("invalid extension")
	}

	ifdPath := "IFD/GPSInfo"

	// Get the requested ifdPath
	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPath)
	if err != nil {
		return "", "", err
	}

	// Find the latitude tag using tag latitude id
	latitudeTag, err := ifdIb.FindTag(exif.TagLatitudeId)
	if err != nil {
		return "", "", err
	}

	latitudeInfo, err := getGPSData(latitudeTag.Value().Bytes())
	if err != nil {
		return "", "", err
	}

	// Find the longitude tag using tag longitude id
	longitudeTag, err := ifdIb.FindTag(exif.TagLongitudeId)
	if err != nil {
		return "", "", err
	}

	longitudeInfo, err := getGPSData(longitudeTag.Value().Bytes())
	if err != nil {
		return "", "", err
	}

	return latitudeInfo, longitudeInfo, nil
}

// Helper function to extract GPS position from EXIF data
func getGPSData(data []byte) (string, error) {
	locationInfo := ""

	parser := new(exifcommon.Parser)
	rationals, err := parser.ParseRationals(data, 3, binary.BigEndian)
	if err != nil {
		return "", err
	}

	for _, rational := range rationals {
		locationInfo += fmt.Sprintf("%v ", float64(rational.Numerator)/float64(rational.Denominator))
	}

	return locationInfo, nil
}

func saveToHTML(data *[]Data) error {
	// Open the HTML file for writing
	file, err := os.Create("output.html")
	if err != nil {
		return err
	}
	defer file.Close()

	// Define the HTML template
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
	<title>Image Details</title>
</head>
<body>
	<h1>People Details</h1>
	<table>
		<tr>
			<th>Image File Path</th>
			<th>GPS Latitude</th>
			<th>GPS Longitude</th>
		</tr>
		{{range .}}
		<tr>
			<td>{{.FilePath}}</td>
			<td>{{.Latitude}}</td>
			<td>{{.Longitude}}</td>
		</tr>
		{{end}}
	</table>
</body>
</html>`

	// Create a new template and parse the HTML template string
	tmpl := template.Must(template.New("peopleDetails").Parse(htmlTemplate))

	// Execute the template with the people data and write it to the HTML file
	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	fmt.Println("HTML file created successfully!")

	return nil
}

func saveToCSV(data []Data) error {
	// Create the CSV file
	file, err := os.Create(csvFilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	// Create the CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"Image File Path", "GPS Latitude", "GPS Longitude"}
	err = writer.Write(header)
	if err != nil {
		fmt.Printf("Failed to write csv header.\n")
		os.Exit(1)
	}

	// Write GPS info to csv
	for _, fileData := range data {
		row := []string{fileData.FilePath, fileData.Latitude, fileData.Longitude}
		if err = writer.Write(row); err != nil {
			log.Printf("Error writing CSV row: %v", err)
		}
	}

	fmt.Println("CSV file created successfully!")

	return nil
}
