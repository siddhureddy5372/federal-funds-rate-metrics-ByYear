package services

import (
	"federal-funds-rate-metrics-ByYear/dto"
	"log"
	"strconv"

	"context"

	"federal-funds-rate-metrics-ByYear/db"
)

func IsCurrentYearDataPresent(year int) (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1 FROM federal_funds_insights WHERE year = $1
	)`
	var exists bool
	err := db.Conn.QueryRow(context.Background(), query, year).Scan(&exists)
	if err != nil {
		return false, err
	}
	
	return exists, nil
}

func GetAllYearsData() ([]dto.YearlyInsight, error) {
	query := `SELECT year, average_rate, highest_rate, lowest_rate, growth_percentage, highest_rate_month, lowest_rate_month
				FROM federal_funds_insights
				ORDER BY 
					year DESC; -- Sort the remaining years in descending order
`

	// Execute the query and get a rows iterator
	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed after iteration

	var insights []dto.YearlyInsight

	// Iterate through the result set
	for rows.Next() {
		var insight dto.YearlyInsight
		err := rows.Scan(
			&insight.Year,
			&insight.AverageRate,
			&insight.HighestRate,
			&insight.LowestRate,
			&insight.GrowthPercentage,
			&insight.HighestRateMonth,
			&insight.LowestRateMonth,
		)
		if err != nil {
			return nil, err
		}
		insights = append(insights, insight)
	}

	// Check for errors during iteration
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return insights, nil
}

func ProcessFederalFundsData(data dto.AlphaVantageResponse) ([]dto.YearlyInsight, error) {
	yearlyData := make(map[int][]float64)
	monthlyRates := make(map[int]map[string]float64)

	// Organize data by year and month
	for _, record := range data.Data {
		date := record["date"]
		rateStr := record["value"]

		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			log.Printf("Error parsing rate: %v\n", err)
			continue
		}

		year, month := parseYearMonth(date)
		if year != 0 {
			yearlyData[year] = append(yearlyData[year], rate)

			if monthlyRates[year] == nil {
				monthlyRates[year] = make(map[string]float64)
			}
			monthlyRates[year][month] = rate
		}
	}

	// Calculate insights
	var insights []dto.YearlyInsight
	for year, rates := range yearlyData {
		averageRate := calculateAverage(rates)
		highestRate, lowestRate, highestMonth, lowestMonth := calculateExtremes(monthlyRates[year])

		var growthPercentage float64
		if prevRates, ok := yearlyData[year-1]; ok {
			lastYearAvg := calculateAverage(prevRates)
			growthPercentage = ((averageRate - lastYearAvg) / lastYearAvg) * 100
		}

		insight := dto.YearlyInsight{
			Year:             year,
			AverageRate:      averageRate,
			HighestRate:      highestRate,
			LowestRate:       lowestRate,
			GrowthPercentage: growthPercentage,
			HighestRateMonth: highestMonth,
			LowestRateMonth:  lowestMonth,
		}
		insights = append(insights, insight)
	}

	return insights, nil
}

func parseYearMonth(date string) (int, string) {
	if len(date) < 7 {
		return 0, ""
	}
	year, _ := strconv.Atoi(date[:4])
	month := date[5:7]
	return year, month
}

func calculateAverage(rates []float64) float64 {
	var sum float64
	for _, rate := range rates {
		sum += rate
	}
	return sum / float64(len(rates))
}

func calculateExtremes(monthlyRates map[string]float64) (float64, float64, string, string) {
	var highestRate, lowestRate float64 = -1e9, 1e9
	var highestMonth, lowestMonth string

	for month, rate := range monthlyRates {
		if rate > highestRate {
			highestRate = rate
			highestMonth = month
		}
		if rate < lowestRate {
			lowestRate = rate
			lowestMonth = month
		}
	}

	return highestRate, lowestRate, highestMonth, lowestMonth
}

func StoreFederalFundsInsights(insights []dto.YearlyInsight) error {
	query := `INSERT INTO federal_funds_insights (year, average_rate, highest_rate, lowest_rate, growth_percentage, highest_rate_month, lowest_rate_month)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)
	          ON CONFLICT (year) DO UPDATE
	          SET average_rate = EXCLUDED.average_rate,
	              highest_rate = EXCLUDED.highest_rate,
	              lowest_rate = EXCLUDED.lowest_rate,
	              growth_percentage = EXCLUDED.growth_percentage,
	              highest_rate_month = EXCLUDED.highest_rate_month,
	              lowest_rate_month = EXCLUDED.lowest_rate_month`

	for _, insight := range insights {
		_, err := db.Conn.Exec(context.Background(), query,
			insight.Year, insight.AverageRate, insight.HighestRate, insight.LowestRate, insight.GrowthPercentage, insight.HighestRateMonth, insight.LowestRateMonth)
		if err != nil {
			log.Printf("Failed to insert insight for year %d: %v\n", insight.Year, err)
			return err
		}
	}

	log.Println("Insights stored successfully.")
	return nil
}
