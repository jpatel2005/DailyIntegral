package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/time/rate"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

var (
	WOLFRAM_URL    string
	AUTH0_DOMAIN   string
	AUTH0_AUDIENCE string
	TOKEN_SCOPES   string
	API_KEY        string
	db             *sql.DB
	DATE_LAYOUT    = "2006-01-02"
	DAY_SECONDS    = 86400.0
	// Create a rate limiter per client
	clients = make(map[string]*rate.Limiter)
	mu      sync.Mutex
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Audience []string `json:"aud"`
	Scope    string   `json:"scope"`
}

// LeaderboardCursor is for the leaderboard cursor pagination
type LeaderboardCursor struct {
	DailyProblems int    `json:"daily_problems"`
	TotalProblems int    `json:"total_problems"`
	Username      string `json:"username"`
}

// LeaderboardEntry contains data for a single leaderboard user
type LeaderboardEntry struct {
	Username      string `json:"username"`
	DailyProblems int    `json:"daily_problems"`
	TotalProblems int    `json:"total_problems"`
	CreationDate  string `json:"creation_date"`
}

// Validate does nothing for this example, but we need it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	scopes := strings.Split(TOKEN_SCOPES, " ")
	for i := range scopes {
		if !strings.Contains(c.Scope, scopes[i]) {
			return fmt.Errorf("invalid token scopes")
		}
	}
	return nil
}

func validateDate(dateStr string) (time.Time, string) {
	date, err := time.Parse(DATE_LAYOUT, dateStr)
	if err != nil {
		return time.Time{}, "Invalid date format, expected YYYY-MM-DD"
	}
	return date, ""
}

type UserInfo struct {
	Username string `json:"username"`
}

func fetchUserInfo(token string) (*UserInfo, error) {
	userInfoURL := fmt.Sprintf("https://%s/userinfo", AUTH0_DOMAIN)
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user info (invalid token provided): %s", resp.Status)
	}
	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

/*
This function should ensure that the following things are completed:
- Perform standard JWT validation
- Verify token audience claims (via validator configuration)
- Verify token scope configuration (via CustomClaims Validate function)
- Validate that the provided username is associated with the given access token
*/
func validateToken(token string, checkUsername bool, username string) bool {
	issuerURL, err := url.Parse("https://" + AUTH0_DOMAIN + "/")
	if err != nil {
		log.Printf("Failed to parse the issuer url: %v", err)
		return false
	}
	// Go through the process of validating the provided JWT access token
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		// ensure that audience is validated
		[]string{AUTH0_AUDIENCE},
		// will check for proper scopes by calling the Validate function
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Printf("Failed to set up the jwt validator")
		return false
	}
	// Finally attempt to validate the jwt token
	_, err = jwtValidator.ValidateToken(context.TODO(), token)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		return false
	}
	// Exit early if username validation not being performed
	if !checkUsername {
		return true
	}
	// Fetch user info from Auth0
	userInfo, err := fetchUserInfo(token)
	if err != nil {
		log.Printf("Failed to fetch user info when validating token: %v", err)
		return false
	}
	if userInfo.Username != username {
		return false
	}
	return true
}

func getLimiter(ip string, cooldownTime int, burstCount int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := clients[ip]
	if !exists {
		// Allow 1 request every "cooldownTime" seconds (with a burstCount as well)
		limiter = rate.NewLimiter(rate.Every(time.Duration(cooldownTime)*time.Second), burstCount)
		clients[ip] = limiter
	}

	return limiter
}

func verifierRateLimiter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		limiter := getLimiter(ip, 15, 1)

		if !limiter.Allow() {
			return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "Too many requests"})
		}

		return next(c)
	}
}

func generalRateLimiter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		limiter := getLimiter(ip, 1, 25)

		if !limiter.Allow() {
			return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "Too many requests"})
		}

		return next(c)
	}
}

func retrieveProblem(c echo.Context) error {
	// will not require access token
	dateStr := c.Param("date")
	inputTime, msg := validateDate(dateStr)
	if len(msg) != 0 {
		return c.JSON(http.StatusBadRequest, msg)
	}
	// If the requested problem is from the future, then exit early (user shouldn't be able to view this problem)
	currTime := time.Now().UTC()
	if inputTime.After(currTime) {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("cannot request problem data from the future: %s", dateStr))
	}
	rows, err := db.Query(`
		SELECT problem
		FROM problems
		WHERE date_used = $1
	`, dateStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to query problem data for %s: %v", dateStr, err))
	}
	defer rows.Close()
	// Check if at least one row is returned
	if !rows.Next() {
		return c.JSON(http.StatusNoContent, fmt.Sprintf("no data found for %s", dateStr))
	}
	// create struct to hold the problem data
	var problemData struct {
		Problem string `json:"problem"`
	}
	// Retrieve data from the first row
	if err := rows.Scan(&problemData.Problem); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to scan problem data from database for %s: %v", dateStr, err))
	}
	// Check if there is more than one row (this shouldn't happen)
	if rows.Next() {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("expected one result row, got more than one for %s", dateStr))
	}
	// Return the problem data as JSON
	return c.JSON(http.StatusOK, problemData)
}

