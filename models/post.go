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


func GetPostsByThread(params utils.SearchParams, thread Thread) ([]Post, error) {
	postsFound := make([]Post, 0)
	var postsUsers []int
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if err != nil {
			log.Fatalln(errRollback)
		}
		return postsFound, err
	}
	var actualId int
	var forumId int
	var forumSlug string
	if thread.Id != 0 {
		rows := transaction.QueryRow("SELECT id, forumid FROM thread WHERE id = $1", thread.Id)
		err = rows.Scan(&actualId, &forumId)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return postsFound, fmt.Errorf("can't find thread with id %d", thread.Id)
		}
	} else {
		rows := transaction.QueryRow("SELECT id, forumid FROM thread WHERE slug = $1", thread.Slug)
		err = rows.Scan(&actualId, &forumId)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return postsFound, fmt.Errorf("can't find thread with slug %s", thread.Slug)
		}
		thread.Id = actualId
	}

	row := transaction.QueryRow("SELECT slug FROM forum WHERE id = $1", forumId)
	err = row.Scan(&forumSlug)
	if err != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return postsFound, err
	}

	var rows *pgx.Rows
	if params.Sort != "" {
		switch params.Sort {
		case "flat":
			if params.Since != ""{
				since, _ := strconv.Atoi(params.Since)
				if params.Limit != -1 {
					if params.Decs {
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 and id  < $2 ORDER BY created DESC ) as a ORDER BY id DESC LIMIT $3", thread.Id, since, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post as a WHERE threadid = $1 and id  > $2 ORDER BY created ) as a ORDER BY id LIMIT $3 ", thread.Id, since, params.Limit)
					}
				} else {
					if params.Decs {
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post as a WHERE threadid = $1 and id  < $2 ORDER BY created DESC) as a ORDER BY id DESC ", thread.Id, since)
					} else {
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post as a WHERE threadid = $1 and id  > $2 ORDER BY created) as a ORDER BY id ", thread.Id, since)
					}
				}
			} else {
				if params.Limit != -1 {
					if params.Decs {
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post as a WHERE threadid = $1 ORDER BY created DESC) as a ORDER BY id DESC LIMIT $2", thread.Id, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post as a WHERE threadid = $1 ORDER BY created) as a ORDER BY id LIMIT $2", thread.Id, params.Limit)
					}
				} else {
					if params.Decs {
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post as a WHERE threadid = $1 ORDER BY created DESC) as a ORDER BY id DESC", thread.Id)
					} else {
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post as a WHERE threadid = $1 ORDER BY created) as a ORDER BY id ", thread.Id)
					}
				}
			}

		case "tree":
			if params.Since != ""{
				since, _ := strconv.Atoi(params.Since)
				if params.Limit != -1 {
					if params.Decs {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id DESC) as a "+
						"LEFT JOIN (SELECT * FROM post WHERE path < (SELECT path FROM post WHERE id = $2) order by path DESC) as b ON ARRAY [a.id] &&  b.path::integer[] WHERE b.id is not null LIMIT $3 ", thread.Id, since, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id) as a "+
							"LEFT JOIN (SELECT * FROM post WHERE path > (SELECT path FROM post WHERE id = $2) order by path) as b ON ARRAY [a.id] &&  b.path::integer[] WHERE b.id is not null LIMIT $3 ", thread.Id, since, params.Limit)
					}
				} else {
					if params.Decs {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id DESC) as a "+
							"LEFT JOIN (SELECT * FROM post WHERE path < (SELECT path FROM post WHERE id = $2) order by path DESC) as b ON ARRAY [a.id] &&  b.path::integer[] WHERE b.id is not null ", thread.Id, since)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id) as a "+
							"LEFT JOIN (SELECT * FROM post WHERE path > (SELECT path FROM post WHERE id = $2) order by path) as b ON ARRAY [a.id] &&  b.path::integer[] WHERE b.id is not null ", thread.Id, since)
					}
				}
			} else {
				if params.Limit != -1 {
					if params.Decs {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id DESC) as a "+
							"LEFT JOIN (SELECT * FROM post order by path DESC) as b ON ARRAY [a.id] &&  b.path::integer[] WHERE b.id is not null LIMIT $2", thread.Id, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id) as a "+
							"LEFT JOIN (SELECT * FROM post order by path) as b ON ARRAY [a.id] &&  b.path::integer[] WHERE b.id is not null LIMIT $2", thread.Id, params.Limit)
					}
				} else {
					if params.Decs {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id DESC) as a "+
							"LEFT JOIN (SELECT * FROM post order by path DESC) as b ON ARRAY [a.id] &&  b.path::integer[] WHERE b.id is not null", thread.Id)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id) as a "+
							"LEFT JOIN (SELECT * FROM post order by path) as b ON ARRAY [a.id] &&  b.path::integer[] WHERE b.id is not null", thread.Id)
					}
				}
			}
		case "parent_tree":
			if params.Since != ""{
				since, _ := strconv.Atoi(params.Since)
				if params.Limit != -1 {
					if params.Decs {
						rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid  FROM post WHERE path[1] IN (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND " +
						"path[1] < (SELECT path[1] FROM post WHERE id = $2) order by id DESC LIMIT $3) AND id IS NOT NULL ORDER BY path[1] DESC, path", thread.Id, since, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid  FROM post WHERE path[1] IN (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND " +
							"path[1] > (SELECT path[1] FROM post WHERE id = $2) order by id LIMIT $3) AND id IS NOT NULL ORDER BY path", thread.Id, since, params.Limit)
					}
				} else {
					if params.Decs {
						rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid  FROM post WHERE path[1] IN (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND " +
							"path[1] < (SELECT path[1] FROM post WHERE id = $2) order by id DESC) AND id IS NOT NULL ORDER BY path[1] DESC, path", thread.Id, since)
					} else {
						rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid  FROM post WHERE path[1] IN (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND " +
							"path[1] > (SELECT path[1] FROM post WHERE id = $2) order by id) AND id IS NOT NULL ORDER BY path", thread.Id, since)
					}
				}
			} else {
				if params.Limit != -1 {
					if params.Decs {
						rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE path[1] IN (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 " +
							"order by id DESC LIMIT $2) AND id IS NOT NULL ORDER BY path[1] DESC, path", thread.Id, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE path[1] IN (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 " +
							"order by id LIMIT $2) AND id IS NOT NULL ORDER BY path", thread.Id, params.Limit)
					}
				} else {
					if params.Decs {
						rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE path[1] IN (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 " +
							"order by id DESC) AND id IS NOT NULL ORDER BY path[1] DESC, path", thread.Id)
					} else {
						rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE path[1] IN (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 " +
							"order by id) AND id IS NOT NULL ORDER BY path", thread.Id)
					}
				}
			}

		}
	} else {
		if params.Since != ""{
			since, _ := strconv.Atoi(params.Since)
			if params.Limit != -1 {
				if params.Decs {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 and id  < $2 ORDER BY id DESC LIMIT $3 ", thread.Id, since, params.Limit)
				} else {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 and id  > $2 ORDER BY id LIMIT $3 ", thread.Id, since, params.Limit)
				}
			} else {
				if params.Decs {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 and id < $2 ORDER BY id DESC ", thread.Id, since)
				} else {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 and id  > $2 ORDER BY id ", thread.Id, since)
				}
			}
		} else {
			if params.Limit != -1 {
				if params.Decs {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 ORDER BY id DESC LIMIT $2 ", thread.Id, params.Limit)
				} else {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1  ORDER BY id LIMIT $2 ", thread.Id, params.Limit)
				}
			} else {
				if params.Decs {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 ORDER BY id DESC ", thread.Id)
				} else {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 ORDER BY id ", thread.Id)
				}
			}
		}
	}

	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return postsFound, fmt.Errorf("can't find posts with thread %s", thread.Id)
	}
	for rows.Next() {
		var postFound Post
		var userId int
		err = rows.Scan(&postFound.Id, &postFound.Message, &postFound.Date, &postFound.Parent, &postFound.Edited, &userId)
		fmt.Println(postFound)
		fmt.Println(params)
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if err != nil {
				log.Println(errRollback)
			}
			return postsFound, err
		}
		postFound.Forum = forumSlug
		postFound.Thread = int(thread.Id)
		postsFound = append(postsFound, postFound)
		postsUsers = append(postsUsers, userId)
	}

	for i:=0; i< len(postsFound); i++ {
		row := transaction.QueryRow("SELECT nickname FROM forum_user WHERE id = $1", postsUsers[i])
		err = row.Scan(&postsFound[i].User)
		if err != nil {
			errRollback := transaction.Rollback()
			if errRollback != nil {
				log.Fatalln(errRollback)
			}
			return postsFound, err
		}

	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return postsFound, nil
}

