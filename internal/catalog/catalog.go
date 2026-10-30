// Package catalog contains everyday and playful emoji reactions, and the
// Noul question that scores each one as a reply.
package catalog

import (
	"fmt"

	"github.com/amscotti/jevmoji/internal/jev"
)

// Entry is one catalog emoji. Hint describes its meaning and common uses
// so the model can judge whether it fits a message.
type Entry struct {
	Emoji string
	Name  string
	Label string
	Hint  string
}

// Catalog mixes playful reactions with everyday chat essentials. Order
// breaks ties in the suggestion engine, and position defines each question ID.
var Catalog = []Entry{
	{Emoji: "🗑️", Name: "wastebasket", Label: "trash", Hint: "garbage, dumpster, messy eater, sloppy"},
	{Emoji: "🐦", Name: "bird", Label: "seagull", Hint: "bird, seagull, pigeon, flying, chirping"},
	{Emoji: "🦝", Name: "raccoon", Label: "trash panda", Hint: "raccoon, trash, stealing, sneaky, night raid"},
	{Emoji: "🐀", Name: "rat", Label: "rat", Hint: "rat, snitch, dirty, sewer, sneaky"},
	{Emoji: "🪳", Name: "cockroach", Label: "cockroach", Hint: "survival, gross, infestation, unkillable"},
	{Emoji: "🦭", Name: "seal", Label: "seal", Hint: "clapping, cute, ocean, balancing"},
	{Emoji: "🦆", Name: "duck", Label: "duck", Hint: "waddling, pond, quack"},
	{Emoji: "🪿", Name: "goose", Label: "goose", Hint: "honk, chaos, angry bird"},
	{Emoji: "🦃", Name: "turkey", Label: "turkey", Hint: "thanksgiving, silly, gobble"},
	{Emoji: "🐔", Name: "chicken", Label: "chicken", Hint: "scared, coward, nuggets"},
	{Emoji: "🦩", Name: "flamingo", Label: "flamingo", Hint: "pink, fancy, standing"},
	{Emoji: "🐧", Name: "penguin", Label: "penguin", Hint: "cold, formal, sliding"},
	{Emoji: "🦤", Name: "dodo", Label: "dodo", Hint: "extinct, outdated, dead idea"},
	{Emoji: "🦞", Name: "lobster", Label: "lobster", Hint: "fancy, red, boiling"},
	{Emoji: "🐙", Name: "octopus", Label: "octopus", Hint: "multitasking, arms, tentacles"},
	{Emoji: "🦑", Name: "squid", Label: "squid", Hint: "ink, weird, deep sea"},
	{Emoji: "🦀", Name: "crab", Label: "crab", Hint: "salty, sideways, pinchy"},
	{Emoji: "🪼", Name: "jellyfish", Label: "jellyfish", Hint: "drifting, sting, no backbone"},
	{Emoji: "🐸", Name: "frog", Label: "frog", Hint: "tea, weird, hopping"},
	{Emoji: "🐌", Name: "snail", Label: "snail", Hint: "slow, slow reply, sluggish"},
	{Emoji: "🦥", Name: "sloth", Label: "sloth", Hint: "lazy, slow, hanging around"},
	{Emoji: "🐢", Name: "turtle", Label: "turtle", Hint: "slow, shell, hiding"},
	{Emoji: "🦟", Name: "mosquito", Label: "mosquito", Hint: "annoying, buzzing, pest"},
	{Emoji: "🪰", Name: "fly", Label: "fly", Hint: "annoying, garbage, buzzing"},
	{Emoji: "🕷️", Name: "spider", Label: "spider", Hint: "creepy, web, crawling"},
	{Emoji: "🦨", Name: "skunk", Label: "skunk", Hint: "stinks, smelly, bad idea"},
	{Emoji: "🐗", Name: "boar", Label: "boar", Hint: "wild, charging, aggressive"},
	{Emoji: "🫏", Name: "donkey", Label: "donkey", Hint: "stubborn, mule, foolish"},
	{Emoji: "🪨", Name: "rock", Label: "rock", Hint: "dumb, pet rock, stone cold"},
	{Emoji: "🎻", Name: "violin", Label: "tiny violin", Hint: "sad, sympathy, sarcastic, crying"},
	{Emoji: "🤡", Name: "clown", Label: "clown", Hint: "joke, circus, ridiculous, fool"},
	{Emoji: "💀", Name: "skull", Label: "dead", Hint: "dead, dying, hilarious, killed me"},
	{Emoji: "🪦", Name: "gravestone", Label: "gravestone", Hint: "rip, dead, grave, over, finished"},
	{Emoji: "🎪", Name: "circus tent", Label: "circus", Hint: "chaos, clown show, tent"},
	{Emoji: "🤹", Name: "juggler", Label: "juggling", Hint: "multitasking, dropping balls"},
	{Emoji: "🧌", Name: "troll", Label: "troll", Hint: "trolling, baiting, feeding"},
	{Emoji: "🛑", Name: "stop sign", Label: "stop", Hint: "stop, halt, enough, no"},
	{Emoji: "🧃", Name: "juice box", Label: "juice box", Hint: "kid, small, elementary"},
	{Emoji: "🍟", Name: "fries", Label: "fries", Hint: "fast food, greasy, snack, stolen food"},
	{Emoji: "🧀", Name: "cheese", Label: "cheese", Hint: "corny, cheesy, mouse"},
	{Emoji: "🍕", Name: "pizza", Label: "pizza", Hint: "party, slices, late night"},
	{Emoji: "🍔", Name: "hamburger", Label: "burger", Hint: "fast food, hungry, whopper"},
	{Emoji: "🌭", Name: "hot dog", Label: "hot dog", Hint: "cookout, silly, ballpark"},
	{Emoji: "🥪", Name: "sandwich", Label: "sandwich", Hint: "lunch, packed, layers"},
	{Emoji: "🌶️", Name: "hot pepper", Label: "spicy", Hint: "spicy, hot take, burning, fire"},
	{Emoji: "🍌", Name: "banana", Label: "banana", Hint: "slipping, silly, peel"},
	{Emoji: "🥑", Name: "avocado", Label: "avocado", Hint: "toast, expensive, millennial"},
	{Emoji: "🍄", Name: "mushroom", Label: "mushroom", Hint: "growing, weird, forest"},
	{Emoji: "☕", Name: "coffee", Label: "coffee", Hint: "morning, tired, caffeine"},
	{Emoji: "🧊", Name: "ice", Label: "ice", Hint: "cold, frozen, chill, iced"},
	{Emoji: "🍩", Name: "doughnut", Label: "donut", Hint: "sweet, hole, dunkin"},
	{Emoji: "🧁", Name: "cupcake", Label: "cupcake", Hint: "sweet, soft, delicate"},
	{Emoji: "🧇", Name: "waffle", Label: "waffle", Hint: "indecisive, breakfast, waffling"},
	{Emoji: "🍍", Name: "pineapple", Label: "pineapple", Hint: "tropical, pizza debate, sweet"},
	{Emoji: "🎺", Name: "trumpet", Label: "sad trombone", Hint: "wah wah, losing, failure, flop"},
	{Emoji: "🧻", Name: "toilet paper", Label: "toilet paper", Hint: "crap, roll, bathroom"},
	{Emoji: "🫠", Name: "melting face", Label: "melting", Hint: "overwhelmed, hot, dying inside"},
	{Emoji: "🙃", Name: "upside-down face", Label: "sarcasm", Hint: "sarcastic, ironic, fine, dead inside"},
	{Emoji: "😤", Name: "huffing face", Label: "huffing", Hint: "angry, mad, steam, frustrated"},
	{Emoji: "💅", Name: "nail polish", Label: "nails", Hint: "sassy, unbothered, done, fancy"},
	{Emoji: "🤏", Name: "pinching hand", Label: "small", Hint: "tiny, little, weak, short"},
	{Emoji: "👑", Name: "crown", Label: "crown", Hint: "king, queen, best, respect, royalty"},
	{Emoji: "🔥", Name: "fire", Label: "fire", Hint: "lit, awesome, hot"},
	{Emoji: "🧠", Name: "brain", Label: "brain", Hint: "smart, big brain, thinking"},
	{Emoji: "🤓", Name: "nerd face", Label: "nerd", Hint: "geek, glasses, actually"},
	{Emoji: "🥱", Name: "yawning face", Label: "yawn", Hint: "boring, tired, sleepy, unimpressed"},
	{Emoji: "🙄", Name: "eye roll", Label: "eye roll", Hint: "annoyed, whatever, duh, rolling eyes"},
	{Emoji: "🤦", Name: "facepalm", Label: "facepalm", Hint: "obvious, mistake, stupid"},
	{Emoji: "🤷", Name: "shrug", Label: "shrug", Hint: "whatever, no idea, don't care"},
	{Emoji: "😵‍💫", Name: "spiral eyes", Label: "spiral", Hint: "dizzy, spiral, confused, overwhelmed"},
	{Emoji: "🫡", Name: "salute", Label: "salute", Hint: "respect, yes, acknowledging"},
	{Emoji: "💨", Name: "dash", Label: "dash", Hint: "gone, fast, poof, vanished"},

	// Greetings, acknowledgments, and gestures.
	{Emoji: "👋", Name: "waving hand", Label: "wave", Hint: "hello, hi, greeting, goodbye, welcome"},
	{Emoji: "👍", Name: "thumbs up", Label: "thumbs up", Hint: "yes, approval, agreed, sounds good, acknowledged"},
	{Emoji: "👎", Name: "thumbs down", Label: "thumbs down", Hint: "no, disapproval, reject, not good"},
	{Emoji: "👏", Name: "clapping hands", Label: "applause", Hint: "well done, congratulations, achievement, praise"},
	{Emoji: "🙌", Name: "raising hands", Label: "hooray", Hint: "celebration, finally, success, great news"},
	{Emoji: "🙏", Name: "folded hands", Label: "thanks", Hint: "thank you, gratitude, please, appreciation"},
	{Emoji: "🤝", Name: "handshake", Label: "handshake", Hint: "agreement, deal, teamwork, partnership"},
	{Emoji: "✌️", Name: "victory hand", Label: "peace", Hint: "peace, goodbye, relaxed, all good"},
	{Emoji: "🤞", Name: "crossed fingers", Label: "fingers crossed", Hint: "hope, good luck, wish, hoping it works"},
	{Emoji: "👌", Name: "OK hand", Label: "okay", Hint: "okay, perfect, works for me, just right"},
	{Emoji: "🤘", Name: "sign of the horns", Label: "rock on", Hint: "awesome, rock music, enthusiasm, keep going"},
	{Emoji: "💪", Name: "flexed biceps", Label: "strong", Hint: "strength, effort, you got this, pushing through"},
	{Emoji: "🫶", Name: "heart hands", Label: "heart hands", Hint: "love, heartfelt thanks, care, support"},
	{Emoji: "👐", Name: "open hands", Label: "open hands", Hint: "welcome, openness, offering help, receiving"},
	{Emoji: "✋", Name: "raised hand", Label: "raised hand", Hint: "volunteer, question, stop, count me in"},
	{Emoji: "👊", Name: "oncoming fist", Label: "fist bump", Hint: "solidarity, nice job, encouragement, agreed"},
	{Emoji: "🤙", Name: "call me hand", Label: "call me", Hint: "phone call, call me, hang loose, relaxed"},
	{Emoji: "👉", Name: "backhand index pointing right", Label: "this", Hint: "this, exactly, pointing something out, direction"},

	// Joy, affection, humor, and surprise.
	{Emoji: "😀", Name: "grinning face", Label: "grin", Hint: "happy, cheerful, pleased, good news"},
	{Emoji: "😃", Name: "grinning face with big eyes", Label: "big smile", Hint: "excited, delighted, yay, happy greeting"},
	{Emoji: "😄", Name: "grinning face with smiling eyes", Label: "cheerful", Hint: "happy, amused, cheerful, great"},
	{Emoji: "😁", Name: "beaming face", Label: "beaming", Hint: "proud, delighted, victory, broad smile"},
	{Emoji: "😆", Name: "grinning squinting face", Label: "laughing", Hint: "funny, silly, laughing, amused"},
	{Emoji: "😂", Name: "face with tears of joy", Label: "tears of joy", Hint: "lol, hilarious, laughing, cracking up"},
	{Emoji: "🤣", Name: "rolling on the floor laughing", Label: "rolling", Hint: "extremely funny, absurd, cannot stop laughing"},
	{Emoji: "😊", Name: "smiling face with smiling eyes", Label: "warm smile", Hint: "friendly, warm greeting, pleased, grateful"},
	{Emoji: "🙂", Name: "slightly smiling face", Label: "smile", Hint: "polite hello, friendly acknowledgment, small smile"},
	{Emoji: "😉", Name: "winking face", Label: "wink", Hint: "playful, teasing, inside joke, knowing"},
	{Emoji: "😍", Name: "smiling face with heart eyes", Label: "heart eyes", Hint: "love it, adore, beautiful, amazing"},
	{Emoji: "🤩", Name: "star-struck face", Label: "star-struck", Hint: "wow, impressed, excited, amazing"},
	{Emoji: "🥳", Name: "partying face", Label: "party", Hint: "celebration, birthday, success, congratulations"},
	{Emoji: "😎", Name: "smiling face with sunglasses", Label: "cool", Hint: "cool, confident, stylish, effortless"},
	{Emoji: "🤗", Name: "hugging face", Label: "hug", Hint: "warm welcome, hug, comfort, friendly support"},
	{Emoji: "🥰", Name: "smiling face with hearts", Label: "loving", Hint: "affection, feeling loved, heartfelt gratitude"},
	{Emoji: "😇", Name: "smiling face with halo", Label: "angel", Hint: "innocent, helpful, good deed, acting innocent"},
	{Emoji: "😋", Name: "face savoring food", Label: "yum", Hint: "delicious, tasty, looking forward to food"},
	{Emoji: "😏", Name: "smirking face", Label: "smirk", Hint: "sly, knowing, smug, cheeky"},
	{Emoji: "🤔", Name: "thinking face", Label: "thinking", Hint: "considering, question, wondering, puzzling"},
	{Emoji: "🧐", Name: "face with monocle", Label: "monocle", Hint: "skeptical, inspecting, checking, curious"},
	{Emoji: "😐", Name: "neutral face", Label: "neutral", Hint: "unimpressed, awkward, speechless, flat reaction"},
	{Emoji: "😑", Name: "expressionless face", Label: "over it", Hint: "resigned, unamused, this again, exhausted patience"},
	{Emoji: "😬", Name: "grimacing face", Label: "grimace", Hint: "awkward, yikes, uncomfortable, nervous"},
	{Emoji: "😅", Name: "grinning face with sweat", Label: "nervous laugh", Hint: "close call, embarrassed, awkward relief, nervous laughter"},
	{Emoji: "😳", Name: "flushed face", Label: "flushed", Hint: "embarrassed, surprised, blushing, caught off guard"},
	{Emoji: "😮", Name: "face with open mouth", Label: "surprised", Hint: "unexpected, surprised, oh, wow"},
	{Emoji: "😲", Name: "astonished face", Label: "astonished", Hint: "unbelievable, astonished, amazing surprise"},
	{Emoji: "😱", Name: "face screaming in fear", Label: "scream", Hint: "shock, horror, alarm, dramatic surprise"},
	{Emoji: "🤯", Name: "exploding head", Label: "mind blown", Hint: "mind blown, revelation, astonishing, cannot believe it"},

	// Sympathy, frustration, rest, and feeling unwell.
	{Emoji: "😢", Name: "crying face", Label: "sad", Hint: "sad, hurt, disappointment, a tear"},
	{Emoji: "😭", Name: "loudly crying face", Label: "sobbing", Hint: "crying, devastated, overwhelmed, crying with laughter"},
	{Emoji: "🥺", Name: "pleading face", Label: "pleading", Hint: "please, hopeful, cute, touched"},
	{Emoji: "😔", Name: "pensive face", Label: "downcast", Hint: "disappointed, regretful, sad, reflective"},
	{Emoji: "😞", Name: "disappointed face", Label: "disappointed", Hint: "let down, failure, bad news, disappointment"},
	{Emoji: "😕", Name: "confused face", Label: "confused", Hint: "uncertain, do not understand, puzzled"},
	{Emoji: "😟", Name: "worried face", Label: "worried", Hint: "concern, worry, uneasy, anxious"},
	{Emoji: "😰", Name: "anxious face with sweat", Label: "anxious", Hint: "pressure, anxiety, nervous, stressed"},
	{Emoji: "😨", Name: "fearful face", Label: "fearful", Hint: "scared, alarmed, apprehensive, fear"},
	{Emoji: "😡", Name: "enraged face", Label: "angry", Hint: "furious, frustrating, angry, unfair"},
	{Emoji: "🤬", Name: "face with symbols on mouth", Label: "furious", Hint: "swearing, outrage, extremely angry"},
	{Emoji: "😴", Name: "sleeping face", Label: "sleeping", Hint: "tired, asleep, bedtime, boring meeting"},
	{Emoji: "🤤", Name: "drooling face", Label: "drooling", Hint: "delicious food, craving, wanting something"},
	{Emoji: "🤒", Name: "face with thermometer", Label: "ill", Hint: "sick, fever, unwell, taking a sick day"},
	{Emoji: "🤕", Name: "face with head bandage", Label: "hurt", Hint: "injury, headache, ouch, recovering"},
	{Emoji: "🤢", Name: "nauseated face", Label: "queasy", Hint: "gross, nauseous, revolting, feeling sick"},
	{Emoji: "🤮", Name: "face vomiting", Label: "vomiting", Hint: "disgusting, vomiting, food poisoning, extreme disgust"},
	{Emoji: "🤧", Name: "sneezing face", Label: "sneeze", Hint: "cold, allergies, sneezing, sniffles"},
	{Emoji: "🥵", Name: "hot face", Label: "overheated", Hint: "hot weather, overheating, intense effort, sweating"},
	{Emoji: "🥶", Name: "cold face", Label: "freezing", Hint: "freezing, cold weather, winter, chilly"},

	// Appreciation, celebrations, and useful signals.
	{Emoji: "❤️", Name: "red heart", Label: "love", Hint: "love, care, heartfelt appreciation, affection"},
	{Emoji: "💛", Name: "yellow heart", Label: "yellow heart", Hint: "friendship, warmth, cheerful affection"},
	{Emoji: "💜", Name: "purple heart", Label: "purple heart", Hint: "purple, affection, heartfelt support"},
	{Emoji: "💙", Name: "blue heart", Label: "blue heart", Hint: "blue, trust, friendship, calm support"},
	{Emoji: "💚", Name: "green heart", Label: "green heart", Hint: "green, nature, environmental care, affection"},
	{Emoji: "💖", Name: "sparkling heart", Label: "sparkling heart", Hint: "delight, affection, feeling special, love"},
	{Emoji: "💔", Name: "broken heart", Label: "heartbroken", Hint: "heartache, loss, disappointment, sad news"},
	{Emoji: "💕", Name: "two hearts", Label: "two hearts", Hint: "mutual love, friendship, affectionate connection"},
	{Emoji: "💯", Name: "hundred points", Label: "perfect", Hint: "exactly, fully agree, perfect, full marks"},
	{Emoji: "✅", Name: "check mark button", Label: "done", Hint: "done, completed, checked, confirmed, yes"},
	{Emoji: "❌", Name: "cross mark", Label: "nope", Hint: "wrong, rejected, canceled, no"},
	{Emoji: "⚠️", Name: "warning", Label: "warning", Hint: "caution, warning, potential problem, risk"},
	{Emoji: "❓", Name: "question mark", Label: "question", Hint: "why, clarification, confusion, a question"},
	{Emoji: "❗", Name: "exclamation mark", Label: "important", Hint: "important, emphasis, urgency, attention"},
	{Emoji: "💡", Name: "light bulb", Label: "idea", Hint: "idea, insight, solution, figured it out"},
	{Emoji: "💬", Name: "speech balloon", Label: "chat", Hint: "conversation, talking, discussion, message"},
	{Emoji: "💭", Name: "thought balloon", Label: "thought", Hint: "thinking, wondering, daydream, pondering"},
	{Emoji: "🎉", Name: "party popper", Label: "party popper", Hint: "celebration, good news, congratulations, achievement"},
	{Emoji: "🎊", Name: "confetti ball", Label: "confetti", Hint: "festive, celebration, success, congratulations"},
	{Emoji: "✨", Name: "sparkles", Label: "sparkles", Hint: "shiny, polished, magic, fresh, brilliant"},

	// Work, plans, learning, and small victories.
	{Emoji: "👀", Name: "eyes", Label: "eyes", Hint: "watching, intrigued, attention, curious, tell me more"},
	{Emoji: "👂", Name: "ear", Label: "listening", Hint: "listening, hearing, all ears, tell me"},
	{Emoji: "🚀", Name: "rocket", Label: "rocket", Hint: "launch, progress, taking off, fast growth"},
	{Emoji: "🏆", Name: "trophy", Label: "trophy", Hint: "win, achievement, success, champion"},
	{Emoji: "🥇", Name: "first place medal", Label: "gold medal", Hint: "first place, best, winning, top result"},
	{Emoji: "🎯", Name: "bullseye", Label: "on target", Hint: "exactly, accurate, on point, goal achieved"},
	{Emoji: "🏅", Name: "sports medal", Label: "medal", Hint: "reward, achievement, effort, recognition"},
	{Emoji: "📈", Name: "chart increasing", Label: "growth", Hint: "improvement, increasing, progress, trending up"},
	{Emoji: "📉", Name: "chart decreasing", Label: "decline", Hint: "decreasing, setback, drop, trending down"},
	{Emoji: "📅", Name: "calendar", Label: "calendar", Hint: "schedule, meeting, date, planning"},
	{Emoji: "⏰", Name: "alarm clock", Label: "alarm", Hint: "early, wake up, deadline, late, time"},
	{Emoji: "⌛", Name: "hourglass done", Label: "hourglass", Hint: "running out of time, waiting, deadline, patience"},
	{Emoji: "💤", Name: "zzz", Label: "zzz", Hint: "sleep, boredom, fatigue, snooze"},
	{Emoji: "💻", Name: "laptop", Label: "laptop", Hint: "work, coding, computer, remote work"},
	{Emoji: "🐛", Name: "bug", Label: "bug", Hint: "software bug, glitch, broken, debugging, caterpillar"},
	{Emoji: "🔧", Name: "wrench", Label: "fixing", Hint: "fixing, repair, maintenance, adjustment"},
	{Emoji: "🛠️", Name: "hammer and wrench", Label: "tools", Hint: "building, troubleshooting, work in progress, repair"},
	{Emoji: "📚", Name: "books", Label: "books", Hint: "reading, studying, learning, research"},
	{Emoji: "📝", Name: "memo", Label: "notes", Hint: "taking notes, writing, plan, remember"},
	{Emoji: "📎", Name: "paperclip", Label: "attachment", Hint: "attachment, document, file, paperwork"},

	// Food, weather, pets, and a little extra whimsy.
	{Emoji: "🫖", Name: "teapot", Label: "tea", Hint: "tea, brewing, warm drink, gossip"},
	{Emoji: "🥤", Name: "cup with straw", Label: "drink", Hint: "drink, soda, hydration, thirsty"},
	{Emoji: "🍪", Name: "cookie", Label: "cookie", Hint: "treat, sweet, baking, reward"},
	{Emoji: "🍰", Name: "shortcake", Label: "cake", Hint: "birthday, dessert, celebration, cake"},
	{Emoji: "🍦", Name: "soft ice cream", Label: "ice cream", Hint: "treat, hot weather, dessert, ice cream"},
	{Emoji: "🥨", Name: "pretzel", Label: "pretzel", Hint: "snack, twisted, knots, salty"},
	{Emoji: "🌮", Name: "taco", Label: "taco", Hint: "taco, lunch, dinner, hungry"},
	{Emoji: "🍿", Name: "popcorn", Label: "popcorn", Hint: "drama, watching, entertainment, movie, snack"},
	{Emoji: "🍜", Name: "steaming bowl", Label: "noodles", Hint: "ramen, noodles, soup, hungry, comfort food"},
	{Emoji: "🥗", Name: "green salad", Label: "salad", Hint: "healthy lunch, greens, salad, vegetables"},
	{Emoji: "🐶", Name: "dog face", Label: "dog", Hint: "dog, puppy, pet, excited, loyal"},
	{Emoji: "🐱", Name: "cat face", Label: "cat", Hint: "cat, kitty, pet, independent"},
	{Emoji: "🦊", Name: "fox", Label: "fox", Hint: "clever, sly, fox, cunning"},
	{Emoji: "🦉", Name: "owl", Label: "owl", Hint: "night owl, late, wise, awake at night"},
	{Emoji: "🦄", Name: "unicorn", Label: "unicorn", Hint: "rare, magical, unique, too good to be true"},
	{Emoji: "🌞", Name: "sun with face", Label: "sunshine", Hint: "sunny, bright, warm, cheerful, good morning"},
	{Emoji: "🌧️", Name: "cloud with rain", Label: "rain", Hint: "rainy, wet, dreary, gloomy weather"},
	{Emoji: "🌈", Name: "rainbow", Label: "rainbow", Hint: "hope, colorful, bright, inclusion, after the rain"},
	{Emoji: "🎈", Name: "balloon", Label: "balloon", Hint: "birthday, party, celebration, festive"},
	{Emoji: "🧯", Name: "fire extinguisher", Label: "fire extinguisher", Hint: "putting out fires, crisis, incident, calming things down"},
}

