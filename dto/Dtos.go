package dto

import (
	"federal-funds-rate-metrics-ByYear/models"
)

// Message represents a response structure with a status, optional message, 
// and optional data payload containing user information.
type Message struct {
	Status  string   `json:"status"`  // Status of the response (e.g., success, error).
	Message string   `json:"message,omitempty"` // Optional message providing additional details.
	Data    *UserDto `json:"data,omitempty"`    // Optional user data payload.
}

// MessageInsights represents a response structure containing insights data.
// This is used to convey yearly insights with additional status and message.
type MessageInsights struct {
	Status  string         `json:"status"`          // Status of the response (e.g., success, error).
	Message string         `json:"message,omitempty"` // Optional message providing additional details.
	Data    []YearlyInsight `json:"data,omitempty"`    // Array of yearly insights data.
}


// YearlyInsight represents a structure to hold insights about a specific year.
// Includes financial rates, growth percentage, and associated months.
type YearlyInsight struct {
	Year              int     // The year the insights pertain to.
	AverageRate       float64 // Average rate for the year.
	HighestRate       float64 // Highest rate recorded in the year.
	LowestRate        float64 // Lowest rate recorded in the year.
	GrowthPercentage  float64 // Growth percentage for the year.
	HighestRateMonth  string  // Month with the highest rate.
	LowestRateMonth   string  // Month with the lowest rate.
}

// UserDto represents a simplified structure for user information to be shared in responses.
type UserDto struct {
	Name  string `json:"name"`  // User's name.
	Email string `json:"email"` // User's email address.
}

// AlphaVantageResponse represents the structure for responses from the Alpha Vantage API.
// Includes metadata and an array of data entries.
type AlphaVantageResponse struct {
	MetaData map[string]string   `json:"Meta Data"` // Metadata about the API response.
	Data     []map[string]string `json:"Data"`      // Array of data entries, where each entry is a map of strings.
}

// ConvertToUserDto converts a user model object into a UserDto object.
// This function simplifies the user data for external use.
func ConvertToUserDto(user models.User) UserDto {
	return UserDto{
		Name:  user.Name,   // Assign the user's name from the models.User structure.
		Email: user.Email,  // Assign the user's email from the models.User structure.
	}
}
