package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const uploadPath = "./uploads/"
const gifOutputPath = "./gifs/"

func main() {
	// Ensure directories exist
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating upload directory: %v", err)
	}
	err = os.MkdirAll(gifOutputPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating gif directory: %v", err)
	}

	// Set up HTTP handlers
	http.HandleFunc("/", uploadFormHandler)
	http.HandleFunc("/upload", uploadFileHandler)

	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// uploadFormHandler renders a simple HTML form to upload a video
func uploadFormHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<html>
		<head>
      <link rel="icon" type="image/png" href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAKAAAACgCAYAAACLz2ctAABJ30lEQVR4AezSAQ2AMBDAwMrBv8HHBdvYNTkHbVXS1NMUy9CAATEgGBADggExIBgQA4IBMSAYEAOCATEgGBADggExIBgQA4IBMSAYEAOCATEgBBrRLMLuMMGYVmJesHFsWkAAAAASUVORK5CYII=">
			<title>Video to GIF Converter</title>
			<script src="https://cdn.tailwindcss.com"></script>
		</head>
		<body class="bg-gray-900 text-white flex items-center justify-center min-h-screen">
			<div class="p-6 bg-gray-800 bg-opacity-60 rounded-lg shadow-lg text-center">
				<h1 class="text-4xl font-bold mb-4">Upload Video File</h1>
				<form enctype="multipart/form-data" action="/upload" method="post" class="space-y-4">
					<input class="block w-full text-white bg-gray-700 rounded-md p-2" type="file" name="video" accept="video/*">
					<br>
					<input class="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md cursor-pointer" type="submit" value="Convert to GIF">
				</form>
			</div>
		</body>
	</html>`
	w.Write([]byte(html))
}

// uploadFileHandler handles the video file upload, converts to GIF, and returns the result as a base64 encoded image
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Limit to 10MB files
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusInternalServerError)
		return
	}

	// Retrieve the uploaded file
	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the uploaded file to the server
	videoFileName := filepath.Join(uploadPath, header.Filename)
	outFile, err := os.Create(videoFileName)
	if err != nil {
		http.Error(w, "Could not save uploaded file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Convert the video to GIF
	gifFileName := getOutputGIFPath(videoFileName)
	err = convertVideoToGIF(videoFileName, gifFileName)
	if err != nil {
		http.Error(w, "Failed to convert video to GIF", http.StatusInternalServerError)
		return
	}

	// Read the generated GIF and encode it as base64
	gifData, err := os.ReadFile(gifFileName)
	if err != nil {
		http.Error(w, "Failed to read GIF file", http.StatusInternalServerError)
		return
	}
	base64GIF := base64.StdEncoding.EncodeToString(gifData)

	// Serve HTML with embedded base64 image and a back button
	html := fmt.Sprintf(`
		<html>
			<head>
				<title>GIF Result</title>
				<script src="https://cdn.tailwindcss.com"></script>
			</head>
			<body class="bg-gray-900 text-white flex flex-col items-center justify-center min-h-screen">
				<h1 class="text-4xl font-bold mb-4">GIF Result</h1>
				<img src="data:image/gif;base64,%s" alt="Generated GIF" class="border-4 border-blue-600 rounded-md mb-4">
				<a href="data:image/gif;base64,%s" download="output.gif" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md cursor-pointer mb-4">Download GIF</a>
				<a href="/" class="px-4 py-2 bg-gray-600 hover:bg-gray-700 rounded-md cursor-pointer">Back to Upload</a>
			</body>
		</html>`, base64GIF, base64GIF)

	w.Write([]byte(html))

	// Schedule cleanup to delete both the video and GIF files after serving
	time.AfterFunc(10*time.Second, func() {
		os.Remove(videoFileName) // Remove video
		os.Remove(gifFileName)   // Remove GIF
	})
}

// convertVideoToGIF converts a video file to a GIF using ffmpeg
func convertVideoToGIF(inputVideo, outputGIF string) error {
	// Check if input file exists
	if _, err := os.Stat(inputVideo); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist")
	}

	// Construct the ffmpeg command
	cmd := exec.Command("ffmpeg", "-i", inputVideo, "-vf", "fps=10,scale=320:-1:flags=lanczos", "-t", "10", outputGIF)

	// Set the output destination
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert video to GIF: %v", err)
	}

	return nil
}

// getOutputGIFPath generates a valid output GIF file path
func getOutputGIFPath(inputVideo string) string {
	ext := filepath.Ext(inputVideo)
	outputGIF := strings.TrimSuffix(filepath.Base(inputVideo), ext) + ".gif"
	return filepath.Join(gifOutputPath, outputGIF)
}