func retrieveSolutionAndSteps(c echo.Context) error {
	// require access token
	token := c.Request().Header.Get("Authorization")
	if !validateToken(token, false, "") {
		return c.JSON(http.StatusUnauthorized, "Invalid token provided")
	}
	dateStr := c.Param("date")
	inputTime, msg := validateDate(dateStr)
	if len(msg) != 0 {
		return c.JSON(http.StatusBadRequest, msg)
	}
	// If the requested problem is from the today or from future, then exit early (user shouldn't be able to view this problem)
	currTime := time.Now().UTC()
	diffSeconds := currTime.Sub(inputTime).Seconds()
	if diffSeconds < DAY_SECONDS {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("cannot request solution data from the current date or future: %s", dateStr))
	}
	rows, err := db.Query(`
		SELECT answer, steps
		FROM problems
		WHERE date_used = $1
	`, dateStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to query steps data for %s: %v", dateStr, err))
	}
	defer rows.Close()
	// Check if at least one row is returned
	if !rows.Next() {
		return c.JSON(http.StatusNoContent, fmt.Sprintf("no steps data found for %s", dateStr))
	}
	// create struct to hold the problem data
	var problemData struct {
		Answer string `json:"answer"`
		Steps  string `json:"steps"`
	}
	// Retrieve data from the first row
	if err := rows.Scan(&problemData.Answer, &problemData.Steps); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to scan steps data from database for %s: %v", dateStr, err))
	}
	// Check if there is more than one row (this shouldn't happen)
	if rows.Next() {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("expected one steps result row, got more than one for %s", dateStr))
	}
	// Return the solution and steps data as JSON
	return c.JSON(http.StatusOK, problemData)
}

func verifySolution(c echo.Context) error {
	// require access token
	username := c.FormValue("username")
	if len(username) == 0 {
		return c.JSON(http.StatusBadRequest, "no username provided for solution verify")
	}
	row := db.QueryRow(`
		SELECT 1
		FROM users
		WHERE username = $1
	`, username)
	if err := row.Scan(new(int)); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, "Invalid username provided")
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to validate username for %s: %v", username, err))
	}
	token := c.Request().Header.Get("Authorization")
	if !validateToken(token, true, username) {
		return c.JSON(http.StatusUnauthorized, "Invalid token provided")
	}
	dateStr := c.Param("date")
	inputTime, msg := validateDate(dateStr)
	if len(msg) != 0 {
		return c.JSON(http.StatusBadRequest, msg)
	}
	// If the requested problem is from the future, then exit early (user shouldn't be able to view this problem)
	submitTime := time.Now().UTC()
	if inputTime.After(submitTime) {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("cannot verify solution from the future: %s", dateStr))
	}
	answer := c.FormValue("answer")
	// First check to see that the answer doesn't contain integral or isn't empty or doesn't contain comparison operators (avoid cheating or wasteful processing)
	if len(answer) == 0 || strings.Contains(answer, "\\int") || strings.Contains(answer, ">") || strings.Contains(answer, "<") {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Invalid answer provided. Cannot provide empty string or \\int: %s or use comparison operators", answer))
	}
	rows, err := db.Query(`
		SELECT answer
		FROM problems
		WHERE date_used = $1
	`, dateStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to query data for %s: %v", dateStr, err))
	}
	defer rows.Close()
	// Check if at least one row is returned
	if !rows.Next() {
		return c.JSON(http.StatusNoContent, fmt.Sprintf("no data found for %s", dateStr))
	}
	// create struct to hold the problem data
	var problemData struct {
		Answer string `json:"answer"`
	}
	// Retrieve data from the first row
	if err := rows.Scan(&problemData.Answer); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to scan data from database for %s: %v", dateStr, err))
	}
	// Check if there is more than one row (this shouldn't happen)
	if rows.Next() {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("expected one result row, got more than one for %s", dateStr))
	}
	// Load provided answer and correct answer as query parameters
	queryParams := url.Values{}
	queryParams.Add("eq1", answer)
	queryParams.Add("eq2", problemData.Answer)
	// Build full URL with query parameters
	fullURL := fmt.Sprintf("%s?%s", WOLFRAM_URL, queryParams.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error making GET request when verifying solution: %v", err))
	}
	// Parse the body of the response
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error parsing response body from verifier: %v", err))
	}
	// Check if the provided answer and the correct answer were identified as equal
	if strings.TrimSpace(string(body)) != "true" {
		return c.JSON(http.StatusOK, "false")
	}
	// Check if problem should be considered as daily
	// First get the difference in time between the submission time and the assignment time of the problem (midnight at the provided date)
	diffSeconds := submitTime.Sub(inputTime).Seconds()
	// If diff seconds is within a day (plus 90 seconds for leeway), then count as a daily problem
	dailyWindow := DAY_SECONDS + 90
	// Update the users table with the appropriate data (both queries account for whether or not the problem has been solved yet or not)
	if diffSeconds < dailyWindow {
		// Count as daily problem solve (as well as normal problem solve)
		_, err = db.Exec(`
			UPDATE users
			SET 
				daily_problems = daily_problems + 1,
				total_problems = total_problems + 1,
				solves = jsonb_set(solves, $2, '2', true)
			WHERE username = $1 
			AND NOT (solves ? $3); 
		`, username, "{"+dateStr+"}", dateStr)
	} else {
		// Count as normal problem solve
		_, err = db.Exec(`
			UPDATE users
			SET 
				total_problems = total_problems + 1,
				solves = jsonb_set(solves, $2, '1', true)
			WHERE username = $1 
			AND NOT (solves ? $3); 
		`, username, "{"+dateStr+"}", dateStr)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to update user: %v", err))
	}
	return c.JSON(http.StatusOK, "true")
}

