package app

type Config struct {
	GithubToken GithubToken
}

func ProvideGithubToken(config Config) GithubToken {
	return config.GithubToken
}
