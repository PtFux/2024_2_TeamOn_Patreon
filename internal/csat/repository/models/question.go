package models

// Repository модель Question
type Question struct {
	// Вопрос
	Question string
	// ИД вопроса
	QuestionID string
	// ИД темы вопроса
	QuestionThemeID string
}
