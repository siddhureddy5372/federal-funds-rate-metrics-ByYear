package handle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"federal-funds-rate-metrics-ByYear/config"
	"federal-funds-rate-metrics-ByYear/dto"
	"federal-funds-rate-metrics-ByYear/services"
)

// FederalFundsHandlerInsight handles requests for federal funds insights, either retrieving them from the database
// or fetching, processing, and storing them if not already present.
func FederalFundsHandlerInsight(w http.ResponseWriter, r *http.Request) {
	currentYear := time.Now().Year() - 1

	// Check if current year's data exists in the database
	exists, err := services.IsCurrentYearDataPresent(currentYear)
	if err != nil {
		handleError(w, fmt.Sprintf("Error checking current year data: %v", err))
		return
	}

	if exists {
		// Retrieve all years' data from the database
		insights, err := services.GetAllYearsData()
		if err != nil {
			handleError(w, fmt.Sprintf("Error retrieving data: %v", err))
			return
		}

		respondWithJSON(w, http.StatusOK, dto.MessageInsights{
			Status:  "success",
			Message: "Data successfully fetched from the database.",
			Data:    insights,
		})
		return
	}

	// Fetch data from the external API
	data, err := fetchFederalFundsRate()
	if err != nil {
		handleError(w, fmt.Sprintf("Error fetching data: %v", err))
		return
	}

	// Process the data to calculate insights
	insights, err := services.ProcessFederalFundsData(data)
	if err != nil {
		handleError(w, fmt.Sprintf("Error processing data: %v", err))
		return
	}

	// Store the insights in the database
	err = services.StoreFederalFundsInsights(insights)
	if err != nil {
		handleError(w, fmt.Sprintf("Error storing insights: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, dto.MessageInsights{
		Status:  "success",
		Message: "Data successfully fetched, processed, and stored.",
		Data:    insights,
	})
}

// FederalFundsHandler handles simple requests to fetch federal funds data directly from the API.
func FederalFundsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchFederalFundsRate()
	if err != nil {
		handleError(w, fmt.Sprintf("Error fetching federal funds rate: %v", err))
		return
	}

	// Process the data to calculate insights
	insights, err := services.ProcessFederalFundsData(data)
	if err != nil {
		handleError(w, fmt.Sprintf("Error processing data: %v", err))
		return
	}

	// Sort the insights slice in descending order by Year
	sort.Slice(insights, func(i, j int) bool {
		return insights[i].Year > insights[j].Year
	})

	// Respond with sorted insights
	respondWithJSON(w, http.StatusOK, dto.MessageInsights{
		Status:  "success",
		Message: "Fetched federal funds rate data successfully.",
		Data:    insights,
	})
}

// fetchFederalFundsRate fetches federal funds rate data from the Alpha Vantage API.
func fetchFederalFundsRate() (dto.AlphaVantageResponse, error) {
	if appConfig == nil {
		return dto.AlphaVantageResponse{}, fmt.Errorf("config not initialized")
	}

	apiURL := fmt.Sprintf("https://www.alphavantage.co/query?function=FEDERAL_FUNDS_RATE&interval=monthly&apikey=%s", appConfig.APIKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		return dto.AlphaVantageResponse{}, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.AlphaVantageResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.AlphaVantageResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var data dto.AlphaVantageResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return dto.AlphaVantageResponse{}, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return data, nil
}

// handleError is a helper function to send error responses to the client.
func handleError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}

// respondWithJSON is a helper function to send JSON responses.
func respondWithJSON(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

var appConfig *config.Config

// InitConfig initializes the configuration for the handler package.
func InitConfig(config *config.Config) {
	appConfig = config
}
