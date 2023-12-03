package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"spacecraft/domain"
)

// Opinion:
// I suspect you got influenced by the  DDD approach of Damiano Petro Ungaro.
// While I respect the choice - and in principle there is nothing wrong with it -
// I personally suggest to keep the "infra" packages inside the `internal/` folder,
// just defined inside their respective packages, eg:
//   spacecraft/
//     cmd/
//       server/
//         main.go
//     internal/
//	     handlers/
//         ...files.go...
//	     store/
//         ...files.go...
// while the cmd/server can contain just the
// main.go which glues all the subpackages from `internal/`

// Spacecraft ...
// (suspected) major: having Spacecraft like this is race prone.
// What about wrapping spacecraft with an accessor or a struct?
// `Accessor pattern` procvdes access to the global []*domain.Spacecraft via provide Get/Set pkg functions
// by using a mutex (similar to the history object, see https://github.com/MyPublicProjects/term-calc/blob/main/internal/history/history.go#L11-L16).
// Using a struct can work too:
//
//		type Spacecrafts stuct{
//		 data  []*domain.Spacecraft
//		 m     sync.Mutex
//	}
//
// I understand that this `major` might not be an issue, given the single threaded nature of the program so far
var Spacecraft []*domain.Spacecraft

func extractIntQueryParam(r *http.Request, paramName string, defaultValue int) (*int, error) {
	urlValues := r.URL.Query()
	var result int
	rawParam := urlValues.Get(paramName)
	if rawParam == "" {
		return &defaultValue, nil
	}
	result, err := strconv.Atoi(rawParam)
	if err != nil {
		return nil, fmt.Errorf("err while fetching value from query string: %v", err)
	}
	return &result, nil
}

// Getspacecraft major: naming convention: should be GetSpacecraft
func Getspacecraft(w http.ResponseWriter, r *http.Request) {
	pageSize, err := extractIntQueryParam(r, "pageSize", 100)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	pageNumber, err := extractIntQueryParam(r, "pageNumber", 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// suggestion:
	// code below is hard to read: extract in its own function
	// increases clarity and allows domain logic testing.

	// Also I don't feel it belongs to the handler portion of the code.
	// it probably should be a "Getter" on the Spacecraft type,
	// see comments on Spacecraft above.
	low := *pageNumber * *pageSize
	high := low + *pageSize
	var result domain.SpacecraftWrapper
	result.PageNumber = *pageNumber
	result.PageSize = *pageSize
	result.NumberOfElements = *pageSize
	result.TotalPages = (len(Spacecraft) / *pageSize) + 1
	result.TotalElements = len(Spacecraft)
	if high >= len(Spacecraft) {
		result.NumberOfElements = len(Spacecraft[low:])
		result.Data = Spacecraft[low:]
	} else {
		result.Data = Spacecraft[low:high]
	}

	data, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	time.Sleep(2 * time.Second) // why? debug leftover ?
	w.Write([]byte(data))       // no need to cast []byte()
}
