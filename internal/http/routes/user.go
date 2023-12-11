package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sublinks/federation/internal/activitypub"
	"sublinks/federation/internal/lemmy"
	"sublinks/federation/internal/logging/logger"

	"github.com/gorilla/mux"
)

func SetupUserRoutes(r *mux.Router) {
	r.HandleFunc("/u/{user}", getUserInfoHandler).Methods("GET")
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.Background()
	c := lemmy.GetLemmyClient(ctx)
	logger.GetLogger().Println(fmt.Sprintf("Looking up user %s", vars["user"]))
	user, err := c.GetUser(ctx, vars["user"])
	if err != nil {
		logger.GetLogger().Println("Error reading user", err)
		return
	}

	userLd := activitypub.ConvertUserToApub(user, r.Host)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(userLd, "", "  ")
	w.Write(content)
}