func (post *Post) GetPost() (error, int, int) {
	var userId int
	var forumId int
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Println(errRollback)
		}
		return err, userId, forumId
	}
	rows := transaction.QueryRow("SELECT id, message, created, parent, isEdited, userid, threadid, forumid FROM post WHERE id = $1", post.Id)
	err = rows.Scan(&post.Id, &post.Message, &post.Date, &post.Parent, &post.Edited, &userId, &post.Thread, &forumId)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Println(errRollback)
		}
		return fmt.Errorf("can't find post with id %d", post.Id), userId, forumId
	}
	rows = transaction.QueryRow("SELECT nickname FROM forum_user WHERE id = $1", userId)
	err = rows.Scan(&post.User)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Println(errRollback)
		}
		return fmt.Errorf("can't find user with id %d", userId), userId, forumId
	}
	rows = transaction.QueryRow("SELECT slug FROM forum WHERE id = $1", forumId)
	err = rows.Scan(&post.Forum)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Println(errRollback)
		}
		return fmt.Errorf("can't find forum with id %d", forumId), userId, forumId
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return nil, userId, forumId
}


func (post *Post) GetPostRelated(related string) (PostRelated, error) {
	var relatedPost PostRelated
	var userId int
	var forumId int
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Println(errRollback)
		}
		return relatedPost, err
	}
	err, userId, forumId = post.GetPost()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Println(errRollback)
		}
		return relatedPost, err
	}
	relatedPost.Post = post
	if strings.Contains(related, "user") {
		var newUser User
		newUser.Id= userId
		rows := transaction.QueryRow("SELECT * FROM forum_user WHERE id = $1", newUser.Id)
		err = rows.Scan(&newUser.Id, &newUser.Nickname, &newUser.Email, &newUser.Fullname, &newUser.About)
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				log.Println(errRollback)
			}
			return relatedPost, err
		}
		relatedPost.User = &newUser
	}
	if strings.Contains(related, "forum") {
		var newForum Forum
		newForum.Id = forumId
		err = newForum.GetForum(post.Forum)
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				log.Println(errRollback)
			}
			return relatedPost, err
		}
		relatedPost.Forum = &newForum
	}
	if strings.Contains(related, "thread") {
		var newThread Thread
		newThread.Id = post.Thread
		err = newThread.GetThread()
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				log.Println(errRollback)
			}
			return relatedPost, err
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
	var userId int
	var forumId int
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
	}
	var postExists int
	var messageExists string
	rows := transaction.QueryRow("SELECT id, message FROM post WHERE id = $1", post.Id)
	err = rows.Scan(&postExists, &messageExists)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return fmt.Errorf("can't find post with id %d", post.Id)
	}
	if post.Message != "" && post.Message != messageExists {
		_, err = transaction.Exec("UPDATE post SET (message, isedited) = ($2, true) WHERE id = $1 ",  post.Id, post.Message)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return err
		}
		post.Edited = true
	}
	rows = transaction.QueryRow("SELECT id, message, created, parent, userid, threadid, forumid FROM post WHERE id = $1", post.Id)
	err = rows.Scan(&post.Id, &post.Message, &post.Date, &post.Parent, &userId, &post.Thread, &forumId)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Println(errRollback)
		}
		return err
	}
	rows = transaction.QueryRow("SELECT nickname FROM forum_user WHERE id = $1", userId)
	err = rows.Scan(&post.User)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Println(errRollback)
		}
		return err
	}
	rows = transaction.QueryRow("SELECT slug FROM forum WHERE id = $1", forumId)
	err = rows.Scan(&post.Forum)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Println(errRollback)
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}