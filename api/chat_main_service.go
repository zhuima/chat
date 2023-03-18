package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	uuid "github.com/iris-contrib/go.uuid"
	"github.com/samber/lo"
	"github.com/swuecho/chatgpt_backend/ai"
	"github.com/swuecho/chatgpt_backend/sqlc_queries"
)

type ChatService struct {
	q *sqlc_queries.Queries
}

// NewChatSessionService creates a new ChatSessionService.
func NewChatService(q *sqlc_queries.Queries) *ChatService {
	return &ChatService{q: q}
}

func (s *ChatService) Chat(chatSessionUuid string, chatUuid, newQuestion string, userID int32) (*sqlc_queries.ChatMessage, error) {
	// no session exists
	//
	// if no session chat_created, create new chat_session with $uuid
	// create a new prompt with topic = $uuid, role = "system", content= req.Prompt

	// if session avaiable,
	// GetChatPromptBySessionID and create Message from Prompt
	// GetLatestMessagesBySessionID  and create Messsage(s) from messages

	// Check if the chat session exists
	ctx := context.Background()

	// no session exists
	// create session and prompt

	chatSession, err := s.q.CreateOrUpdateChatSessionByUUID(ctx, sqlc_queries.CreateOrUpdateChatSessionByUUIDParams{
		Uuid:   chatSessionUuid,
		UserID: userID,
		Topic:  firstN(newQuestion, 30),
	})

	if err != nil {
		return nil, fmt.Errorf("fail to create or update session: %w", err)
	}

	log.Println(chatSession)

	existingPrompt := true

	log.Println(chatSessionUuid)
	_, err = s.q.GetOneChatPromptBySessionUUID(ctx, chatSessionUuid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			existingPrompt = false
		} else {
			return nil, fmt.Errorf("error when get prompt: %w", err)
		}
	}

	if existingPrompt {
		user := ai.User
		_, err := s.q.CreateChatMessage(ctx,
			sqlc_queries.CreateChatMessageParams{
				ChatSessionUuid: chatSession.Uuid,
				Uuid:            chatUuid,
				Role:            user.String(),
				Content:         newQuestion,
				Raw:             json.RawMessage([]byte("{}")),
				UserID:          userID,
				CreatedBy:       userID,
				UpdatedBy:       userID,
			})

		if err != nil {
			return nil, fmt.Errorf("add user message when not new session: %w", err)
		}
	} else {
		uuidVar, _ := uuid.NewV4()
		chatPrompt, err := s.q.CreateChatPrompt(ctx,
			sqlc_queries.CreateChatPromptParams{
				Uuid:            uuidVar.String(),
				ChatSessionUuid: chatSessionUuid,
				Role:            "system",
				Content:         newQuestion,
				UserID:          userID,
				CreatedBy:       userID,
				UpdatedBy:       userID,
			})
		if err != nil {
			return nil, fmt.Errorf("fail to create prompt: %w", err)
		}
		log.Printf("%+v\n", chatPrompt)
	}

	chat_prompts, err := s.q.GetChatPromptsBySessionUUID(ctx, chatSessionUuid)

	if err != nil {
		return nil, fmt.Errorf("fail to get prompt: %w", err)
	}

	chat_massages, err := s.q.GetLatestMessagesBySessionUUID(ctx,
		sqlc_queries.GetLatestMessagesBySessionUUIDParams{ChatSessionUuid: chatSession.Uuid, Limit: 5})

	if err != nil {
		return nil, fmt.Errorf("fail to get latest message: %w", err)
	}
	chat_prompt_msgs := lo.Map(chat_prompts, func(m sqlc_queries.ChatPrompt, _ int) Message {
		return Message{Role: m.Role, Content: m.Content}
	})
	chat_message_msgs := lo.Map(chat_massages, func(m sqlc_queries.ChatMessage, _ int) Message {
		return Message{Role: m.Role, Content: m.Content}
	})
	msgs := append(chat_prompt_msgs, chat_message_msgs...)

	if existingPrompt {
		msgs = append(msgs, NewUserMessage(newQuestion))
	}
	if len(msgs) == 0 {
		return nil, fmt.Errorf("fail to collect messages: %w", err)
	}
	var aiAnswer ChatCompletionResponse
	// if system message is test_demo_bestqa, return a demo message
	if msgs[0].Content == "test_demo_bestqa" || msgs[len(msgs)-1].Content == "test_demo_bestqa" {
		message := Message{Role: "assitant", Content: "Hi, I am a chatbot. I can help you to find the best answer for your question. Please ask me a question."}
		aiAnswer.Choices = []Choice{{Message: message}}
	} else {
		aiAnswer, err = GetAiAnswerOpenApi(msgs)
		if err != nil {
			return nil, fmt.Errorf("error when try get answer from service %w", err)
		}
	}

	// Generate a response message based on the input prompt
	answer := aiAnswer.Choices[0].Message

	jsonMsg, err := json.Marshal(aiAnswer)
	if err != nil {
		return nil, fmt.Errorf("error when try to serialize answer %w", err)
	}

	answerUuid, _ := uuid.NewV4()
	answer_msg, err := s.q.CreateChatMessage(ctx,
		sqlc_queries.CreateChatMessageParams{
			ChatSessionUuid: chatSession.Uuid,
			Uuid:            answerUuid.String(),
			Role:            answer.Role,
			Content:         answer.Content,
			Raw:             json.RawMessage(jsonMsg),
			UserID:          userID,
			CreatedBy:       userID,
			UpdatedBy:       userID,
		})
	if err != nil {
		return nil, fmt.Errorf("add ai answer %w", err)
	}
	return &answer_msg, err
}

