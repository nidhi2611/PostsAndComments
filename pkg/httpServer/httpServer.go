package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.eng.vmware.com/nidhig1/goassignment-postandcomments/pkg/api"
	"gitlab.eng.vmware.com/nidhig1/goassignment-postandcomments/pkg/cache"
)

type Server struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewServer(address string, readTimeout, writeTimeout time.Duration) *Server {
	return &Server{Addr: address, ReadTimeout: readTimeout, WriteTimeout: writeTimeout}
}
func (s *Server) Start(ctx context.Context) error {
	//initializing cache
	c := cache.NewPostsAndCommentsCache()
	srv := http.Server{
		Addr:         s.Addr,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		Handler:      s.initHandlers(c),
	}

	go func() {
		ticker(ctx, c)
		<-ctx.Done()
		fmt.Println("attempting graceful shutdown of server")
		srv.SetKeepAlivesEnabled(false)
		closeCtx, closeFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer closeFn()
		_ = srv.Shutdown(closeCtx)
	}()

	return srv.ListenAndServe()
}

func (s *Server) initHandlers(c *cache.PostsAndCommentsCache) http.Handler {
	apiReciever := api.CacheApi{
		Cache: c,
	}
	r := mux.NewRouter()
	r.HandleFunc("/posts", apiReciever.GetAllPostsAndComments).Methods("GET")
	r.HandleFunc("/getPosts", apiReciever.GetOnePosts).Methods("GET")
	r.HandleFunc("/addPostsAndComments", apiReciever.AddOnePostAndComment).Methods("POST")
	r.HandleFunc("/addPost", api.AddOnePost).Methods("POST")
	return r
}
func ticker(ctx context.Context, c *cache.PostsAndCommentsCache) error {
	fmt.Println("Starting Ticker")
	apiReciever := api.CacheApi{
		Cache: c,
	}
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := apiReciever.ReadPosts()
			if err != nil {
				fmt.Println("read posts error", err)
			}
			er := apiReciever.ReadComments()
			if er != nil {
				fmt.Println("read comments error", er)
			}
		case <-ctx.Done():
			fmt.Printf("Closing Ticker Function")
			return ctx.Err()
		}
	}
}
