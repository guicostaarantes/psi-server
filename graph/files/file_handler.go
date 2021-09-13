package files

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/guicostaarantes/psi-server/graph/resolvers"
)

type FileHandler struct {
	Resolvers *resolvers.Resolver
}

func (f FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userID").(string)

	if userID == "" {
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("403 forbidden"))
		return
	}

	if r.Method == "GET" {
		fileName := chi.URLParam(r, "name")

		data, readErr := f.Resolvers.ReadFileService().Execute(fileName)
		if readErr != nil {
			w.WriteHeader(500)
			w.Write([]byte("500 internal server error"))
			return
		}

		if len(data) == 0 {
			w.WriteHeader(404)
			w.Write([]byte("404 page not found"))
			return
		}

		w.Write(data)
		return
	}

	w.WriteHeader(404)
	w.Write([]byte("404 page not found"))
}
