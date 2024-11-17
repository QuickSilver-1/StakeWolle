package handlers

// User - структура данных о пользователе
type User struct {
    Email string `json:"email"`     // Адрес электронной почты пользователя
    Pass  string `json:"password"`  // Пароль пользователя
    Ref   string `json:"ref"`       // Реферальный код пользователя, если есть
}
