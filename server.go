// server.go
//
// REST APIs with Go and MySql.
//
// Usage:
//
//   # run go server in the background
//   $ go run server.go 

package main

import (
	"fmt"
    	"io/ioutil"
	"strconv"
	"log"
    	"strings"
	"net/http"
    	"encoding/json"
    	"database/sql"
    	_ "github.com/go-sql-driver/mysql"
)


type Panda  struct {
    bloque string
}

//Handle all requests
func Handler(response http.ResponseWriter, request *http.Request){
    response.Header().Set("Content-type", "text/html")
    webpage, err := ioutil.ReadFile("index.html")
    if err != nil {
    http.Error(response, fmt.Sprintf("home.html file error %v", err), 500)
    }
    fmt.Fprint(response, string(webpage));
}

// Respond to URLs of the form /generic/...
func APIHandler(response http.ResponseWriter, request *http.Request){

    //Connect to database
    db, e := sql.Open("mysql", "iacoman:Iaco.2010@tcp(marta.ctrenoefei46.sa-east-1.rds.amazonaws.com:3306)/martadb")
     if( e != nil){
      fmt.Print(e)
     }

    //set mime type to JSON
    response.Header().Set("Content-type", "application/json")
    

	err := request.ParseForm()
	if err != nil {
		http.Error(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

    //can't define dynamic slice in golang
    var result = make([]string,1000)

    switch request.Method {
        case "GET":
            st, err := db.Prepare("select bloque from bloques")
             if err != nil{
              fmt.Print( err );	
             }
             rows, err := st.Query()
             if err != nil {
              fmt.Print( err )
             }
             i := 0
             for rows.Next() {
              var bloque string
              err = rows.Scan( &bloque )
              panda := Panda{"bloque : "+bloque}
              fmt.Printf("%s", panda)
                
               result[i] = fmt.Sprintf("%s", panda)
              i++
             }

            result = result[:i]

        case "POST":
            bloque := request.PostFormValue("bloque")
            st, err := db.Prepare("select bloque from bloques")
             if err != nil{
              fmt.Print( err );
             }
             res, err := st.Exec(bloque)
             if err != nil {
              fmt.Print( err )
             }

             if res!=nil{
                 result[0] = "true"
             }
            result = result[:1]

        case "PUT":
            bloque := request.PostFormValue("bloque")
            id := request.PostFormValue("id")

            st, err := db.Prepare("select bloque from bloques")
             if err != nil{
              fmt.Print( err );
             }
             res, err := st.Exec(bloque,id)
             if err != nil {
              fmt.Print( err )
             }

             if res!=nil{
                 result[0] = "true"
             }
            result = result[:1]
        case "DELETE":
            id := strings.Replace(request.URL.Path,"/getBlocks","",-1)
            st, err := db.Prepare("select bloque from bloques")
             if err != nil{
              fmt.Print( err );
             }
             res, err := st.Exec(id)
             if err != nil {
              fmt.Print( err )
             }

             if res!=nil{
                 result[0] = "true"
             }
            result = result[:1]

        default:
    }
    
    json, err := json.Marshal(result)
    if err != nil {
        fmt.Println(err)
        return
    }

    //fmt.Sprintf("%s", string(json))
    //fmt.Sprintf("%s", result)


	// Send the text diagnostics to the client.
    //fmt.Fprintf(response,"%v",result)
    fmt.Fprintf(response,"%v",string(json))

	//fmt.Fprintf(response, " request.URL.Path   '%v'\n", request.Method)
    db.Close()
}


func main(){
	port := 80
    var err string
	portstring := strconv.Itoa(port)

	mux := http.NewServeMux()
	mux.Handle("/getBlocks", http.HandlerFunc( APIHandler ))
	mux.Handle("/", http.HandlerFunc( Handler ))

	// Start listing on a given port with these routes on this server.
	log.Print("Listening on port " + portstring + " ... ")
	errs := http.ListenAndServe(":" + portstring, mux)
	if errs != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
