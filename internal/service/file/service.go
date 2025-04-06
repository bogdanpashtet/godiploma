package file

import (
	"context"
	"fmt"
	"time"

	filev1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/file/v1"
	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
	"go.uber.org/zap"
)

const kind = "FileService"

type Service struct {
	l *zap.Logger
}

func New(
	l *zap.Logger,
) *Service {
	return &Service{
		l: l,
	}
}

//func (s *Service) UploadFiles(ctx context.Context, req *filed.UploadFilesRequest) (*filev1.UploadDocumentsResponse, error) {
//	processedCount := 0
//	var processingError error
//
//	for i, doc := range req.GetDocuments() {
//		docType := doc.GetDocumentType()
//		docData := doc.GetDocumentData()
//
//		docLogger := s.l.With(
//			zap.Any("documentIndex", i),
//			zap.Any("documentType", docType.String()),
//			zap.Any("documentSize", len(docData)),
//		)
//
//		// 3. Вызов функции симуляции обработки для текущего документа
//		err := s.simulateDocumentProcessing(ctx, docType, docData, docLogger)
//		if err != nil {
//			// Обработка ошибки: Записываем первую ошибку и прерываем цикл
//			docLogger.Errorw("Failed to process document", zap.Error(err))
//			processingError = fmt.Errorf("failed processing document at index %d (%s): %w", i, docType.String(), err)
//			s.metrics.ProcessingErrorsTotal.Inc() // Увеличиваем счетчик ошибок (если есть)
//			break                                 // Прекращаем обработку пакета при первой ошибке
//		}
//
//		// Успешная обработка документа
//		docLogger.Info("Document processed successfully")
//		processedCount++
//		s.metrics.ProcessedFilesTotal.Inc() // Увеличиваем счетчик успеха (если есть)
//		// Можно добавить метрику по типу документа: s.metrics.ProcessedDocumentsByType.WithLabelValues(docType.String()).Inc()
//	}
//
//	// 4. Формирование ответа
//	if processingError != nil {
//		// Если была ошибка при обработке хотя бы одного документа
//		l.Errorw("Request processing failed", zap.Error(processingError))
//		return &filehandlerv1.ProcessDocumentsResponse{
//			Success:            false,
//			Message:            processingError.Error(),
//			DocumentsProcessed: int64(processedCount), // Сколько успели обработать до ошибки
//		}, nil // Возвращаем nil, т.к. это не ошибка выполнения gRPC метода, а ошибка бизнес-логики
//		// Или можно вернуть статус gRPC: return nil, status.Error(codes.Internal, processingError.Error())
//	}
//
//	// Если все документы обработаны успешно
//	successMessage := fmt.Sprintf("Successfully processed %d documents", processedCount)
//	l.Info("Request processing finished successfully", zap.Int("processedCount", processedCount))
//}
//
//// simulateDocumentProcessing имитирует обработку одного документа
//// Принимает тип документа, данные и логгер (уже с контекстом запроса)
//func (s *Server) simulateDocumentProcessing(ctx context.Context, docType filed.DocumentType, docData []byte, l *zap.Logger) error {
//	startTime := time.Now()
//	dataSize := len(docData)
//	l.Debug("Starting document processing simulation") // Используем переданный логгер
//
//	// --- Варьируем нагрузку в зависимости от типа документа ---
//	var cpuIterations int
//	var memAllocMultiplier float64 = 0.1 // По умолчанию выделяем 10% от размера
//	var baseDelay time.Duration = 10 * time.Millisecond
//
//	switch docType {
//	case filev1.DocumentType_DOCUMENT_TYPE_PASSPORT:
//		cpuIterations = 15000 + dataSize/256
//		memAllocMultiplier = 0.25
//		baseDelay = 25 * time.Millisecond
//		l.Debug("Using PASSPORT simulation profile")
//	case filev1.DocumentType_DOCUMENT_TYPE_BANK_STATEMENT:
//		cpuIterations = 8000 + dataSize/512
//		memAllocMultiplier = 0.05
//		baseDelay = 15 * time.Millisecond
//		l.Debug("Using BANK_STATEMENT simulation profile")
//	default: // DOCUMENT_TYPE_DRIVER_LICENSE и другие
//		cpuIterations = 10000 + dataSize/512
//		baseDelay = 10 * time.Millisecond
//		l.Debug("Using default simulation profile")
//	}
//
//	// --- A. Симуляция CPU нагрузки ---
//	var cpuResult uint64
//	for i := 0; i < cpuIterations; i++ {
//		if dataSize > 0 {
//			cpuResult = (cpuResult << 3) + uint64(docData[i%dataSize]) + uint64(i)
//		}
//		if i%10000 == 0 { // Периодическая проверка контекста
//			select {
//			case <-ctx.Done():
//				l.Warnw("Processing cancelled by context during CPU simulation")
//				return ctx.Err()
//			default:
//			}
//		}
//	}
//	l.Debugw("CPU simulation finished", "pseudo_result", cpuResult)
//
//	// --- Б. Симуляция использования памяти ---
//	tempBufferSize := int(float64(dataSize) * memAllocMultiplier)
//	if tempBufferSize > 512 { // Выделяем только если размер значимый
//		l.Debugw("Allocating temporary buffer", "size", tempBufferSize)
//		tempBuffer := make([]byte, tempBufferSize)
//		copy(tempBuffer, docData[:min(tempBufferSize, dataSize)]) // Копируем часть данных
//		// Имитируем работу с буфером (удерживаем в памяти)
//		time.Sleep(time.Duration(5+rand.Intn(25)) * time.Millisecond)
//		_ = tempBuffer // Используем, чтобы компилятор не выкинул
//	}
//
//	// --- В. Симуляция ожидания (I/O и т.п.) ---
//	randomDelay := time.Duration(rand.Intn(35)) * time.Millisecond
//	totalDelay := baseDelay + randomDelay
//	l.Debugw("Simulating I/O wait", "delay", totalDelay)
//	time.Sleep(totalDelay)
//
//	processingDuration := time.Since(startTime)
//	l.Info("Document processing simulation finished", "duration", processingDuration)
//
//	// // Опционально: Имитация случайной ошибки
//	// if rand.Intn(100) < 3 { // 3% шанс на ошибку
//	//     return fmt.Errorf("simulated random processing error")
//	// }
//
//	return nil // Успех
//}
//
//// Вспомогательная функция min для Go < 1.21
//func min(a, b int) int {
//	if a < b {
//		return a
//	}
//	return b
//}
