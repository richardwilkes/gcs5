package library

// Config holds the configuration information for a library of data files.
type Config struct {
	Title    string  `json:"title"`
	GitHub   string  `json:"github"`
	Repo     string  `json:"repo"`
	Path     string  `json:"path"`
	LastSeen Version `json:"last_seen"`
}
