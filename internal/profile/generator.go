package profile

import (
	"fmt"
	"math/rand"
	"net/url"
)

var adjectives = []string{
	"Bouncy", "Fizzy", "Wobbly", "Zesty", "Giggly", "Bubbly",
	"Wiggly", "Peppy", "Jolly", "Snappy", "Zippy", "Sproingy",
}

var nouns = []string{
	"Pickle", "Muffin", "Pancake", "Jellybean", "Cupcake", "Waffle",
	"Buttercup", "Marshmallow", "Pumpkin", "Doodle", "Sprinkle", "Noodle",
}

func GenerateRandomName() string {
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	return fmt.Sprintf("%s %s", adj, noun)
}

func GenerateRandomAvatarURL() string {
	colors := []string{"d1d4f9", "c0aede", "b6e3f4"}

	// Avataaars options (excluding mouth "vomit")
	topTypes := []string{"noHair", "eyepatch", "hat", "hijab", "turban", "winterHat1", "winterHat2", "winterHat3", "winterHat4", "longHairBigHair", "longHairBob", "longHairBun", "longHairCurly", "longHairCurvy", "LongHairDreads", "longHairFrida", "longHairFro", "longHairFroBand", "longHairNotTooLong", "longHairShavedSides", "longHairMiaWallace", "longHairStraight", "longHairStraight2", "longHairStraightStrand", "shortHairDreads01", "shortHairDreads02", "shortHairFrizzle", "shortHairShaggyMullet", "shortHairShortCurly", "shortHairShortFlat", "shortHairShortRound", "shortHairShortWaved", "shortHairSides", "shortHairTheCaesar", "shortHairTheCaesarSidePart"}

	eyeTypes := []string{"closed", "cry", "default", "xDizzy", "eyeRoll", "happy", "hearts", "side", "squint", "surprised", "wink", "winkWacky"}

	eyebrowTypes := []string{"angry", "angryNatural", "default", "defaultNatural", "flatNatural", "raisedExcited", "raisedExcitedNatural", "sadConcerned", "sadConcernedNatural", "unibrowNatural", "upDown", "upDownNatural"}

	mouthTypes := []string{"concerned", "default", "disbelief", "eating", "grimace", "sad", "screamOpen", "serious", "smile", "tongue", "twinkle"}

	skinColors := []string{"tanned", "yellow", "pale", "light", "brown", "darkBrown", "black"}

	baseURL := "https://api.dicebear.com/7.x/avataaars-neutral/svg"
	params := url.Values{}

	params.Add("backgroundColor", colors[rand.Intn(len(colors))])
	params.Add("top", topTypes[rand.Intn(len(topTypes))])
	params.Add("eyes", eyeTypes[rand.Intn(len(eyeTypes))])
	params.Add("eyebrows", eyebrowTypes[rand.Intn(len(eyebrowTypes))])
	params.Add("mouth", mouthTypes[rand.Intn(len(mouthTypes))])
	params.Add("skin", skinColors[rand.Intn(len(skinColors))])

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}
