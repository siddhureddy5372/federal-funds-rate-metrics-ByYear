package handle

import (
	"encoding/json"
	"net/http"
	"strings"

	"federal-funds-rate-metrics-ByYear/dto"
	"federal-funds-rate-metrics-ByYear/services"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// CreateUser handles POST requests to create a new user.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestDto dto.UserDto
	var response dto.Message
	err := json.NewDecoder(r.Body).Decode(&requestDto)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate DTO
	err = validate.Struct(requestDto)
	if err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	savedDto, err := services.CreateUser(&requestDto)
	if err != nil {
		http.Error(w, "Creation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	response = dto.Message{Status: "success", Data: savedDto}
	w.WriteHeader(http.StatusCreated)
	jsonResponse(w, response)
}

// UserInfo handles GET requests to retrieve user information.
// It expects the URL in the form "/auth/<email>".
func UserInfo(w http.ResponseWriter, r *http.Request) {
	// Extract the email from the URL path.
	// For example, if the URL is "/auth/siddhureddy", then email will be "siddhureddy".
	const prefix = "/auth/"
	email := ""
	if strings.HasPrefix(r.URL.Path, prefix) {
		email = r.URL.Path[len(prefix):]
	}

	var response dto.Message

	if email == "" {
		// Respond with an error message if email is not provided.
		response = dto.Message{
			Status:  "fail",
			Message: "Invalid email. Please provide a valid email.",
		}
	} else {
		// Fetch user details.
		details, err := services.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "Retrieval error: "+err.Error(), http.StatusBadRequest)
			return
		} else if details.ID == 0 {
			// User not found.
			response = dto.Message{
				Status:  "fail",
				Message: "User not found",
			}
		} else {
			// Success: convert details to DTO.
			convertedDto := dto.ConvertToUserDto(details)
			response = dto.Message{
				Status: "success",
				Data:   &convertedDto,
			}
		}
	}
	jsonResponse(w, response)
}

// HealthCheck handles GET requests to check server health.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// jsonResponse is a utility function to send JSON responses.
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
