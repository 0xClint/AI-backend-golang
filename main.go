package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Prompt struct {
	Query string `json:"query"`
}

type Response struct {
	Answer            string   `json:"answer"`
	SourceDocs        []string `json:"source_docs"`
	GeneratedQuestion string   `json:"generated_question"`
}

func main() {
	fmt.Println("Welcome to Clint server")
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	//Listen to port 4000
	r.HandleFunc("/", postRequest).Methods("POST")
	log.Fatal(http.ListenAndServe(":4000", r))

}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to Clint's Server! :)</h2>"))
	json.NewEncoder(w).Encode("Welcome to Clint's Server! :)")

}

func postRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	//if response data is empty or not Valid
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send valid data")
		return
	}

	var data Prompt

	_ = json.NewDecoder(r.Body).Decode(&data)

	fmt.Println("query : ", data.Query)
	response := AIPostRequest(data.Query)

	json.NewEncoder(w).Encode(response)
}

// Helper Func
func AIPostRequest(query string) *Response {
	requestData := map[string]string{
		// "question":            "What is Polygon?",
		"question":            query,
		"knowledge_source_id": "0x7521b754a946844c720a4772f16b0574680223a8",
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return nil
	}

	req, err := http.NewRequest("POST", "https://rag-chat-ml-backend-prod.flock.io/chat/conversational_rag_chat", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}

	// Add API key to request header
	req.Header.Set("x-api-key", "54c488b4a94b415b989c19e0f29d199c")

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	var response Response
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil
	}

	// Return the Response struct
	return &response
}
