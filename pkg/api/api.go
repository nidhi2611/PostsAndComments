package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil" //nolint
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"gitlab.eng.vmware.com/nidhig1/goassignment-postandcomments/pkg/cache"
	types "gitlab.eng.vmware.com/nidhig1/goassignment-postandcomments/pkg/types"
)

var postsData []types.Post
var commentsData []types.Comment
var postMap = make(map[int]*types.Post)

const postUrl = "https://jsonplaceholder.typicode.com/posts"
const commentsUrl = "https://jsonplaceholder.typicode.com/comments"

type CacheApi struct {
	Cache *cache.PostsAndCommentsCache
}

func (c *CacheApi) GetAllPostsAndComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	out := json.NewEncoder(w)
	out.SetIndent("", "  ")
	err := out.Encode(c.Cache.GetAll())
	if err != nil {
		fmt.Println("Error while writing data")
	}
}
func GetOnePostsAndComment(w http.ResponseWriter, r *http.Request) {
	var found bool
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, val := range postsData {
		//handle error
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			fmt.Println("cannot fetch id from url")
		}
		if postsData[index].Id == id {
			out := json.NewEncoder(w)
			out.SetIndent("", "  ")
			err := out.Encode(val)
			found = true
			if err != nil {
				fmt.Println("Error while writing data")
			}
		}
	}
	if found == false {
		out := json.NewEncoder(w)
		out.SetIndent("", "  ")
		err := out.Encode("NO POST PRESENT WITH THIS ID")
		if err != nil {
			fmt.Println("Error while writing data")
		}
	}
}
func (c *CacheApi) GetOnePosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s := r.URL.Query().Get("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("cannot fetch id from url")
	}
	result := c.Cache.GetOnePost(id)
	if result == nil {
		c.ReadPosts()
		c.ReadComments()
		res := c.Cache.GetOnePost(id)
		if res == nil {
			io.WriteString(w, "NO SUCH POST FOUND IN DATABASE")
		} else {
			out := json.NewEncoder(w)
			out.SetIndent("", "  ")
			er := out.Encode(res)
			if er != nil {
				fmt.Println("error while writing data")
			}
		}
	} else {
		out := json.NewEncoder(w)
		out.SetIndent("", "  ")
		er := out.Encode(result)
		if er != nil {
			fmt.Println("error while writing data")
		}
	}
}

// pass the below in json in body
//
//	{
//		"userId": 789,
//		"title": "nidhi has added a post",
//		"body": "nidhi is new to banglore .she plans to expllore the city for next few days",
//		"comments": [
//			{
//			"postId": 789,
//			"id":1,
//			"name":"nidhi",
//			"body":"nidhi is new to banglore .she plans to expllore the city for next few days",
//			"email":"nidhi@vmware.com"
//			}
//		]
//	}
func (c *CacheApi) AddOnePostAndComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post types.Post
	if r.Body == nil {
		fmt.Println("Please enter something inside the body")
	} else {
		k := json.NewDecoder(r.Body).Decode(&post)
		if k != nil {
			fmt.Println("cannot read data")
		}
		rand.Seed(time.Now().UnixNano())
		post.Id = rand.Intn(500)
		postsData = append(postsData, post)
		//
		io.WriteString(w, "Post added to database successfully")
		//if err != nil {
		//	fmt.Println("Error while writing data")
		//}
		c.Cache.AddAll(postsData)
	}

}

// pass the below in json in body
//
//	{
//		"userId": 789,
//		"title": "nidhi has added a post",
//		"body": "nidhi is new to banglore .she plans to expllore the city for next few days",
//	}
func AddOnePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post types.Post
	if r.Body == nil {
		fmt.Println("Please enter something inside the body")
	} else {
		k := json.NewDecoder(r.Body).Decode(&post)
		if k != nil {
			fmt.Println("cannot read data")
		}
		rand.Seed(time.Now().UnixNano())
		post.Id = rand.Intn(500)
		postsData = append(postsData, post)
		out := json.NewEncoder(w)
		out.SetIndent("", "  ")
		err := out.Encode(postsData)
		if err != nil {
			fmt.Println("Error while writing data")
		}
	}
}
func (c *CacheApi) ReadPosts() error {
	fmt.Println("inside read posts")
	resPost, err := http.Get(postUrl)
	if err != nil {
		return err
	}
	dataBytes, er := ioutil.ReadAll(resPost.Body)
	if er != nil {
		return er
	}
	e := json.Unmarshal(dataBytes, &postsData)
	if e != nil {
		return e
	}
	for i := range postsData {
		postMap[postsData[i].Id] = &postsData[i]
	}
	return nil
}

func (c *CacheApi) ReadComments() error {
	var found bool
	fmt.Println("inside read comments")
	resComment, err := http.Get(commentsUrl)
	if err != nil {
		return err
	}
	dataBytes, er := ioutil.ReadAll(resComment.Body)
	if er != nil {
		return er
	}
	e := json.Unmarshal(dataBytes, &commentsData)
	if e != nil {
		return e
	}
	for i := range commentsData {
		postRef := postMap[commentsData[i].PostId]
		for j := range postRef.Comments {
			if postRef.Comments[j].Id == commentsData[i].Id {
				found = true
				break
			}
		}
		if found == false {
			postRef.Comments = append(postRef.Comments, &commentsData[i])
		}
	}
	for i := range commentsData {
		postRef := postMap[commentsData[i].PostId]
		for j := range postRef.Comments {
			if postRef.Comments[j].Id == commentsData[i].Id {
				found = true
				break
			}
		}
		if found == false {
			postRef.Comments = append(postRef.Comments, &commentsData[i])
		}
	}
	c.Cache.AddAll(postsData)
	return nil
}
func (c *CacheApi) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	out := json.NewEncoder(w)
	out.SetIndent("", "  ")
	err := out.Encode(postsData)
	if err != nil {
		fmt.Println("Error while writing data")
	}
}
