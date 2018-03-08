package routes

import (
	"net/http"
	"strconv"

	"github.com/maksadbek/dumbwall/internal/actions"
	"github.com/maksadbek/dumbwall/internal/posts"
	"go.uber.org/zap"
)

func (r *Routes) NewPost(w http.ResponseWriter, req *http.Request) {
	context := struct {
		flash flash
	}{
		flash: flash{
			Notice: "you're creating a post",
			Alert:  "be careful!",
			Custom: map[string]string{
				"first_ever_post": "hey, post anything what you want",
			},
		},
	}

	r.templates.ExecuteTemplate(w, "new_post", context)
}

func (r *Routes) CreatePost(w http.ResponseWriter, req *http.Request) {
	userID, err := r.validateToken(w, req)
	if err != nil {
		r.logger.Error("failed to validate token", zap.Error(err))
		return
	}

	req.ParseForm()

	title := req.PostForm.Get("title")
	body := req.PostForm.Get("body")

	post, err := r.db.CreatePost(userID, posts.Post{
		Title: title,
		Body:  body,
	})

	if err != nil {
		r.logger.Error("failed to create a post", zap.Error(err))
		return
	}

	r.logger.Debug("created post", zap.Any("post", post))
}

func (r *Routes) UpdatePost(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) DeletePost(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) UpPost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	userID, err := r.validateToken(w, req)
	if err != nil {
		http.Redirect(w, req, req.Referer(), http.StatusFound)
		return
	}

	postID, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil {
		http.Redirect(w, req, req.Referer(), http.StatusFound)
		return
	}

	err = r.db.VotePost(userID, postID, actions.ActionUp)
	if err != nil {
		http.Redirect(w, req, req.Referer(), http.StatusFound)
		return
	}

	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func (r *Routes) DownPost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	userID, err := r.validateToken(w, req)
	if err != nil {
		http.Redirect(w, req, req.Referer(), http.StatusFound)
		return
	}

	postID, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil {
		http.Redirect(w, req, req.Referer(), http.StatusFound)
		return
	}

	err = r.db.VotePost(userID, postID, actions.ActionDown)
	if err != nil {
		http.Redirect(w, req, req.Referer(), http.StatusFound)
		return
	}

	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func (r *Routes) Post(w http.ResponseWriter, req *http.Request) {

}