// QuestionSuffix frames the whole set: the model sees this phrasing on every
// question and the state carries the message itself.
const QuestionSuffix = "is a natural, fitting reply to this message"

// ID returns the question key for index i, with at least two digits
// ("r00", "r01", ..., "r99", "r100", ...).
func ID(i int) string {
	return fmt.Sprintf("r%02d", i)
}

// Question builds the Noul question that scores one entry as a reply.
func Question(e Entry) jev.Question {
	return jev.NoulWithCriteria(
		fmt.Sprintf("Reacting with %s (%s: %s) %s.", e.Emoji, e.Label, e.Hint, QuestionSuffix),
		"the emoji naturally fits the message's intent or emotion, acknowledges it, or lands as an apt joke",
		"the emoji feels unrelated, contradicts the tone, or only matches an incidental word",
	)
}

// Questions builds the full question map, keyed by ID in catalog order.
func Questions() map[string]jev.Question {
	questions := make(map[string]jev.Question, len(Catalog))
	for i, e := range Catalog {
		questions[ID(i)] = Question(e)
	}
	return questions
}

// State wraps the typed message into the state object the model evaluates.
func State(message string) map[string]any {
	return map[string]any{
		"message":   message,
		"context":   "one message in a casual group chat or Slack conversation",
		"reactions": "Slack-style emoji reactions a friend could tap to reply",
	}
}
