package api_test

import (
	"bytes"
	"gitlab.eng.vmware.com/nidhig1/goassignment-postandcomments/pkg/api"
	"gitlab.eng.vmware.com/nidhig1/goassignment-postandcomments/pkg/cache"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}

var _ = Describe("Posts API Testing", func() {

	c := cache.NewPostsAndCommentsCache()
	apiReceiver := api.CacheApi{
		Cache: c,
	}

	Describe("Fetching One Posts", func() {
		//Test Case that POST ID=0
		Context("With ID:0", func() {
			It("Should return NOT FOUND", func() {
				req, _ := http.NewRequest("GET", "/getPosts", nil)
				q := req.URL.Query()
				q.Add("id", "0")
				req.URL.RawQuery = q.Encode()
				response := httptest.NewRecorder()
				handler := http.HandlerFunc(apiReceiver.GetOnePosts)
				handler.ServeHTTP(response, req)
				expected := `NO SUCH POST FOUND IN DATABASE`
				Expect(response.Body.String()).To(Equal(expected))
			})
		})
		//Test Case that POST ID=-89
		Context("With ID:-89", func() {
			It("Should return NOT FOUND", func() {
				req, _ := http.NewRequest("GET", "/getPosts", nil)
				q := req.URL.Query()
				q.Add("id", "-89")
				req.URL.RawQuery = q.Encode()
				response := httptest.NewRecorder()
				handler := http.HandlerFunc(apiReceiver.GetOnePosts)
				handler.ServeHTTP(response, req)
				expected := `NO SUCH POST FOUND IN DATABASE`
				Expect(response.Body.String()).To(Equal(expected))
			})
		})
		//Test Case that POST ID=67
		Context("With ID:67", func() {
			It("Should return FOUND", func() {
				req, _ := http.NewRequest("GET", "/getPosts", nil)
				q := req.URL.Query()
				q.Add("id", "67")
				req.URL.RawQuery = q.Encode()
				response := httptest.NewRecorder()
				handler := http.HandlerFunc(apiReceiver.GetOnePosts)
				handler.ServeHTTP(response, req)
				Expect(response.Body.String()).To(Equal(api.Post67))
			})
		})
		//Test Case that fetches all Posts
		Context("Get All Posts", func() {
			It("Should return FOUND All", func() {
				req, _ := http.NewRequest("GET", "/posts", nil)
				response := httptest.NewRecorder()
				handler := http.HandlerFunc(apiReceiver.GetAllPosts)
				handler.ServeHTTP(response, req)
				Expect(response.Body.String()).To(Equal(api.AllPostsData))
			})
		})
		////Test Case that adds post
		Context("Add One Post", func() {
			It("Should return FOUND All", func() {
				newPosts2 := []byte(`{
					"userId": 789,
					"title": "nidhi has added a post",
					"body": "nidhi is new to banglore .she plans to expllore the city for next few days",
					"comments": [
						{
						"postId": 789,
						"id":1,
						"name":"nidhi",
						"body":"nidhi is new to banglore .she plans to expllore the city for next few days",
						"email":"nidhi@vmware.com"
						}
					]
		         }`)
				req, _ := http.NewRequest("POST", "/addPostsAndComments", bytes.NewBuffer(newPosts2))
				req.Header.Set("Content-Type", "application/json")
				response := httptest.NewRecorder()
				handler := http.HandlerFunc(apiReceiver.AddOnePostAndComment)
				handler.ServeHTTP(response, req)
				Expect(response.Body.String()).To(Equal("Post added to database successfully"))
			})
		})
	})

})