func GetAiAnswerProxyLightsail(msgs []Message) (ChatCompletionResponse, error) {
	openaiRequest := OpenaiChatRequest{Model: "gpt-3.5-turbo", Messages: msgs}
	req_bytes, err := json.Marshal(openaiRequest)
	if err != nil {
		panic(err)
	}
	url := "http://14.214.224.18:8085/v1/chat/completions"
	req_str := string(req_bytes)
	fmt.Println(req_str)
	ai_req, _ := http.NewRequest("POST", url, strings.NewReader(req_str))

	ai_req.Header.Add("Content-Type", "application/json")
	ai_req.Header.Add("Authorization", "Basic ZWNob191c2VyX2E6N2UyMmE4YTQ0MTJmNDU2MDU0ODY5NDI1MjNlNDdkYjNmZTJlNzdhYQ==")

	ai_res, err := http.DefaultClient.Do(ai_req)

	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("request error: %w", err)
	}
	defer ai_res.Body.Close()
	var aiAnswer ChatCompletionResponse
	err = json.NewDecoder(ai_res.Body).Decode(&aiAnswer)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("decode request body err: %w", err)
	}
	return aiAnswer, nil
}

func GetAiAnswerOpenApi(msgs []Message) (ChatCompletionResponse, error) {
	openaiRequest := OpenaiChatRequest{Model: "gpt-3.5-turbo", Messages: msgs}
	req_bytes, err := json.Marshal(openaiRequest)
	log.Println(string(req_bytes))
	if err != nil {
		panic(err)
	}

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	// curl https://api.openai.com/v1/chat/completions \
	//   -H 'Content-Type: application/json' \
	//   -H 'Authorization: Bearer YOUR_API_KEY' \
	//   -d '{
	//   "model": "gpt-3.5-turbo",
	//   "messages": [{"role": "user", "content": "Hello!"}]
	// }'

	body := bytes.NewReader(req_bytes)

	ai_req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", body)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("request error: %w", err)
	}
	ai_req.Header.Set("Content-Type", "application/json")
	ai_req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", OPENAI_API_KEY))

	ai_res, err := http.DefaultClient.Do(ai_req)

	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("request error: %w", err)
	}
	defer ai_res.Body.Close()
	var aiAnswer ChatCompletionResponse
	err = json.NewDecoder(ai_res.Body).Decode(&aiAnswer)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("decode request body err: %w", err)
	}
	return aiAnswer, nil
}