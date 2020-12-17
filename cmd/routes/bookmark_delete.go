package routes

import "net/http"

//DeleteBookmark removes a bookmark from the users list if it exists.
func (s *Server) DeleteBookmark(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
