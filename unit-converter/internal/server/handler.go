package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unit-converter/cmd/web"
)

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the converter type from URL query parameter
	convertType := r.URL.Query().Get("converter")

	// If we're at root path with no converter parameter, redirect to length converter
	if r.URL.Path == "/" && convertType == "" {
		http.Redirect(w, r, "/?converter=length", http.StatusSeeOther)
		return
	}

	err := web.Base(convertType).Render(ctx, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) WeightHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	weightStr := r.FormValue("weight")
	from := r.FormValue("from")
	to := r.FormValue("to")

	// Convert weight string to float
	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		http.Error(w, "Invalid weight value", http.StatusBadRequest)
		return
	}

	// Convert to base unit (kg) first
	var baseWeight float64
	switch strings.ToLower(from) {
	case "kg":
		baseWeight = weight
	case "g":
		baseWeight = weight / 1000
	case "mg":
		baseWeight = weight / 1000000
	case "lb":
		baseWeight = weight * 0.45359237
	case "oz":
		baseWeight = weight * 0.0283495231
	default:
		http.Error(w, "Invalid source unit", http.StatusBadRequest)
		return
	}

	// Convert from base unit to target unit
	var result float64
	switch strings.ToLower(to) {
	case "kg":
		result = baseWeight
	case "g":
		result = baseWeight * 1000
	case "mg":
		result = baseWeight * 1000000
	case "lb":
		result = baseWeight / 0.45359237
	case "oz":
		result = baseWeight / 0.0283495231
	default:
		http.Error(w, "Invalid target unit", http.StatusBadRequest)
		return
	}

	// Create result HTML matching the exact format
	resultHTML := fmt.Sprintf(`
    <div id="result" class="px-8 py-4">
        <div >
            <p>
                %.2f %s = %.2f %s
            </p>
        </div>

    </div>
    `, weight, from, result, to)

	// Set content type and write response
	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write([]byte(resultHTML))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (s *Server) TempHandler(w http.ResponseWriter, r *http.Request) {

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	tempStr := r.FormValue("temperature")
	from := r.FormValue("from")
	to := r.FormValue("to")

	// Convert temperature string to float
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		http.Error(w, "Invalid temperature value", http.StatusBadRequest)
		return
	}

	// Convert to base unit (Celsius) first
	var baseTemp float64
	switch strings.ToLower(from) {
	case "c":
		baseTemp = temp
	case "f":
		baseTemp = (temp - 32) * 5 / 9
	case "k":
		baseTemp = temp - 273.15
	default:
		http.Error(w, "Invalid source unit", http.StatusBadRequest)
	}

	// Convert from base unit to target unit
	var result float64
	switch strings.ToLower(to) {
	case "c":
		result = baseTemp
	case "f":
		result = (baseTemp * 9 / 5) + 32
	case "k":
		result = baseTemp + 273.15
	default:
		http.Error(w, "Invalid target unit", http.StatusBadRequest)
	}

	// Create result HTML matching the exact format
	resultHTML := fmt.Sprintf(`
	<div id="result" class="px-8 py-4">
		<div >
			<p>	
				%.2f %s = %.2f %s
			</p>	
		</div>	
	</div>	
	`, temp, from, result, to)

	// Set content type and write response
	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write([]byte(resultHTML))
}

func (s *Server) LengthHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	lengthStr := r.FormValue("length")
	from := r.FormValue("from")
	to := r.FormValue("to")

	// Convert length string to float
	length, err := strconv.ParseFloat(lengthStr, 64)
	if err != nil {

		http.Error(w, "Invalid length value", http.StatusBadRequest)
		return
	}

	// Convert to base unit (meter) first
	var baseLength float64
	switch strings.ToLower(from) {
	case "m":
		baseLength = length
	case "cm":
		baseLength = length / 100
	case "mm":
		baseLength = length / 1000
	case "km":
		baseLength = length * 1000
	case "ft":
		baseLength = length * 0.3048
	default:
		http.Error(w, "Invalid source unit", http.StatusBadRequest)
	}

	// Convert from base unit to target unit

	var result float64
	switch strings.ToLower(to) {
	case "m":
		result = baseLength
	case "cm":
		result = baseLength * 100
	case "mm":
		result = baseLength * 1000
	case "km":
		result = baseLength / 1000
	case "ft":
		result = baseLength / 0.3048
	default:
		http.Error(w, "Invalid target unit", http.StatusBadRequest)
	}

	// Create result HTML matching the exact format
	resultHTML := fmt.Sprintf(`
	<div id="result" class="px-8 py-4">
		<div >
			<p>

				%.2f %s = %.2f %s
			</p>

		</div>
	</div>
	`, length, from, result, to)

	// Set content type and write response
	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write([]byte(resultHTML))
}
