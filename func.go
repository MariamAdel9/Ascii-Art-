package main

import (
	"bufio"
	"net/http"
	"os"
	"slices"
	"strings"
)

type ERROR struct {
	Message string
}

func errorcustomize(w http.ResponseWriter, stat int, Message string) {

	data := ERROR{Message: Message}
	w.WriteHeader(stat)
	templates.ExecuteTemplate(w, "error.html", data)
}

// Handle GET and POST requests for ASCII art generation
func Post(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	// if r.URL.Path != "/asciiart" {
	// 		errorcustomize(w, http.StatusNotFound, "404 PAGE NOT FOUND")
	// 		return
	// 	}
		templates.ExecuteTemplate(w, "index.html", nil)

	case "POST":

		
		text := r.FormValue("text")
		if len(text) > 300 {
			errorcustomize(w, http.StatusBadRequest, "ERROR-400\nText too long! Maximum length is 300 characters.")
			return
		}

		if len(text) < 1 {
			errorcustomize(w, http.StatusBadRequest, "ERROR-400\nText Please enter your text!.")
			return
		}

		banner := r.FormValue("banner")
		finalText := MultiLines(text)

		arr := []string{"standard", "shadow", "thinkertoy"}
		if !slices.Contains(arr, banner) {
			errorcustomize(w, http.StatusBadRequest, "ERROR-400\nPlease select a valid font!")
			return
		}

		if !isValid(finalText) {
			errorcustomize(w, http.StatusBadRequest, "ERROR-400\nBad request!")
			return
		}

		fontFile, err := os.Open(banner + ".txt")
		if err != nil {
			errorcustomize(w, http.StatusInternalServerError, "ERROR-500\nInternal Server Error")
			return
		}
		defer fontFile.Close()

		lines, err := readLines(fontFile)
		if err != nil {

			errorcustomize(w, http.StatusInternalServerError, "ERROR-500\nInternal Server Error")
			return
		}

		char := make(map[int][]string)
		dec := 31
		for _, line := range lines {
			if line == "" {
				dec++
			} else {
				char[dec] = append(char[dec], line)
			}
		}

		// Generate ASCII art
		output := ""
		for _, line := range strings.Split(finalText, "\n") {
			output += PrintArt(line, char) + "\n"
		}

		if len(output) < 0 {
			errorcustomize(w, http.StatusInternalServerError, "ERROR-500\nInternal Server Error")
			return
		} else {
			err := templates.ExecuteTemplate(w, "result.html", output)
			if err != nil {
				errorcustomize(w, http.StatusInternalServerError, "ERROR-500\nInternal Server Error")
			}
			
		}
	default:
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}

func IndexHandler (w http.ResponseWriter, r *http.Request){
    if r.URL.Path == "/" {
      
        err := templates.ExecuteTemplate(w, "index.html", nil)

        if err != nil  {
			errorcustomize(w, http.StatusInternalServerError, "ERROR-500\nInternal Server Error")
        }

    } else {
		if _, err := os.Stat("./templates" + r.URL.Path); os.IsNotExist(err) {
            errorcustomize(w, http.StatusNotFound, "404 PAGE NOT FOUND")
            return
        }
        http.ServeFile(w, r, "./templates"+r.URL.Path)

    }
}

// Normalize newlines in input text
func MultiLines(text string) string {
	return strings.ReplaceAll(text, "\r\n", "\n")
}

// Check if the text contains valid printable ASCII characters
func isValid(text string) bool {
	for _, v := range text {
		if !(v >= 32 && v <= 126) && v != '\n' {
			return false
		}
	}
	return true
}

// Read all lines from the font file
func readLines(file *os.File) ([]string, error) {
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Generate ASCII art for the input string
func PrintArt(str string, m map[int][]string) string {
	if str == "" {
		return ""
	}

	result := make([]string, len(m[32])) // Assume space character has an entry

	for _, r := range str {
		if r == '\n' {
			result = append(result, "\n")
			continue
		}

		if artLines, exists := m[int(r)]; exists {
			for i := 0; i < len(artLines); i++ {
				result[i] += artLines[i]
			}
		}
	}

	return strings.Join(result, "\n")
}
