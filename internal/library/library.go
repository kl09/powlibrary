package library

import (
	"context"
	"math/rand"
)

// Library is a collection of motivational quotes.
type Library struct {
	collection []string
}

// NewLibrary creates a new Library instance.
func NewLibrary() *Library {
	return &Library{
		collection: []string{
			`Don't settle for just anything, make sure you really want it.`,
			`There are no shortcuts to success.`,
			`Failure is only a setback, not an obstacle.`,
			`You must continue to learn and grow as a person.`,
			`Don't let other people tell you what to do, you are in control of your own life.`,
			`Make it happen.`,
			`Persistence is key.`,
			`Always be patient.`,
			`Don’t wait, just do, it just might be too late.`,
			`Know that you are worth something, don't let anyone tell you different.`,
			`You can't please everyone, there is always someone that won't like the choices you make.`,
			`Admit that you will be wrong sometimes.`,
			`There is no such thing as perfection, everyone has flaws.`,
			`Always be the bigger person.`,
			`Giving up is the enemy.`,
			`Stand up for yourself.`,
			`Don't dwell on the past.`,
			`Make sure to live in the moment.`,
			`Your past mistakes don't define you.`,
			`It doesn't matter what people think of you, it's what you think of yourself that counts.`,
			`Don't let your emotions get the best of you.`,
			`Don't take advantage of others.`,
			`Your future hasn’t been written yet, so make sure it’s a good one.`,
			`Don’t let hatred consume you.`,
			`Live your life with confidence.`,
			`Don’t live your life with doubts.`,
			`There will be some things you can’t control.`,
			`Before you think about quitting, think about why you started.`,
			`Personality is more important than looks.`,
			`Nothing is impossible if you believe you can do it.`,
			`Don’t waste your voice on people who won’t listen.`,
			`Dreams come true if you work for it.`,
			`Positivity gets you places, while negativity brings you down.`,
			`Sometimes you just got to relax.`,
			`Learning never stops; it continues as you grow older.`,
			`It takes common sense to just walk away.`,
			`Step out of your comfort zone.`,
			`If you don’t want to win, don’t try.`,
			`Nothing worthwhile ever came easy.`,
			`You can’t be afraid to fail, if you don’t give it a shot.`,
			`Keep your head held high.`,
			`Luxury isn’t everything, money doesn’t buy happiness.`,
			`Keep moving forward.`,
		},
	}
}

// GetRandomQuote returns a random quote from the collection.
func (l *Library) GetRandomQuote(ctx context.Context) string {
	return l.collection[rand.Intn(len(l.collection))]
}
