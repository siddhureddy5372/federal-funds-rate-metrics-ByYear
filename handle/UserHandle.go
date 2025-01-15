package handle

import (
	"encoding/json"
	"net/http"
	"strconv"

	"federal-funds-rate-metrics-ByYear/dto"
	"federal-funds-rate-metrics-ByYear/services"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate = validator.New()

// Data structure for JSON response

// Handlers

// POST Handlers

var Message dto.Message

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
	} else {
		savedDto, err := services.CreateUser(&requestDto)
		if err != nil {
			http.Error(w, "Creation error: "+err.Error(), http.StatusBadRequest)
			return
		} else {
			response = dto.Message{Status: "success", Data: savedDto}
			w.WriteHeader(http.StatusCreated)
			jsonResponse(w, response)
		}
	}
}

// GET Handlers
// User Information by ID
func UserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	num := vars["id"]

	// Convert ID from string to integer
	id, err := strconv.Atoi(num)
	var response dto.Message

	if err != nil {
		// Respond with an error message
		response = dto.Message{
			Status:  "fail",
			Message: "Invalid ID format",
		}
	} else {
		// Fetch user details
		details, err := services.GetUserByID(id)
		if err != nil {
			http.Error(w, "Retrieval error: "+err.Error(), http.StatusBadRequest)
			return
		} else if details.ID == 0 {
			// User not found case
			response = dto.Message{
				Status:  "fail",
				Message: "User not found",
			}
		} else {
			// Success case
			convertedDto := dto.ConvertToUserDto(details)
			response = dto.Message{
				Status: "success",
				Data:   &convertedDto,
			}
		}
	}
	// Send JSON response
	jsonResponse(w, response)
}

// Utility function to send JSON response
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
