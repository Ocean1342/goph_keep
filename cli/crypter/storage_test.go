package crypter

import (
	"errors"
	"sync"
	"testing"
)

// TestStorage_NewStorage проверяет создание нового хранилища
func TestStorage_NewStorage(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	if s.storage == nil {
		t.Error("storage map should be initialized")
	}
}

// TestStorage_Set проверяет добавление ключа
func TestStorage_Set(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	login := "testuser"
	key := []byte("secret_key_123")

	s.Set(login, key)

	s.mu.Lock()
	storedKey, exists := s.storage[login]
	s.mu.Unlock()

	if !exists {
		t.Error("key should be stored in storage")
	}

	if string(storedKey) != string(key) {
		t.Errorf("stored key mismatch: got %s, want %s", storedKey, key)
	}
}

// TestStorage_SetOverwrite проверяет перезапись существующего ключа
func TestStorage_SetOverwrite(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	login := "testuser"
	key1 := []byte("old_key")
	key2 := []byte("new_key")

	s.Set(login, key1)
	s.Set(login, key2)

	s.mu.Lock()
	storedKey, exists := s.storage[login]
	s.mu.Unlock()

	if !exists {
		t.Error("key should exist after overwrite")
	}

	if string(storedKey) != string(key2) {
		t.Errorf("key should be overwritten: got %s, want %s", storedKey, key2)
	}
}

// TestStorage_GetExisting проверяет получение существующего ключа
func TestStorage_GetExisting(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	login := "testuser"
	expectedKey := []byte("test_key")

	s.Set(login, expectedKey)

	key, err := s.Get(login)
	if err != nil {
		t.Errorf("Get should not return error for existing key: %v", err)
	}

	if string(key) != string(expectedKey) {
		t.Errorf("Get returned wrong key: got %s, want %s", key, expectedKey)
	}
}

// TestStorage_GetNonExisting проверяет получение несуществующего ключа
func TestStorage_GetNonExisting(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	login := "nonexistent"

	key, err := s.Get(login)
	if err == nil {
		t.Error("Get should return error for non-existing key")
	}

	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("Get should return ErrKeyNotFound, got: %v", err)
	}

	if len(key) != 0 {
		t.Errorf("Get should return empty slice for non-existing key, got: %v", key)
	}
}

// TestStorage_GetEmptyLogin проверяет получение ключа с пустым логином
func TestStorage_GetEmptyLogin(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	key, err := s.Get("")
	if err == nil {
		t.Error("Get should return error for empty login")
	}

	if !errors.Is(err, ErrEmptyLogin) {
		t.Errorf("Get should return ErrKeyNotFound for empty login, got: %v", err)
	}

	if len(key) != 0 {
		t.Errorf("Get should return empty slice for empty login, got: %v", key)
	}
}

// TestStorage_ConcurrentSet проверяет конкурентную запись
func TestStorage_ConcurrentSet(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	var wg sync.WaitGroup
	users := 100

	// Конкурентно записываем 100 пользователей
	for i := 0; i < users; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			login := string(rune('a' + (id % 26)))
			key := []byte{byte(id)}
			s.Set(login, key)
		}(i)
	}

	wg.Wait()

	// Проверяем, что все записи завершились без паники
	s.mu.Lock()
	count := len(s.storage)
	s.mu.Unlock()

	// Не можем точно знать сколько уникальных логинов было создано
	// но можем проверить что нет паники при конкурентном доступе
	t.Logf("Created %d unique logins concurrently", count)
}

// TestStorage_ConcurrentGetSet проверяет конкурентное чтение и запись
func TestStorage_ConcurrentGetSet(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	var wg sync.WaitGroup
	operations := 1000

	// Предварительно добавляем один ключ
	s.Set("shared", []byte("initial_value"))

	// Конкурентные операции чтения и записи
	for i := 0; i < operations; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			if id%2 == 0 {
				// Операция записи
				s.Set("shared", []byte{byte(id)})
			} else {
				// Операция чтения
				s.Get("shared")
			}
		}(i)
	}

	wg.Wait()

	// Проверяем, что последнее значение записано
	value, err := s.Get("shared")
	if err != nil {
		t.Errorf("Should be able to get shared key: %v", err)
	}

	t.Logf("Final value after concurrent operations: %v", value)
}

// TestStorage_ThreadSafety проверяет потокобезопасность
func TestStorage_ThreadSafety(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	// Тест на гонку данных (запускать с флагом -race)
	var wg sync.WaitGroup

	// Писатели
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				login := string(rune('a' + (j % 26)))
				key := []byte{byte(id), byte(j)}
				s.Set(login, key)
			}
		}(i)
	}

	// Читатели
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				login := string(rune('a' + (j % 26)))
				s.Get(login)
			}
		}(i)
	}

	wg.Wait()
	// Если тест проходит без паники при -race, значит потокобезопасность обеспечена
}

