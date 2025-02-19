package tests

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Talonmortem/http-file-server/internal/config"
	"github.com/Talonmortem/http-file-server/internal/database"
	"github.com/Talonmortem/http-file-server/internal/middleware"
	"github.com/Talonmortem/http-file-server/internal/router"

	"golang.org/x/crypto/bcrypt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexHandlerWithAuth(t *testing.T) {
	// Тестовый конфиг
	cfg := &config.Config{
		Server: config.Server{
			Host: "localhost",
			Port: 8083,
		},
		Storage: config.Storage{
			UploadDir:   t.TempDir(), // Временная директория для тестов
			PublicDir:   "public",
			TemplateDir: "../web/templates/*.html",
			WebDir:      "../web",
		},
		JWT: config.JWT{
			SecretKey: "correct-secret-key-123", // Правильный секрет
			ExpiresIn: 1,                        // 1 час
		},
		Database: config.Database{
			Path: "test.db",
		},
	}
	// Мок баз данных
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	database.DB = db

	router := router.SetupRouter(cfg)

	t.Run("Валидация токена", func(t *testing.T) {
		// Генерируем токен
		token, err := middleware.GenerateToken("testuser", cfg)
		require.NoError(t, err)

		// Проверяем, что токен валиден
		claims, err := middleware.ValidateToken(token, cfg.JWT.SecretKey)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", claims.Username)
	})

	t.Run("Отображение формы", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/login", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "<form method=\"POST\"")
	})

	t.Run("Успешная авторизация", func(t *testing.T) {
		// Мокируем базу данных
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
		rows := sqlmock.NewRows([]string{"password_hash"}).AddRow(string(hashedPass))
		mock.ExpectQuery("SELECT ...").WillReturnRows(rows)

		// Формируем запрос
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "testpass")

		req, _ := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Проверяем редирект и куки
		assert.Equal(t, http.StatusFound, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "/")
		cookies := w.Result().Cookies()
		assert.True(t, len(cookies) > 0)
	})

	t.Run("Доступ к защищенному ресурсу с валидным токеном", func(t *testing.T) {
		// Генерируем тестовый токен
		token, err := middleware.GenerateToken("testuser", cfg)
		require.NoError(t, err)

		// Запрос к защищенному ресурсу
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Проверка с невалидным токеном", func(t *testing.T) {
		// Конфиг для генерации НЕвалидного токена
		invalidCfg := &config.Config{
			Server: config.Server{
				Host: "localhost",
				Port: 8083,
			},
			Storage: config.Storage{
				UploadDir: t.TempDir(),
				PublicDir: "public",
			},
			JWT: config.JWT{
				SecretKey: "wrong-secret-key-456", // Неправильный секрет
				ExpiresIn: 1,
			},
			Database: config.Database{
				Path: "test.db",
			},
		}

		token, _ := middleware.GenerateToken("testuser", invalidCfg)
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Доступ к скачиванию файлов с валидным токеном", func(t *testing.T) {
		// Генерируем тестовый токен
		token, _ := middleware.GenerateToken("testuser", cfg)
		request := struct {
			Files []string `json:"files"`
		}{
			Files: []string{"file1.txt", "file2.txt"},
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		// Запрос к защищенному ресурсу
		req, _ := http.NewRequest("POST", "/download-zip", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Доступ к скачиванию файлов с невалидным токеном", func(t *testing.T) {

		// Генерируем тестовый токен
		token := "wrong-token"
		request := struct {
			Files []string `json:"files"`
		}{
			Files: []string{"file1.txt", "file2.txt"},
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		// Запрос к защищенному ресурсу
		req, _ := http.NewRequest("POST", "/download-zip", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Проверка загрузки с токеном", func(t *testing.T) {
		// Генерируем тестовый токен
		token, err := middleware.GenerateToken("testuser", cfg)
		require.NoError(t, err)

		// Создаем запрос с токеном
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.txt")
		part.Write([]byte("content"))
		writer.Close()

		req, _ := http.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Проверка загрузки с невалидным токеном", func(t *testing.T) {
		// Генерируем тестовый токен
		token := "wrong-token"
		require.NoError(t, err)

		// Создаем запрос с токеном
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.txt")
		part.Write([]byte("content"))
		writer.Close()

		req, _ := http.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("Проверка редиректа", func(t *testing.T) {

		testCases := []struct {
			name         string
			path         string
			method       string
			expectedCode int
			redirectPath string
		}{
			{
				name:         "Главная без авторизации",
				path:         "/",
				method:       "GET",
				expectedCode: http.StatusFound,
				redirectPath: "/login",
			},
			{
				name:         "Загрузка файлов без авторизации",
				path:         "/upload",
				method:       "POST",
				expectedCode: http.StatusFound,
				redirectPath: "/login",
			},
			{
				name:         "Список файлов без авторизации",
				path:         "/files",
				method:       "GET",
				expectedCode: http.StatusFound,
				redirectPath: "/login",
			},
			{
				name:         "Удаление файлов без авторизации",
				path:         "/delete",
				method:       "POST",
				expectedCode: http.StatusFound,
				redirectPath: "/login",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req, _ := http.NewRequest(tc.method, tc.path, nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedCode, w.Code)
				if tc.redirectPath != "" {
					location, _ := w.Result().Location()
					assert.Equal(t, tc.redirectPath, location.Path)
				}
			})
		}
	})
}
