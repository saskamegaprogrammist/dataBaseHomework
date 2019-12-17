package models

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
	"strconv"
	"time"
)

type Post struct {
	Id int32 `json:"id"`
	Message string `json:"message"`
	Date time.Time `json:"created"`
	Parent int32 `json:"parent"`
	Edited bool `json:"isEdited"`
	User string `json:"author"`
	Forum string `json:"forum"`
	Thread int32 `json:"thread"`
}

func GetPostsByThread(params utils.SearchParams, thread Thread) ([]Post, error) {
	postsFound := make([]Post, 0)
	var postsUsers []int32
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
	var actualId int32
	var forumId int32
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
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 and id  < $2 ORDER BY created DESC LIMIT $3) as a ORDER BY id DESC", thread.Id, since, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post as a WHERE threadid = $1 and id  > $2 ORDER BY created LIMIT $3) as a ORDER BY id ", thread.Id, since, params.Limit)
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
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post as a WHERE threadid = $1 ORDER BY created DESC LIMIT $2) as a ORDER BY id DESC", thread.Id, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT * FROM (SELECT id, message, created, parent, isEdited, userid FROM post as a WHERE threadid = $1 ORDER BY created LIMIT $2) as a ORDER BY id", thread.Id, params.Limit)
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
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND id > $2 ORDER BY id DESC) as a "+
						"LEFT JOIN (SELECT * FROM post order by path DESC) as b ON ARRAY [a.id] &&  b.path::integer[] LIMIT $3", thread.Id, since, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND id < $2 ORDER BY id) as a "+
							"LEFT JOIN (SELECT * FROM post order by path) as b ON ARRAY [a.id] &&  b.path::integer[] LIMIT $3", thread.Id, since, params.Limit)
					}
				} else {
					if params.Decs {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND id > $2 ORDER BY id DESC) as a "+
							"LEFT JOIN (SELECT * FROM post order by path DESC) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id, since)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND id < $2 ORDER BY id) as a "+
							"LEFT JOIN (SELECT * FROM post order by path) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id, since)
					}
				}
			} else {
				if params.Limit != -1 {
					if params.Decs {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id DESC) as a "+
							"LEFT JOIN (SELECT * FROM post order by path DESC) as b ON ARRAY [a.id] &&  b.path::integer[] LIMIT $2", thread.Id, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id) as a "+
							"LEFT JOIN (SELECT * FROM post order by path) as b ON ARRAY [a.id] &&  b.path::integer[] LIMIT $2", thread.Id, params.Limit)
					}
				} else {
					if params.Decs {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id DESC) as a "+
							"LEFT JOIN (SELECT * FROM post order by path DESC) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id) as a "+
							"LEFT JOIN (SELECT * FROM post order by path) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id)
					}
				}
			}
		case "parent_tree":
			if params.Since != ""{
				since, _ := strconv.Atoi(params.Since)
				if params.Limit != -1 {
					if params.Decs {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND id < $2 ORDER BY id DESC LIMIT $3) as a "+
							"LEFT JOIN (SELECT * FROM post order by path ) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id, since, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND id > $2 ORDER BY id LIMIT $3) as a "+
							"LEFT JOIN (SELECT * FROM post order by path) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id, since, params.Limit)
					}
				} else {
					if params.Decs {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND id < $2 ORDER BY id DESC) as a "+
							"LEFT JOIN (SELECT * FROM post order by path ) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id, since)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND id > $2 ORDER BY id) as a "+
							"LEFT JOIN (SELECT * FROM post order by path) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id, since)
					}
				}
			} else {
				if params.Limit != -1 {
					if params.Decs {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id DESC LIMIT $2) as a "+
							"LEFT JOIN (SELECT * FROM post order by path ) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id, params.Limit)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 AND id > $2 ORDER BY id LIMIT $2) as a "+
							"LEFT JOIN (SELECT * FROM post order by path) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id, params.Limit)
					}
				} else {
					if params.Decs {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id DESC) as a "+
							"LEFT JOIN (SELECT * FROM post order by path ) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id)
					} else {
						rows, err = transaction.Query("SELECT b.id, b.message, b.created, b.parent, b.isEdited, b.userid FROM (SELECT id FROM post WHERE array_length(path, 1) = 1 AND threadid = $1 ORDER BY id) as a "+
							"LEFT JOIN (SELECT * FROM post order by path) as b ON ARRAY [a.id] &&  b.path::integer[] ", thread.Id)
					}
				}
			}

		}
	} else {
		if params.Since != ""{
			since, _ := strconv.Atoi(params.Since)
			if params.Limit != -1 {
				if params.Decs {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 and id  > $2 ORDER BY id DESC LIMIT $3 ", thread.Id, since, params.Limit)
				} else {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 and id  < $2 ORDER BY id LIMIT $3 ", thread.Id, since, params.Limit)
				}
			} else {
				if params.Decs {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 and id  > $2 ORDER BY id DESC ", thread.Id, since)
				} else {
					rows, err = transaction.Query("SELECT id, message, created, parent, isEdited, userid FROM post WHERE threadid = $1 and id  < $2 ORDER BY id ", thread.Id, since)
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
		var userId int32
		err = rows.Scan(&postFound.Id, &postFound.Message, &postFound.Date, &postFound.Parent, &postFound.Edited, &userId)
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if err != nil {
				log.Fatalln(errRollback)
			}
			return postsFound, err
		}
		postFound.Forum = forumSlug
		postFound.Thread = thread.Id
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