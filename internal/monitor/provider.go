package monitor

import (
	"strings"

	"github.com/turnage/graw/internal/operator"
)

// PostMonitor returns a monitor for new posts in a subreddit(s).
func PostMonitor(
	op operator.Operator,
	handlePost postHandler,
	subreddits []string,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/r/"+strings.Join(subreddits, "+"),
		handlePost,
		nil,
		nil,
		dir,
	)
}

// UserMonitor returns a monitor for new posts or comments by a user.
func UserMonitor(
	op operator.Operator,
	handlePost postHandler,
	handleComment commentHandler,
	user string,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/user/"+user,
		handlePost,
		handleComment,
		nil,
		dir,
	)
}

// MessageMonitor returns a monitor for new private messages to the bot.
func MessageMonitor(
	op operator.Operator,
	handleMessage messageHandler,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/message/messages",
		nil,
		nil,
		handleMessage,
		dir,
	)
}

// CommentReplyMonitor returns a monitor for new replies to the bot's comments.
func CommentReplyMonitor(
	op operator.Operator,
	handleComment commentHandler,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/message/comments",
		nil,
		handleComment,
		nil,
		dir,
	)
}

// PostReplyMonitor returns a monitor for new replies to the bot's posts.
func PostReplyMonitor(
	op operator.Operator,
	handleComment commentHandler,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/message/selfreply",
		nil,
		handleComment,
		nil,
		dir,
	)
}

// MentionMonitor returns a monitor for new mentions of the bot's username
// across Reddit.
func MentionMonitor(
	op operator.Operator,
	handleComment commentHandler,
	dir Direction,
) (Monitor, error) {
	return baseFromPath(
		op,
		"/message/mentions",
		nil,
		handleComment,
		nil,
		dir,
	)
}