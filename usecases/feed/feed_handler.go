package feed

import (
	"net/http"
	"strconv"

	"github.com/ranjabi/twitter-clone-backend/app"
	"github.com/ranjabi/twitter-clone-backend/errmsg"
	"github.com/ranjabi/twitter-clone-backend/response"
	"github.com/ranjabi/twitter-clone-backend/utils"
)

type Handler struct {
	feedService Service
}

func NewHandler(feedService Service) Handler {
	return Handler{
		feedService,
	}
}

func (h Handler) HandleGetFeed(w http.ResponseWriter, r *http.Request) *app.Error {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	query := r.URL.Query()
	pageStr := query.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	feed, err := h.feedService.GetFeed(id, page)
	if err != nil {
		return utils.HandleErr(err)
	}

	if err := response.Json(w, http.StatusOK, response.Payload{Data: feed}); err != nil {
		return err
	}

	return nil
}
