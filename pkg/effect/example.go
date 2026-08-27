package effect

import (
	"errors"
)

func effectExample() {
	var err error = errors.New("error message")
	var ef, prior Effect

	ef = Bless("data")
	ef = Curse("error message")
	ef = Cursed(err)
	ef = Judge("data")
	ef = Judge(err)
	ef = Choose("data", err)

	prior = Bless("abc")
	ef = Bless("efg").
		PriorAs(prior).
		NameAs("3rd, 4th, and 5th letters").
		Dispel()

	_ = ef.Name()      // "3rd, 4th, and 5th letters"
	_ = ef.Prior()     // Effect{result: "abc"}
	_ = ef.Result()    // "efg"
	_ = ef.Cursed()    // false
	_ = ef.Dispelled() // true
	_, _ = ef.Values() // "efg", nil

	ef = Curse("error message")

	_ = ef.Cursed()    // true
	_ = ef.Error()     // error{message:"error message"}
	_, _ = ef.Values() // nil, error{message:"error message"}
}