// Get user's solve status on a problem
func getProblemStatus(dateStr string, username string) int {
	row := db.QueryRow(`
		SELECT solves->$1
		FROM users
		WHERE username = $2
	`, dateStr, username)
	// create struct to hold the problem data
	var problemStatus int
	// Retrieve data from the first row
	if err := row.Scan(&problemStatus); err != nil {
		// Check if at least one row is returned (if not then the user has no attempted the problem)
		if err == sql.ErrNoRows {
			return 0
		}
		return 0
	}
	return problemStatus
}

func getUser(c echo.Context) error {
	inputUsername := c.Param("username")
	if len(inputUsername) == 0 {
		return c.JSON(http.StatusBadRequest, "no username provided for get user")
	}
	row := db.QueryRow(`
		SELECT username, daily_problems, total_problems, creation_date, solves
		FROM users
		WHERE username = $1
	`, inputUsername)
	// Retrieve data from the first row
	var dbUsername string
	var dailyProblems int
	var totalProblems int
	var creationDate string
	var solvesData []byte
	if err := row.Scan(&dbUsername, &dailyProblems, &totalProblems, &creationDate, &solvesData); err != nil {
		// Check if at least one row is returned
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNoContent, fmt.Sprintf("no data found for user %s", inputUsername))
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to scan data from database for user %s: %v", inputUsername, err))
	}
	// Decode solves field into a map
	var solves map[string]interface{}
	if err := json.Unmarshal(solvesData, &solves); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to parse solves JSON for user %s: %v", inputUsername, err))
	}
	// get the rank of the user (same as position in leaderboard)
	var rank int
	row = db.QueryRow(`
		SELECT COUNT(*) + 1
		FROM users
		WHERE 
			(daily_problems > $1) OR
			(daily_problems = $1 AND total_problems > $2) OR
			(daily_problems = $1 AND total_problems = $2 AND username < $3)
	`, dailyProblems, totalProblems, inputUsername)
	if err := row.Scan(&rank); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to retrieve rank for user %s: %v", inputUsername, err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"username":       dbUsername,
		"daily_problems": dailyProblems,
		"total_problems": totalProblems,
		"creation_date":  creationDate,
		"solves":         solves,
		"rank":           rank,
	})
}

