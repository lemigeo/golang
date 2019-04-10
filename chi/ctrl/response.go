package ctrl

import ( 
    "encoding/json" 
	"fmt" 
	"strconv"
    "net/http" 
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
    respondwithJSON(w, map[string]string{"code": strconv.Itoa(code), "message": msg})
}

func respondwithJSON(w http.ResponseWriter, payload interface{}) {
    response, _ := json.Marshal(payload)
    fmt.Println(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(response)
}