package internal

type Todo struct {
	ID          string `json:"id" bson:"_id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Status      string `json:"status" bson:"status"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
