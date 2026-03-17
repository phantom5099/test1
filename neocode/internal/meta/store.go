package meta

// Store：未来用于元数据存储的简单实现，占位以便扩展
type Store struct {
	History History
}

func NewStore() *Store { return &Store{} }
