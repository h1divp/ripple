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
	topTypes := []string{"NoHair", "Eyepatch", "Hat", "Hijab", "Turban", "WinterHat1", "WinterHat2", "WinterHat3", "WinterHat4", "LongHairBigHair", "LongHairBob", "LongHairBun", "LongHairCurly", "LongHairCurvy", "LongHairDreads", "LongHairFrida", "LongHairFro", "LongHairFroBand", "LongHairNotTooLong", "LongHairShavedSides", "LongHairMiaWallace", "LongHairStraight", "LongHairStraight2", "LongHairStraightStrand", "ShortHairDreads01", "ShortHairDreads02", "ShortHairFrizzle", "ShortHairShaggyMullet", "ShortHairShortCurly", "ShortHairShortFlat", "ShortHairShortRound", "ShortHairShortWaved", "ShortHairSides", "ShortHairTheCaesar", "ShortHairTheCaesarSidePart"}

	eyeTypes := []string{"Close", "Cry", "Default", "Dizzy", "EyeRoll", "Happy", "Hearts", "Side", "Squint", "Surprised", "Wink", "WinkWacky"}

	eyebrowTypes := []string{"Angry", "AngryNatural", "Default", "DefaultNatural", "FlatNatural", "RaisedExcited", "RaisedExcitedNatural", "SadConcerned", "SadConcernedNatural", "UnibrowNatural", "UpDown", "UpDownNatural"}

	mouthTypes := []string{"Concerned", "Default", "Disbelief", "Eating", "Grimace", "Sad", "ScreamOpen", "Serious", "Smile", "Tongue", "Twinkle"}

	skinColors := []string{"Tanned", "Yellow", "Pale", "Light", "Brown", "DarkBrown", "Black"}

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
