// Package wizzard provides a simple ORM-based interface to
// managing SQLite databases with Go struct's.
//
// # Rincewind the Wizzard
//
// No, wizzard is not a mis-spelliing of wizard! Well ok
// it is, but it's intentional. Wizzard is an old English
// rock band that are famous for the song "I Wish It Could
// be Christmas Everyday". UK supermarkets blast the song
// on repeat every Christmas along with Mariah Carey's "All
// I Want for Christmas is You", The Progues "Fairytale Of
// New York", and Wham's "Last Christmas".
//
// This package is actually named after the words
// haphazardly crafted in sequins on Rincewind's hat.
// Rincewind is a failed wizard of the Discworld who lives
// in a never-ending series of interesting times. They say
// that if you've never read a Discworld novel you're not
// really real. But since very few people are really real,
// 'they' included, I wouldn't worry about it. But I do
// recommend reading the entire series.
//
// # TODO
//   - Add Select function to Model that accepts string
//     WHERE clause followed by series of arguments to
//     SQL parameters and returns all results.
//   - Add SelectFirst function to Model that accepts
//     string  WHERE clause followed by series of arguments to
//     SQL parameters and returns a single result.
//   - Add Delete function to Model that accepts string
//     WHERE clause followed by series of arguments to
//     SQL parameters.
package wizzard
