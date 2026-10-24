package site

// Profile holds the non-localized facts about the site owner.
// Localized strings (title, tagline, CV labels) live in locales/*.json.
type Profile struct {
	Links       []Link
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
	DefaultLang: "en",
	Languages: []Language{
		{Code: "en", Path: "/", Label: "English"},
	},
	Links: []Link{
		{Label: "CV", URL: "https://cv.magzhan.me", Rel: "me"},
		{Label: "GitHub", URL: "https://github.com/mzhakuda", Rel: "me"},
		{Label: "LinkedIn", URL: "https://www.linkedin.com/in/mzhakuda", Rel: "me"},
		{Label: "Email", URL: "mailto:mzhakuda@gmail.com", Rel: "me"},
	},
}
