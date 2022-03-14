package wisdom

import "math/rand"

// quotes is list of wise (more or less) quotes.
var quotes = []string{
	`Don't trust random "wise" quotes from 1st google's result`,
	`Be confident in yourself`,
	`Always be looking forward`,
	`Live a life of purpose`,
	`Be brave. Be bold`,
	`Use your time wisely`,
	`Value yourself for who you are`,
	`Hone your skills`,
	`Keep your head up`,
	`Learn to speak well and listen better`,
	`Have fun. You’ll accomplish more`,
	`Build genuine connections`,
	`Give more than you take`,
	`Seek your purpose`,
	`Pique your curiosity`,
	`Search for more meaning`,
	`Unleash your personal momentum`,
	`Focus on the future`,
	`Excel in your own way`,
	`Don’t forget to live`,
}

// randomQuote returns one of wise quotes.
func randomQuote() string {
	return quotes[rand.Int()%len(quotes)]
}
