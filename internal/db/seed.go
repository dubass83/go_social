package db

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	lorem "github.com/derektata/lorem/ipsum"
	"github.com/dubass83/go_social/internal/store"
	"github.com/rs/zerolog/log"
)

func Seed(store *store.Storage, num int) {
	lg := NewLoremGenerator()
	ctx := context.Background()

	titles := TitleGenerator(lg, num)
	contents := ContentGenerator(lg, num)
	tags := TitleGenerator(lg, 20)
	cs := CommentsGenerator(lg, num*2)

	users := generateUsers(num)
	// Seed users
	for _, user := range users {
		err := store.User.Create(ctx, user)
		if err != nil {
			log.Error().Err(err).Msgf("error seeding user: %s", user.Username)
		}
	}

	posts := generatePosts(num, users, titles, contents, tags)
	// Seed posts
	for _, post := range posts {
		err := store.Post.Create(ctx, post)
		if err != nil {
			log.Error().Err(err).Msgf("error seeding post: %s", post.Title)
		}
	}

	comments := generateComments(cs, users, posts)
	// Seed comments
	for _, comment := range comments {
		err := store.Comment.Create(ctx, comment)
		if err != nil {
			log.Error().Err(err).Msgf("error seeding comment: %s", comment.Content)
		}
	}

}

func generateComments(cs []string, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, len(cs))
	for i := range cs {
		usr := users[rand.Intn(len(users))]
		comments[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  usr.ID,
			Content: cs[i],
			User:    *usr,
		}
	}
	return comments
}

func generatePosts(n int, users []*store.User, titles, contents, tags []string) []*store.Post {
	posts := make([]*store.Post, n)
	for i := range n {
		posts[i] = &store.Post{
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			UserID:  users[rand.Intn(len(users))].ID,
			Version: 0,
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)
	runes := generateRuneList()
	for i := range num {
		usr := randonUserName(runes, 8)
		users[i] = &store.User{
			Username: usr + fmt.Sprintf("%d", i),
			Email:    usr + fmt.Sprintf("%d@example.me", i),
			Password: fmt.Sprintf("some_pass_%d", i),
		}
	}
	return users
}

func randonUserName(r []rune, num int) string {
	result := make([]rune, num)
	for i := range num {
		result[i] = r[rand.Intn(len(r))]
	}
	return string(result)

}

func generateRuneList() []rune {
	// Lowercase English letters
	var lowercase []rune
	for i := 'a'; i <= 'z'; i++ {
		lowercase = append(lowercase, i)
	}
	return lowercase
}

func NewLoremGenerator() *lorem.Generator {
	g := lorem.NewGenerator()
	g.WordsPerSentence = 10     // Customize how many words per sentence
	g.SentencesPerParagraph = 5 // Customize how many sentences per paragraph
	g.CommaAddChance = 3        // Customize the chance of a comma being added to a sentence
	return g
}

// TitleGenerator generate rundom number of titles
func TitleGenerator(g *lorem.Generator, num int) []string {
	titles := strings.Split(g.Generate(num), " ")
	result := []string{}
	for i := range titles {
		// Trim commas, and periods
		cleaned := strings.Trim(titles[i], ",.")
		if cleaned != "" {
			result = append(result, cleaned)
		}
	}
	return result

}

// ContentGenerator generate rundom number of contents
func ContentGenerator(g *lorem.Generator, num int) []string {
	contents := make([]string, num)
	for i := range num {
		contents[i] = g.GenerateParagraphs(3)
	}
	return contents
}

func CommentsGenerator(g *lorem.Generator, num int) []string {
	comments := make([]string, num)
	for i := range num {
		comments[i] = g.Generate(10)
	}
	return comments
}
