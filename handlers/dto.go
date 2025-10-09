package handlers
type NoteRequest struct {
    Title   string `json:"title"`
    Content string `json:"content"`
}

type NoteResponse struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}