func decodeCursor(s string, out *LeaderboardCursor) error {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

func encodeCursor(c LeaderboardCursor) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func getLeaderboard(c echo.Context) error {
	cursorStr := c.QueryParam("cursor")
	baseQuery := `
		SELECT username, daily_problems, total_problems, creation_date 
		FROM users
	`
	// keep track of args for the query
	args := []interface{}{}
	// note starting rank for the current page
	startingRank := 1
	// if cursor exists, add WHERE clause
	if cursorStr != "" && cursorStr != "start" {
		var cursor LeaderboardCursor
		if err := decodeCursor(cursorStr, &cursor); err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid cursor: %v", err))
		}
		baseQuery += `
			WHERE 
				(daily_problems < $1) OR
				(daily_problems = $1 AND total_problems < $2) OR
				(daily_problems = $1 AND total_problems = $2 AND username > $3)
		`
		args = append(args, cursor.DailyProblems, cursor.TotalProblems, cursor.Username)
		// get the starting rank for the page
		// +2 = +1 for cursor user and +1 for 1-based ranking
		row := db.QueryRow(`
			SELECT COUNT(*) + 2
			FROM users
			WHERE 
				(daily_problems > $1) OR
				(daily_problems = $1 AND total_problems > $2) OR
				(daily_problems = $1 AND total_problems = $2 AND username < $3)
		`, cursor.DailyProblems, cursor.TotalProblems, cursor.Username)
		if err := row.Scan(&startingRank); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to retrieve starting rank: %v", err))
		}
	}
	// ensure deterministic ordering for leaderboard pagination and add limit
	baseQuery += `
		ORDER BY daily_problems DESC, total_problems DESC, username ASC
		LIMIT 11
	`
	// query the database with the constructed query and args
	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to retrieve leaderboard: %v", err))
	}
	defer rows.Close()
	// Store all of the users for the leaderboardData
	var leaderboardData []LeaderboardEntry
	var firstItem *LeaderboardEntry
	// Process of the rows from the query and repopulate the response
	for rows.Next() {
		var username string
		var dailyProblems int
		var totalProblems int
		var creationDate string
		// process row data
		err := rows.Scan(&username, &dailyProblems, &totalProblems, &creationDate)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to parse leaderboard data: %v", err))
		}
		entry := LeaderboardEntry{
			Username:      username,
			DailyProblems: dailyProblems,
			TotalProblems: totalProblems,
			CreationDate:  creationDate,
		}
		if firstItem == nil {
			firstItem = &entry
		}
		leaderboardData = append(leaderboardData, entry)
	}
	if err = rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error occurred when processing leaderboard rows: %v", err))
	}
	// check if there is another page (11 entries found)
	hasNextPage := len(leaderboardData) == 11
	// trim to 10 entries after storing next page flag
	if hasNextPage {
		leaderboardData = leaderboardData[:10]
	}
	// get previous cursor
	var prevCursor *string
	if cursorStr != "" {
		prevQuery := `
			SELECT username, daily_problems, total_problems, creation_date
			FROM users
			WHERE
				(daily_problems > $1) OR
				(daily_problems = $1 AND total_problems > $2) OR
				(daily_problems = $1 AND total_problems = $2 AND username < $3)
			ORDER BY daily_problems ASC, total_problems ASC, username DESC
			LIMIT 11
		`
		rows, err := db.Query(prevQuery, firstItem.DailyProblems, firstItem.TotalProblems, firstItem.Username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to query previous cursor data: %v", err))
		}
		defer rows.Close()
		var prevPage []LeaderboardEntry
		for rows.Next() {
			var username string
			var dailyProblems int
			var totalProblems int
			var creationDate string
			// process row data
			err := rows.Scan(&username, &dailyProblems, &totalProblems, &creationDate)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to parse previous cursor data: %v", err))
			}
			prevPage = append(prevPage, LeaderboardEntry{
				Username:      username,
				DailyProblems: dailyProblems,
				TotalProblems: totalProblems,
				CreationDate:  creationDate,
			})
		}
		if len(prevPage) == 11 {
			first := prevPage[len(prevPage)-1] // prevPage[10]
			encoded, err := encodeCursor(LeaderboardCursor{
				DailyProblems: first.DailyProblems,
				TotalProblems: first.TotalProblems,
				Username:      first.Username,
			})
			if err != nil {
				return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to encode previous cursor: %v", err))
			}
			prevCursor = &encoded
		} else if len(prevPage) > 0 {
			// we have less than a full page, so we are on the first page
			// use a special "start" cursor value to indicate this (necesssary to differentiate from null)
			startStr := "start"
			prevCursor = &startStr
		}
	}
	if hasNextPage {
		last := leaderboardData[len(leaderboardData)-1]
		nextCursor, err := encodeCursor(LeaderboardCursor{
			DailyProblems: last.DailyProblems,
			TotalProblems: last.TotalProblems,
			Username:      last.Username,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to encode next cursor: %v", err))
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"starting_rank": startingRank,
			"data":          leaderboardData,
			"prev_cursor":   prevCursor,
			"next_cursor":   nextCursor,
		})
	} else {
		// leaderboard has no additional pages
		return c.JSON(http.StatusOK, map[string]interface{}{
			"starting_rank": startingRank,
			"data":          leaderboardData,
			"prev_cursor":   prevCursor,
			"next_cursor":   nil,
		})
	}
}

