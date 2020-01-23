package models
const  (
	userLimitSinceDesc = `SELECT about, fullname, nickname, email FROM forum_user 
	JOIN (SELECT DISTINCT usernick as merge_nick FROM forum_user_new 
	WHERE forumslug = $1 ) 
	as u ON u.merge_nick = forum_user.nickname WHERE nickname   < $2 ORDER BY (nickname ) DESC LIMIT $3`

	userLimitDesc = `SELECT about, fullname, nickname, email FROM forum_user 
		JOIN (SELECT DISTINCT usernick as merge_nick FROM forum_user_new 
		WHERE forumslug = $1 ) 
		as u ON u.merge_nick = forum_user.nickname  ORDER BY (nickname ) DESC LIMIT $2`

	userSinceDesc = `SELECT about, fullname, nickname, email FROM forum_user 
	JOIN (SELECT DISTINCT usernick as merge_nick FROM forum_user_new 
	WHERE forumslug = $1 ) 
	as u ON u.merge_nick = forum_user.nickname  WHERE nickname  < $2 ORDER BY (nickname ) DESC`

	userDesc =	`SELECT about, fullname, nickname, email FROM forum_user 
	JOIN (SELECT DISTINCT usernick as merge_nick FROM forum_user_new 
	WHERE forumslug = $1 ) 
	as u ON u.merge_nick = forum_user.nickname ORDER BY (nickname ) DESC`

	userLimitSince = `SELECT about, fullname, nickname, email FROM forum_user 
	JOIN (SELECT DISTINCT usernick as merge_nick FROM forum_user_new 
	WHERE forumslug = $1 ) 
	as u ON u.merge_nick = forum_user.nickname WHERE nickname   > $2 ORDER BY (nickname ) LIMIT $3`

	userLimit = `SELECT about, fullname, nickname, email FROM forum_user 
	JOIN (SELECT DISTINCT usernick as merge_nick FROM forum_user_new 
	WHERE forumslug = $1 ) 
	as u ON u.merge_nick = forum_user.nickname  ORDER BY (nickname ) LIMIT $2`

	userSince = `SELECT about, fullname, nickname, email FROM forum_user 
	JOIN (SELECT DISTINCT usernick as merge_nick FROM forum_user_new 
	WHERE forumslug = $1 ) 
	as u ON u.merge_nick = forum_user.nickname  WHERE nickname   > $2 ORDER BY (nickname )`

	userSimple = `SELECT about, fullname, nickname, email FROM forum_user 
	JOIN (SELECT DISTINCT usernick as merge_nick FROM forum_user_new 
	WHERE forumslug = $1 ) 
	as u ON u.merge_nick = forum_user.nickname  ORDER BY (nickname )`

_)