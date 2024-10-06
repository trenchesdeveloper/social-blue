package main

import "net/http"

// healthCheckHandler godoc
//	@Summary		Check the health status of the service
//	@Description	Check the health status of the service
//	@ID				health-check
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (s *server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok",
		"environment": s.config.Environment,
		"version":     version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		s.internalServerError(w, r, err)
	}

}
