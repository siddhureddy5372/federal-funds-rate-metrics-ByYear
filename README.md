# User and Metrics Management API

## Overview

This project is a RESTful API that provides user management functionalities such as registering users and retrieving user information. Additionally, the API displays metrics related to financial data, fetched and processed from the free **Alpha Vantage API**. If the financial data for the requested year(s) is unavailable in the database, the API fetches it from Alpha Vantage, processes it, and stores it in the database for future use.

---

## Key Features

### User Management
1. **Register Users**  
   - Endpoint: `POST /create`  
   - Description: Allows new users to register with their information.

2. **Get User Information**  
   - Endpoint: `GET /id/{id}`  
   - Description: Retrieves the information of a user based on their unique ID.

---

### Financial Metrics
1. **Homepage Metrics**  
   - Endpoint: `GET /`  
   - Description: Displays financial metrics for different years. An example response:  

     ```json
     [{
        "status": "success",
        "message": ".........",
        "data":[ 
                {
                    "Year": 2024,
                    "AverageRate": 5.143333333333333,
                    "HighestRate": 5.33,
                    "LowestRate": 4.48,
                    "GrowthPercentage": 2.3718692983910916,
                    "HighestRateMonth": "05",
                    "LowestRateMonth": "12"
                },
                ....
                ]
    } ]
     ```

   - **How it works**:  
     - On each homepage request, the API checks if the metrics for the requested years are stored in the database.
     - If not found:  
       - Fetches the data from the **Alpha Vantage API**.
       - Processes the data into meaningful metrics.
       - Stores the results in the database for future use.

2. **Data Source**:  
   - The financial data is sourced from the free **Alpha Vantage API** using a secret API token.

---

## Tech Stack

1. **Programming Language**: Go (Golang)
2. **Database**: PostgreSQL (or any preferred database backend)
3. **Environment Configuration**: `.env` file is used to securely store sensitive data like API keys and database credentials.

---

## Environment Variables

The following environment variables should be configured in a `.env` file:

```plaintext
DATABASE_URL = "YOUR_DATABASE_URL"
API_KEY = "YOUR_API_TOKEN"

```

---

## Setup and Installation

1. **Clone the Repository**
   ```bash
   git clone https://github.com/siddhureddy5372/federal-funds-rate-metrics-ByYear.git
   cd federal-funds-rate-metrics-ByYear
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Setup Environment Variables**
   - Create a `.env` file in the root of your project.
   - Add the environment variables listed above.

4. **Run the Application**
   ```bash
   go run main.go
   ```

5. **Access the API**
   - Base URL: `http://localhost:8080` (default)

---

## API Endpoints

### User Management
| Method | Endpoint          | Description                     |
|--------|-------------------|---------------------------------|
| POST   | `/create`         | Register a new user            |
| GET    | `/id/{id}`        | Get information of a user by ID|

### Financial Metrics
| Method | Endpoint          | Description                     |
|--------|-------------------|---------------------------------|
| GET    | `/`               | Get financial metrics by year   |

---

## Credits

This project uses the **Alpha Vantage API** for financial data. All financial data belongs to Alpha Vantage. For more information, visit [Alpha Vantage](https://www.alphavantage.co/).

---

## License

This project is open-source and available under the [MIT License](LICENSE).
