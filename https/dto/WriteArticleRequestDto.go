package dto

type WriteArticleRequestDto struct {
	Category    string            `json:"category"`
	SubCategory string            `json:"subcategory"`
	ID          string            `json:"id"`
	Site        string            `json:"website"`
	Meta        map[string]string `json:"meta"`
	Content     JsonPayload       `json:"content"`
	Resume      string            `json:"resume"`
}

type JsonPayload struct {
	Value string
}

func (b *JsonPayload) UnmarshalJSON(data []byte) error {
	b.Value = string(data)
	return nil
}
