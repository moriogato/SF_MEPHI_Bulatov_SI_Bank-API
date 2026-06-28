package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/beevik/etree"
)

type CBRService struct {
	client *http.Client
}

func NewCBRService() *CBRService {
	return &CBRService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// buildSOAPRequest формирует SOAP-запрос для получения ключевой ставки
func (s *CBRService) buildSOAPRequest() string {
	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
        <soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
            <soap12:Body>
                <KeyRate xmlns="http://web.cbr.ru/">
                    <fromDate>%s</fromDate>
                    <ToDate>%s</ToDate>
                </KeyRate>
            </soap12:Body>
        </soap12:Envelope>`, fromDate, toDate)
}

// sendRequest отправляет SOAP-запрос и возвращает ответ
func (s *CBRService) sendRequest(soapRequest string) ([]byte, error) {
	req, err := http.NewRequest(
		"POST",
		"https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx",
		bytes.NewBuffer([]byte(soapRequest)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://web.cbr.ru/KeyRate")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return rawBody, nil
}

// parseXMLResponse парсит XML-ответ и извлекает ключевую ставку
func (s *CBRService) parseXMLResponse(rawBody []byte) (float64, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(rawBody); err != nil {
		return 0, fmt.Errorf("failed to parse XML: %w", err)
	}

	// Ищем элементы в ответе
	krElements := doc.FindElements("//diffgram/KeyRate/KR")
	if len(krElements) == 0 {
		return 0, errors.New("no key rate data found")
	}

	latestKR := krElements[0]
	rateElement := latestKR.FindElement("./Rate")
	if rateElement == nil {
		return 0, errors.New("Rate element not found")
	}

	rateStr := rateElement.Text()
	var rate float64
	if _, err := fmt.Sscanf(rateStr, "%f", &rate); err != nil {
		return 0, fmt.Errorf("failed to parse rate: %w", err)
	}

	return rate, nil
}

// GetKeyRate получает ключевую ставку ЦБ РФ и добавляет маржу банка
func (s *CBRService) GetKeyRate() (float64, error) {
	soapRequest := s.buildSOAPRequest()
	rawBody, err := s.sendRequest(soapRequest)
	if err != nil {
		return 0, err
	}

	rate, err := s.parseXMLResponse(rawBody)
	if err != nil {
		return 0, err
	}

	// Добавляем маржу банка (+5%)
	rate += 5.0
	return rate, nil
}
