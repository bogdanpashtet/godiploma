package v1

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	filev1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/file/v1"
)

const kind = "FileServer"

type Params struct {
	fx.In

	Logger *zap.Logger
}
type Server struct {
	l *zap.Logger

	filev1.UnimplementedFileServiceServer
}

func NewServer(params Params) *Server {
	return &Server{
		l: params.Logger,
	}
}

func (s *Server) Register(gRPCServer *grpc.Server) {
	filev1.RegisterFileServiceServer(gRPCServer, s)
}

// UploadDocuments теперь использует парсинг "сложного" файла
func (s *Server) UploadDocuments(ctx context.Context, req *filev1.UploadDocumentsRequest) (*filev1.UploadDocumentsResponse, error) {
	const handlerName = "UploadDocuments"
	s.l = s.l.With(zap.String("request_id", "test req"))

	s.l.Info(fmt.Sprintf("%s: received request (parsing simulation)", handlerName))
	//timer := prometheus.NewTimer(s.metrics.ProcessingDuration.WithLabelValues(handlerName))
	//defer timer.ObserveDuration()
	//
	//if req == nil || len(req.GetFileContent()) == 0 {
	//	s.l.Warn(fmt.Sprintf("%s: request is nil or file content is empty", handlerName), zap.String("request_id", req.GetRequestId()))
	//	s.metrics.ProcessingErrorsTotal.Inc()
	//	return &filev1.UploadDocumentsResponse{Success: false, Message: "No file content"}, nil
	//}
	//fileData := req.GetFileContent()
	//fileSize := len(fileData)
	//s.l.Debug(fmt.Sprintf("%s: starting complex parse for %d bytes", handlerName, fileSize), zap.String("request_id", req.GetRequestId()))
	//
	//// --- 3. Вызов функции парсинга ---
	//parsedEvents, err := s.parseSessionLogFile(ctx, fileData)
	//if err != nil {
	//	s.l.Error(fmt.Sprintf("%s: session log parsing failed", handlerName), zap.Error(err), zap.String("request_id", req.GetRequestId()))
	//	s.metrics.ProcessingErrorsTotal.Inc()
	//	return &filev1.UploadDocumentsResponse{Success: false, Message: fmt.Sprintf("Parsing failed: %v", err)}, nil
	//}
	//
	//// --- 4. Успешное завершение ---
	//s.metrics.ProcessedFilesTotal.Inc()
	//eventCount := len(parsedEvents)
	//s.l.Info(fmt.Sprintf("%s: successfully parsed %d bytes into %d events", handlerName, fileSize, eventCount), zap.String("request_id", req.GetRequestId()))
	//
	//// --- 5. Формируем ответ ---
	//response := &filev1.UploadDocumentsResponse{
	//	Success: true,
	//	Message: fmt.Sprintf("Parsed %d bytes into %d events successfully", fileSize, eventCount),
	//	// Можно добавить в ответ количество событий или другую статистику
	//}

	return nil, nil
}

//// parseSessionLogFile имитирует парсинг сложного лог-файла
//func (s *Server) parseSessionLogFile(ctx context.Context, fileData []byte) ([]SessionEvent, error) {
//	startTime := time.Now()
//	// Предварительно выделяем память под слайс, чтобы уменьшить количество реаллокаций
//	// Эвристика: средняя длина строки ~150 байт
//	estimatedLines := len(fileData) / 150
//	if estimatedLines == 0 {
//		estimatedLines = 100
//	}
//	events := make([]SessionEvent, 0, estimatedLines)
//
//	scanner := bufio.NewScanner(bytes.NewReader(fileData))
//	lineNumber := 0
//
//	for scanner.Scan() {
//		lineNumber++
//		line := scanner.Text()
//
//		if line == "" { // Пропускаем пустые строки
//			continue
//		}
//
//		// --- A. Парсинг строки ---
//		parts := strings.SplitN(line, "|", 5) // Ожидаем 5 частей
//		if len(parts) != 5 {
//			s.l.Warn("Invalid line format", zap.Int("line", lineNumber), zap.String("content", line))
//			continue // Пропускаем некорректную строку
//		}
//
//		event := SessionEvent{LineNumber: lineNumber}
//		var err error
//
//		// 1. Timestamp
//		tsNano, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
//		if err != nil {
//			s.l.Warn("Invalid timestamp", zap.Int("line", lineNumber), zap.String("part", parts[0]), zap.Error(err))
//			continue
//		}
//		event.Timestamp = time.Unix(0, tsNano)
//
//		// 2. Session ID (просто строка)
//		event.SessionID = strings.TrimSpace(parts[1])
//		// Можно добавить валидацию UUID, если нужно: uuid.Parse(event.SessionID)
//
//		// 3. User ID
//		event.UserID, err = strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 64)
//		if err != nil {
//			s.l.Warn("Invalid UserID", zap.Int("line", lineNumber), zap.String("part", parts[2]), zap.Error(err))
//			continue
//		}
//
//		// 4. Event Type
//		event.EventType = strings.TrimSpace(parts[3])
//
//		// 5. JSON Payload (Самая ресурсоемкая часть)
//		payloadStr := strings.TrimSpace(parts[4])
//		event.RawPayload = json.RawMessage(payloadStr) // Сохраняем сырой JSON
//		// Пытаемся распарсить JSON - это создает нагрузку на CPU и временные аллокации
//		if err := json.Unmarshal(event.RawPayload, &event.Payload); err != nil {
//			s.l.Warn("Invalid JSON payload", zap.Int("line", lineNumber), zap.String("payload", payloadStr), zap.Error(err))
//			// Можно не считать это фатальной ошибкой строки, если остальное нужно
//			event.Payload = nil // Очищаем, если не распарсился
//			// continue // Или пропускаем всю строку
//		}
//
//		// --- Б. Накопление результата ---
//		// Добавляем распарсенное событие в общий слайс. Слайс растет, память используется.
//		events = append(events, event)
//
//		// --- В. Проверка контекста (важно для больших файлов) ---
//		if lineNumber%5000 == 0 { // Проверяем каждые 5000 строк
//			select {
//			case <-ctx.Done():
//				s.l.Warn("Parsing cancelled by context", zap.Int("lines_parsed", lineNumber))
//				return nil, ctx.Err()
//			default: // продолжаем
//			}
//			s.l.Debug("Parsing progress", zap.Int("lines_parsed", lineNumber), zap.Int("events_collected", len(events)))
//		}
//	} // Конец цикла по строкам
//
//	if err := scanner.Err(); err != nil {
//		s.l.Error("Error reading file data", zap.Error(err))
//		return nil, fmt.Errorf("scanner error: %w", err)
//	}
//
//	// Имитируем небольшую задержку, как будто мы что-то делаем с накопленными 'events'
//	time.Sleep(50 * time.Millisecond)
//
//	duration := time.Since(startTime)
//	s.l.Info("Finished parsing session log file",
//		zap.Duration("duration", duration),
//		zap.Int("lines_total", lineNumber),
//		zap.Int("events_parsed", len(events)),
//	)
//	return events, nil
//}
