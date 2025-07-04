package Model

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	FileName    string  `json:"file_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CreatedAt   string  `json:"created_at"`
}
