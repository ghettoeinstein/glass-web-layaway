package router

func SetUserRoutes(router *mux.Router) *mux.Router {
	//Disable external account createion, re-enable once roles & permissions are set correctly.
	router.HandleFunc("/users/register", controllers.Register).Methods("POST")

	router.HandleFunc("/login", controllers.Login).Methods("POST")

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "Auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		w.Header()["Location"] = []string{"/login"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	}).Methods("GET")

	return router
}
