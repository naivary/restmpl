package ctrl

import "github.com/go-chi/chi/v5"

func GetEndpoints(views *Views) chi.Router {
	main := chi.NewRouter()

	sys := sysEndpoints(views)
	main.Mount("/sys", sys)

	return main
}

func sysEndpoints(views *Views) chi.Router {
	r := chi.NewRouter()
	r.Get("/health", views.Sys.Health)

	return r
}
