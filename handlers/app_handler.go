package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"
)

type AppHandler func(http.ResponseWriter, *http.Request) *models.AppError

// TODO: confirm this by looking at error trace when err is nul at ServiceError <--- The ServeHTTP method called by the appHandler function and displays the returned error
func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		if e.Code == 0 {
			e.Code = http.StatusInternalServerError
		}

		fmt.Println(utils.ColorLog(strconv.Itoa(e.Code), utils.RED), utils.ColorLog(http.StatusText(e.Code), utils.RED))
		fmt.Println(utils.ColorLog(e.Error(), utils.RED))

		res, err := json.Marshal(models.ErrorResponse{Message: e.Message})
		if err != nil {
			http.Error(w, utils.ErrMsgFailedToSerializeResponseBody, http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.Code)
		w.Write(res)
	}
}
