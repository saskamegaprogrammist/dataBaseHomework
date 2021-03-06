package models

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	Id int `json:"id"`
	Message string `json:"message"`
	Date time.Time `json:"created"`
	Parent int `json:"parent"`
	Edited bool `json:"isEdited,omitempty"`
	User string `json:"author"`
	Forum string `json:"forum"`
	Thread int`json:"thread"`
}

type PostRelated struct {
	Post *Post `json:"post"`
	User *User `json:"author,omitempty"`
	Forum *Forum `json:"forum,omitempty"`
	Thread *Thread`json:"thread,omitempty"`
}


func GetPostsByThread(limit int, sinceStr string, desc bool, sort string, thread Thread) ([]Post, error) {
	postsFound := make([]Post, 0)
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		//log.Println(err)
		errRollback := transaction.Rollback()
		if err != nil {
			log.Fatalln(errRollback)
		}
		return postsFound, err
	}
	var forumSlug string
	if thread.Id != 0 {
		rows := transaction.QueryRow("SELECT forumslug FROM thread WHERE id = $1", thread.Id)
		err = rows.Scan(&forumSlug)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return postsFound, fmt.Errorf("can't find thread with id %d", thread.Id)
		}
	} else {
		rows := transaction.QueryRow("SELECT id, forumslug FROM thread WHERE slug = $1", thread.Slug)
		err = rows.Scan(&thread.Id, &forumSlug)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return postsFound, fmt.Errorf("can't find thread with slug %s", thread.Slug)
		}
	}

	var rows *pgx.Rows
	if sort != "" {
		switch sort {
		case "flat":
			if sinceStr != ""{
				since, _ := strconv.Atoi(sinceStr)
				if limit != -1 {
					if desc {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 and id  < $2 ORDER BY id DESC LIMIT $3`, thread.Id, since, limit)
					} else {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 and id  > $2 ORDER BY id LIMIT $3 `, thread.Id, since, limit)
					}
				} else {
					if desc {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 and id  < $2 ORDER BY id DESC `, thread.Id, since)
					} else {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 and id  > $2 ORDER BY id `, thread.Id, since)
					}
				}
			} else {
				if limit != -1 {
					if desc {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 ORDER BY id DESC LIMIT $2`, thread.Id, limit)
					} else {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post  WHERE threadid = $1 ORDER BY id LIMIT $2`, thread.Id, limit)
					}
				} else {
					if desc {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 ORDER BY id DESC`, thread.Id)
					} else {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 ORDER BY id `, thread.Id)
					}
				}
			}

		case "tree":
			if sinceStr != ""{
				since, _ := strconv.Atoi(sinceStr)
				if limit != -1 {
					if desc {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post 
							WHERE threadid = $1 AND (path < (SELECT path FROM post WHERE id = $2)) ORDER BY path DESC LIMIT $3 `, thread.Id, since, limit)
					} else {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post 
							WHERE threadid = $1 AND (path > (SELECT path FROM post WHERE id = $2)) ORDER BY path LIMIT $3 `, thread.Id, since, limit)
					}
				} else {
					if desc {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post 
							WHERE threadid = $1 AND (path < (SELECT path FROM post WHERE id = $2)) ORDER BY path DESC  `, thread.Id, since)
					} else {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post 
							WHERE threadid = $1 AND (path > (SELECT path FROM post WHERE id = $2)) ORDER BY path  `, thread.Id, since)
					}
				}
			} else {
				if limit != -1 {
					if desc {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post 
							WHERE threadid = $1 ORDER BY path DESC LIMIT $2 `, thread.Id, limit)
					} else {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post 
							WHERE threadid = $1 ORDER BY path LIMIT $2 `, thread.Id, limit)
					}
				} else {
					if desc {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post 
							WHERE threadid = $1 ORDER BY path DESC  `, thread.Id)
					} else {
						rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post 
							"WHERE threadid = $1 ORDER BY path  `, thread.Id)
					}
				}
			}
		case "parent_tree":
			if sinceStr != ""{
				since, _ := strconv.Atoi(sinceStr)
				if limit != -1 {
					if desc {
						rows, err = transaction.Query(`SELECT  id, message, created, parent, isEdited, usernick  FROM post p WHERE p.threadid = $1 and p.path[1] IN ( 
							SELECT p2.path[1] FROM post p2 WHERE p2.threadid = $1 AND p2.parent = 0 and p2.path[1] < 
							(SELECT p3.path[1] from post p3 where p3.id = $2) ORDER BY p2.path DESC LIMIT $3) ORDER BY p.path[1] DESC, p.path[2:]`, thread.Id, since, limit)
					} else {
						rows, err = transaction.Query(`SELECT  id, message, created, parent, isEdited, usernick  FROM post p WHERE p.threadid = $1 and p.path[1] IN ( 
							SELECT p2.path[1] FROM post p2 WHERE p2.threadid = $1 AND p2.parent = 0 and p2.path[1] > 
							(SELECT p3.path[1] from post p3 where p3.id = $2) ORDER BY p2.path LIMIT $3) ORDER BY p.path `, thread.Id, since, limit)
					}
				} else {
					if desc {
						rows, err = transaction.Query(`SELECT  id, message, created, parent, isEdited, usernick  FROM post p WHERE p.threadid = $1 and p.path[1] IN ( 
							SELECT p2.path[1] FROM post p2 WHERE p2.threadid = $1 AND p2.parent = 0 and p2.path[1] < 
							(SELECT p3.path[1] from post p3 where p3.id = $2) ORDER BY p2.path DESC) ORDER BY p.path[1] DESC, p.path[2:]`, thread.Id, since, limit)
					} else {
						rows, err = transaction.Query(`SELECT  id, message, created, parent, isEdited, usernick  FROM post p WHERE p.threadid = $1 and p.path[1] IN ( 
							SELECT p2.path[1] FROM post p2 WHERE p2.threadid = $1 AND p2.parent = 0 and p2.path[1] > 
							(SELECT p3.path[1] from post p3 where p3.id = $2) ORDER BY p2.path) ORDER BY p.path `, thread.Id, since, limit)
					}
				}
			} else {
				if limit != -1 {
					if desc {
						rows, err = transaction.Query(`SELECT  id, message, created, parent, isEdited, usernick  FROM post WHERE threadid = $1 and path[1] IN ( 
							SELECT path[1] FROM post WHERE threadid = $1 GROUP BY path[1] 
							ORDER BY path[1] DESC LIMIT $2) ORDER BY path[1] DESC, path`, thread.Id, limit)
					} else {
						rows, err = transaction.Query(`SELECT  id, message, created, parent, isEdited, usernick  FROM post WHERE threadid = $1 and path[1] IN ( 
							SELECT path[1] FROM post WHERE threadid = $1 GROUP BY path[1] 
							ORDER BY path[1] LIMIT $2) ORDER BY path`, thread.Id, limit)
					}
				} else {
					if desc {
						rows, err = transaction.Query(`SELECT  id, message, created, parent, isEdited, usernick  FROM post WHERE threadid = $1 and path[1] IN ( 
							SELECT path[1] FROM post WHERE threadid = $1 GROUP BY path[1] 
							ORDER BY path[1] DESC) ORDER BY path[1] DESC, path`, thread.Id)
					} else {
						rows, err = transaction.Query(`SELECT  id, message, created, parent, isEdited, usernick  FROM post WHERE threadid = $1 and path[1] IN ( 
							SELECT path[1] FROM post WHERE threadid = $1 GROUP BY path[1] 
							ORDER BY path[1]) ORDER BY path`, thread.Id)
					}
				}
			}

		}
	} else {
		if sinceStr != ""{
			since, _ := strconv.Atoi(sinceStr)
			if limit != -1 {
				if desc {
					rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 and id  < $2 ORDER BY id DESC LIMIT $3 `, thread.Id, since, limit)
				} else {
					rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 and id  > $2 ORDER BY id LIMIT $3 `, thread.Id, since, limit)
				}
			} else {
				if desc {
					rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 and id < $2 ORDER BY id DESC `, thread.Id, since)
				} else {
					rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 and id  > $2 ORDER BY id `, thread.Id, since)
				}
			}
		} else {
			if limit != -1 {
				if desc {
					rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 ORDER BY id DESC LIMIT $2 `, thread.Id, limit)
				} else {
					rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1  ORDER BY id LIMIT $2 `, thread.Id, limit)
				}
			} else {
				if desc {
					rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 ORDER BY id DESC `, thread.Id)
				} else {
					rows, err = transaction.Query(`SELECT id, message, created, parent, isEdited, usernick FROM post WHERE threadid = $1 ORDER BY id `, thread.Id)
				}
			}
		}
	}

	if err != nil {
		//log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return postsFound, fmt.Errorf("can't find posts with thread %s", thread.Id)
	}
	for rows.Next() {
		var postFound Post
		err = rows.Scan(&postFound.Id, &postFound.Message, &postFound.Date, &postFound.Parent, &postFound.Edited, &postFound.User)
		if err != nil {
			//log.Println(err)
			//errRollback := transaction.Rollback()
			if err != nil {
				//log.Println(errRollback)
			}
			return postsFound, err
		}
		postFound.Forum = forumSlug
		postFound.Thread = thread.Id
		postsFound = append(postsFound, postFound)
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return postsFound, nil
}

func (post *Post) GetPost() (error, string, string) {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		//log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			//log.Println(errRollback)
		}
		return err, post.User, post.Forum
	}
	rows := transaction.QueryRow("SELECT id, message, created, parent, isEdited, usernick, threadid, forumslug FROM post WHERE id = $1", post.Id)
	err = rows.Scan(&post.Id, &post.Message, &post.Date, &post.Parent, &post.Edited, &post.User, &post.Thread, &post.Forum)
	if err != nil {
		//log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			//log.Println(errRollback)
		}
		return fmt.Errorf("can't find post with id %d", post.Id), post.User, post.Forum
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return nil, post.User, post.Forum
}


func (post *Post) GetPostRelated(related string) (PostRelated, error) {
	var relatedPost PostRelated
	var userStr string
	var forumSlug string
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		//log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			//log.Println(errRollback)
		}
		return relatedPost, err
	}
	err, userStr, forumSlug = post.GetPost()
	if err != nil {
		//log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			//log.Println(errRollback)
		}
		return relatedPost, err
	}
	relatedPost.Post = post
	if strings.Contains(related, "user") {
		var newUser User
		newUser.Nickname = userStr
		rows := transaction.QueryRow("SELECT * FROM forum_user WHERE nickname = $1", newUser.Nickname)
		err = rows.Scan(&newUser.Id, &newUser.Nickname, &newUser.Email, &newUser.Fullname, &newUser.About)
		if err != nil {
			//log.Println(err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				//log.Println(errRollback)
			}
			return relatedPost, err
		}
		relatedPost.User = &newUser
	}
	if strings.Contains(related, "forum") {
		var newForum Forum
		rows := transaction.QueryRow("SELECT * FROM forum WHERE slug = $1", forumSlug)
		err = rows.Scan(&newForum.Id, &newForum.Slug, &newForum.Title, &newForum.Posts, &newForum.Threads, &newForum.User)
		if err != nil {
			//log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return relatedPost, fmt.Errorf("can't find forum with slug %s", forumSlug)
		}
		relatedPost.Forum = &newForum
	}
	if strings.Contains(related, "thread") {
		var newThread Thread
		rows := transaction.QueryRow("SELECT * FROM thread WHERE id = $1", post.Thread)
		err = rows.Scan(&newThread.Id, &newThread.Slug, &newThread.Date, &newThread.Title, &newThread.Message, &newThread.Votes,  &newThread.Forum, &newThread.User)
		if err != nil {
			//log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return relatedPost, fmt.Errorf("can't find thread with id %d", post.Thread)
		}
		relatedPost.Thread = &newThread
	}

	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return relatedPost, nil
}



func (post *Post) UpdatePost() error {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		//log.Println(err)
	}
	newMessage := post.Message
	rows := transaction.QueryRow("SELECT id, created, parent, isEdited, message, usernick, threadid, forumslug FROM post WHERE id = $1", post.Id)
	err = rows.Scan(&post.Id, &post.Date, &post.Parent, &post.Edited, &post.Message, &post.User, &post.Thread, &post.Forum)
	if err != nil {
		//log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return fmt.Errorf("can't find post with id %d", post.Id)
	}
	if len(newMessage) != 0 {
		rows = transaction.QueryRow(`UPDATE post SET (message, isedited) = (coalesce($2, message), $2 IS NOT NULL AND $2 <> message) 
			WHERE id = $1 RETURNING message, isedited`, post.Id, newMessage)
		err = rows.Scan(&post.Message, &post.Edited)
		if err != nil {
			//log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return err
		}
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}