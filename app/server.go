package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"kubrabkaf/business/chain"
	"kubrabkaf/business/kafkaClient"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type Config struct {
	ListenAddr string
}

type APIServer struct {
	*Config
	ETHClient *chain.ETHClient
	EthKafka  *kafkaClient.EthKafka
}

type ApiError struct {
	Error string `json:"error"`
}

func NewServer(cfg *Config) (*APIServer, error) {
	// Connect to the Ethereum network using Geth or Ganache client.
	ethc, err := ethclient.Dial("http://localhost:7545")
	if err != nil {
		return nil, fmt.Errorf("oops! There was a problem: %w", err)
	}

	// Create a new ETHClient wrapper instance.
	cEth := &chain.ETHClient{
		Client: ethc,
	}

	kf, err := kafkaClient.NewEthKafka()
	if err != nil {
		log.Fatal(err)
	}

	return &APIServer{
		Config:    cfg,
		ETHClient: cEth,
		EthKafka:  kf,
	}, nil
}

func (s *APIServer) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/latestblock", makeHTTPHandleFunc(s.LatestBlock)).Methods(http.MethodGet)
	router.HandleFunc("/sendeth", makeHTTPHandleFunc(s.SendETH)).Methods(http.MethodPost)

	log.Println("JSON API server running on port: ", s.Config.ListenAddr)

	err := http.ListenAndServe(s.Config.ListenAddr, router)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) LatestBlock(w http.ResponseWriter, r *http.Request) error {
	b, err := s.ETHClient.GetLatestBlock()
	if err != nil {
		return fmt.Errorf("error getting latest block: %v", err)
	}

	return WriteJSON(w, http.StatusOK, b)
}

func (s *APIServer) SendETH(w http.ResponseWriter, r *http.Request) error {
	var transferData chain.TransferEthRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&transferData)
	if err != nil {
		return fmt.Errorf("error decoding request body: %v", err)
	}

	marshalled, err := json.Marshal(transferData)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	// Produce sends a message to the Kafka broker.
	err = s.EthKafka.Produce(kafkaClient.EthTopic, marshalled)
	if err != nil {
		fmt.Printf("failed to produce message: %v\n", err)
	}

	h, err := s.ETHClient.TransferEth(transferData.PrivateKey, transferData.To, transferData.Amount)
	if err != nil {
		return fmt.Errorf("error sending eth: %v", err)
	}

	return WriteJSON(w, http.StatusOK, h)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