// TestStorage_MultipleKeys проверяет работу с несколькими ключами
func TestStorage_MultipleKeys(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	// Добавляем несколько ключей
	users := map[string][]byte{
		"alice":   []byte("alice_key"),
		"bob":     []byte("bob_key"),
		"charlie": []byte("charlie_key"),
	}

	for login, key := range users {
		s.Set(login, key)
	}

	// Проверяем все ключи
	for login, expectedKey := range users {
		key, err := s.Get(login)
		if err != nil {
			t.Errorf("Failed to get key for %s: %v", login, err)
			continue
		}

		if string(key) != string(expectedKey) {
			t.Errorf("Wrong key for %s: got %s, want %s", login, key, expectedKey)
		}
	}

	// Проверяем общее количество
	s.mu.Lock()
	count := len(s.storage)
	s.mu.Unlock()

	if count != len(users) {
		t.Errorf("Wrong number of keys in storage: got %d, want %d", count, len(users))
	}
}

// TestStorage_EmptyKey проверяет сохранение пустого ключа
func TestStorage_EmptyKey(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	login := "testuser"
	emptyKey := []byte{}

	s.Set(login, emptyKey)

	key, err := s.Get(login)
	if err != nil {
		t.Errorf("Should be able to get empty key: %v", err)
	}

	if len(key) != 0 {
		t.Errorf("Should return empty key, got: %v", key)
	}
}

// TestStorage_NilKey проверяет сохранение nil ключа
func TestStorage_NilKey(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	login := "testuser"

	s.Set(login, nil)

	key, err := s.Get(login)
	if err != nil {
		t.Errorf("Should be able to get nil key: %v", err)
	}

	if key != nil {
		t.Errorf("Should return nil, got: %v", key)
	}
}

// TestStorage_KeyModification проверяет, что ключ защищен от внешних изменений
func TestStorage_KeyModification(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	login := "testuser"
	originalKey := []byte("original")

	// Сохраняем ключ
	s.Set(login, originalKey)

	// Модифицируем оригинальный слайс
	originalKey[0] = 'X'

	// Получаем ключ обратно
	retrievedKey, err := s.Get(login)
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	// Проверяем, что полученный ключ не изменился
	if retrievedKey[0] == 'X' {
		t.Error("Storage should protect against external modification")
	}

	if retrievedKey[0] != 'o' {
		t.Errorf("Key was modified externally: got %c, want 'o'", retrievedKey[0])
	}
}

// TestStorage_ParallelReadWrite проверяет параллельное чтение/запись разных ключей
func TestStorage_ParallelReadWrite(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	var wg sync.WaitGroup
	iterations := 100

	for i := 0; i < iterations; i++ {
		wg.Add(2)

		// Писатель
		go func(id int) {
			defer wg.Done()
			login := string(rune('a' + (id % 26)))
			s.Set(login, []byte{byte(id)})
		}(i)

		// Читатель
		go func(id int) {
			defer wg.Done()
			login := string(rune('z' - (id % 26)))
			s.Get(login)
		}(i)
	}

	wg.Wait()
	// Тест пройден, если нет паники или deadlock
}

// TestStorage_GetAfterDelete проверяет поведение после удаления (косвенный тест)
func TestStorage_GetAfterDelete(t *testing.T) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	login := "testuser"
	key := []byte("test_key")

	// Добавляем ключ
	s.Set(login, key)

	// "Удаляем" ключ путем перезаписи
	s.Set(login, []byte{})

	// Получаем ключ
	retrievedKey, err := s.Get(login)
	if err != nil {
		t.Errorf("Should find empty key: %v", err)
	}

	if len(retrievedKey) != 0 {
		t.Errorf("Should get empty key after 'deletion', got: %v", retrievedKey)
	}
}

// BenchmarkStorage_Set бенчмарк для операции Set
func BenchmarkStorage_Set(b *testing.B) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	login := "testuser"
	key := []byte("benchmark_key")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set(login, key)
	}
}

// BenchmarkStorage_Get бенчмарк для операции Get
func BenchmarkStorage_Get(b *testing.B) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	login := "testuser"
	key := []byte("benchmark_key")
	s.Set(login, key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get(login)
	}
}

// BenchmarkStorage_GetParallel бенчмарк для параллельного Get
func BenchmarkStorage_GetParallel(b *testing.B) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	// Подготавливаем данные
	users := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		users[i] = string(rune('a' + (i % 26)))
		s.Set(users[i], []byte{byte(i)})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			s.Get(users[i%len(users)])
			i++
		}
	})
}

// BenchmarkStorage_SetParallel бенчмарк для параллельного Set
func BenchmarkStorage_SetParallel(b *testing.B) {
	s := &storage{
		storage: make(map[string][]byte),
		mu:      sync.Mutex{},
	}

	counter := 0
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter++
			login := string(rune('a' + (counter % 26)))
			s.Set(login, []byte{byte(counter)})
		}
	})
}
