package site

// Profile holds the non-localized facts about the site owner.
// Localized strings (bio, tagline, nav labels) live in locales/*.json.
type Profile struct {
	Timezone    string
	Links       []Link
	Socials     []Link
	DefaultLang string
	Languages   []Language
}

type Link struct {
	Label string
	URL   string
	Rel   string
}

type Language struct {
	Code  string
	Path  string
	Label string
}

// Edit this to make the site yours.
var Me = Profile{
	Timezone:    "Asia/Almaty",
	DefaultLang: "en",
	Languages: []Language{
		{Code: "en", Path: "/", Label: "English"},
	},
	Links: []Link{
		{Label: "cv", URL: "https://cv.magzhan.me", Rel: "me"},
	},
	Socials: []Link{
		{Label: "GitHub", URL: "https://github.com/mzhakuda", Rel: "me"},
		{Label: "LinkedIn", URL: "https://www.linkedin.com/in/mzhakuda", Rel: "me"},
		{Label: "X", URL: "https://x.com/mzhakuda", Rel: "me"},
		{Label: "Telegram", URL: "https://t.me/mzhakuda", Rel: "me"},
		{Label: "Email", URL: "mailto:mzhakuda@gmail.com", Rel: "me"},
	},
}