// Function for creating user entry in users table
func createUser(c echo.Context) error {
	// will not require access token (calls via auth0 action)
	username := c.Param("username")
	if len(username) == 0 {
		return c.JSON(http.StatusBadRequest, "no username provided for create user")
	}
	// instead, verify custom API_KEY created only for this endpoint
	key := c.Request().Header.Get("Authorization")
	if key != API_KEY {
		return c.JSON(http.StatusUnauthorized, "invalid api key provided")
	}
	// Get the creation date
	dateStr := time.Now().UTC().Format(DATE_LAYOUT)
	_, err := db.Exec(`
		INSERT INTO users
			(username, daily_problems, total_problems, creation_date, solves)
		VALUES
			($1, 0, 0, $2, $3)
	`, username, dateStr, "{}")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to create user for %s: %v", username, err))
	}
	return c.JSON(200, map[string]string{"message": "username created sucessfully", "username": username})
}

// Function for retrieving the start date
func getStartDate(c echo.Context) error {
	rows := db.QueryRow(`
		SELECT MIN(date_used)
		FROM problems
	`)
	var problemData struct {
		StartDate string `json:"date_used"`
	}
	// Retrieve data from the first row
	if err := rows.Scan(&problemData.StartDate); err != nil {
		// Check if at least one row is returned
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNoContent, "no start date found")
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to scan start date from database: %v", err))
	}
	return c.JSON(http.StatusOK, problemData.StartDate)
}

func main() {
	var err error
	// Attempt to connect to database first
	connStr := os.Getenv("DATABASE_URL")
	WOLFRAM_URL = os.Getenv("WOLFRAM_URL")
	AUTH0_DOMAIN = os.Getenv("AUTH0_DOMAIN")
	AUTH0_AUDIENCE = os.Getenv("AUTH0_AUDIENCE")
	TOKEN_SCOPES = os.Getenv("TOKEN_SCOPES")
	API_KEY = os.Getenv("API_KEY")
	// If string is empty, then script is being run locally (or secrets are not configured with hosting)
	if len(connStr) == 0 {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
			return
		}
		connStr = os.Getenv("DATABASE_URL")
		WOLFRAM_URL = os.Getenv("WOLFRAM_URL")
		AUTH0_DOMAIN = os.Getenv("AUTH0_DOMAIN")
		AUTH0_AUDIENCE = os.Getenv("AUTH0_AUDIENCE")
		TOKEN_SCOPES = os.Getenv("TOKEN_SCOPES")
		API_KEY = os.Getenv("API_KEY")
	}
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// Create Echo instance
	e := echo.New()
	// Configure CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))
	// Define empty route (for testing purposes)
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Welcome to Daily Integral Server!")
	})
	// Define routes
	e.GET("/problemStatus/:date", func(c echo.Context) error {
		dateStr := c.Param("date")
		username := c.QueryParam("username")
		_, msg := validateDate(dateStr)
		if len(msg) != 0 {
			return c.JSON(http.StatusBadRequest, msg)
		}
		if len(username) == 0 {
			return c.JSON(http.StatusBadRequest, "no username provided for problem status")
		}
		status := getProblemStatus(dateStr, username)
		return c.JSON(http.StatusOK, status)
	})
	e.GET("/startdate", getStartDate)
	e.GET("/user/:username", getUser)
	e.GET("/leaderboard", getLeaderboard)
	e.GET("/problem/:date", retrieveProblem)
	e.GET("/solution/:date", retrieveSolutionAndSteps)
	e.POST("/verify/:date", verifySolution, verifierRateLimiter)
	e.POST("/register/:username", createUser)
	// Handle shutdown (gracefully)
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down server: ", err)
		}
	}()
	// Wait for signal to shutdown (via error or manually)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// Close database connection
	db.Close()
	// Stop server (gracefully)
	fmt.Println("Shutting down server...")
	if err := e.Shutdown(context.Background()); err != nil {
		fmt.Printf("Error when shutting down: %v\n", err)
	}
}
