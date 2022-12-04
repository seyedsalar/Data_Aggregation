package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/grokify/mogo/fmt/fmtutil"

	stackexchange "github.com/grokify/go-stackoverflow/client"

	"github.com/grokify/go-stackoverflow/util"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Ether struct {
	Transaction_Hash string
	Transfer_Amount  int
	Action           string
	Time             time.Time
}

type Total struct {
	Total_amount int
}

func Running(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found. hi", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "views/index.html")

	case "POST":
		fmt.Fprintf(w, "Sorry, only GET methods are supported.")

	default:
		fmt.Fprintf(w, "Oops!")

	}

}

func Ethereum_Transactions(w http.ResponseWriter, r *http.Request) {

	array := [14]interface{}{}
	var obj [13]Ether
	var TA [1]Total
	switch r.Method {

	case "POST":

		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		var iter = 1
		var sum = 0

		address_ether := string(r.FormValue("Ether_address"))

		//client, err := ethclient.Dial("http://127.0.0.1:7545")

		client, err := ethclient.Dial(address_ether)

		if err != nil {
			//log.Fatal(err)
			return

		}

		for i := 15; i > 2; i-- {
			j := 15 - i
			blockNumber := big.NewInt(int64(i))
			block, err := client.BlockByNumber(context.Background(), blockNumber)
			if err != nil {
				log.Fatal(err)
			}

			var date = int64(block.Time())
			unixTimeUTC := time.Unix(date, 0)
			obj[j].Time = unixTimeUTC
			fmt.Println(unixTimeUTC)
			//fmt.Fprintf(w, "\n Transaction Time : %v " , unixTimeUTC)

			if iter <= 4 || iter == 11 {
				for _, tx := range block.Transactions() {
					obj[j].Transaction_Hash = tx.Hash().Hex()
					obj[j].Action = "Withdrawal"
					fmt.Println(tx.Hash().Hex()) //
					//fmt.Fprintf(w, "\n Transaction Hash : %s " , tx.Hash().Hex())
					var data = tx.Data() // []
					obj[j].Transfer_Amount = int(data[35])
					fmt.Println("withdrawal :", data[35])
					sum -= int(data[35])
					//fmt.Fprintf(w, "\n withdrawal : %d " , data[35])
					array[j] = obj[j]

				}
			} else if (iter > 4 && iter <= 10) || iter == 12 || iter == 13 {
				for _, tx := range block.Transactions() {
					obj[j].Transaction_Hash = tx.Hash().Hex()
					obj[j].Action = "Deposit"
					fmt.Println(tx.Hash().Hex()) //
					//fmt.Fprintf(w, "\n Transaction Hash : %s " , tx.Hash().Hex())
					var data = tx.Data() // []
					obj[j].Transfer_Amount = int(data[35])
					fmt.Println("Deposit: ", data[35])
					sum += int(data[35])
					//fmt.Fprintf(w, "\n Deposit : %d " , data[35])
					array[j] = obj[j]
				}
			}

			iter++
		}
		TA[0].Total_amount = sum
		fmt.Println(" \nTotal Customer Balance = ", sum)
		//fmt.Fprintf(w," \nTotal Customer Balance = %d", sum)
		array[13] = TA[0]
		temp1, _ := json.Marshal(array)
		temp1, _ = json.MarshalIndent(array, "", "  ")
		w.Write(temp1)
		fmt.Println(string(temp1))

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func Stackoverflow(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		apiClient := stackexchange.NewAPIClient(
			stackexchange.NewConfiguration())

		site := "stackoverflow"
		//userId := "1908967"
		userId := string(r.FormValue("Stack_id"))
		history, err := util.GetReputationHistoryAll(apiClient, site, userId)
		if err != nil {
			log.Fatal(err)
		}

		fmtutil.PrintJSON(history)
		fmt.Printf("COUNT [%v]\n", len(history.Items))
		temp2, _ := json.Marshal(history)
		temp2, _ = json.MarshalIndent(history, "", "  ")
		w.Write(temp2)
		fmt.Println("DONE")

	}
}

func Hyperledger_Transactions(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		address_hyper := string(r.FormValue("Hyper_address"))
		bearer_temp := string(r.FormValue("Bearer"))
		bearer := "Bearer" + " " + bearer_temp

		client := &http.Client{}
		req, err := http.NewRequest("GET", address_hyper, nil)
		req.Header.Add("Authorization", bearer)

		if err != nil {
			log.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(body)
	}

}

func main() {

	http.HandleFunc("/", Running)
	http.HandleFunc("/api/Ethereum", Ethereum_Transactions)
	http.HandleFunc("/api/Hyperledger", Hyperledger_Transactions)
	http.HandleFunc("/api/Stackoverflow", Stackoverflow)
	fmt.Printf("Starting server for HTTP POST...\n")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
