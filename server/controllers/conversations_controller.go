package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/config"
	"server/models"
	"server/utils"
)

// CreateConversationRequest representa a payload para criar uma conversa
type CreateConversationRequest struct {
	ParticipantID string `json:"participant_id" binding:"required"`
}

// SendMessageRequest representa a payload para enviar uma mensagem
type SendMessageRequest struct {
	EncryptedContents map[string]models.ElGamalContent `json:"encrypted_contents" binding:"required"`
}

type ConversationResponse struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`           // Nome do grupo ou do contato
	LastMessage  *Message  `json:"last_message"`   // Pode ser nulo
	UnreadCount  int       `json:"unread_count"`
	UpdatedAt    time.Time `json:"updated_at"`     // Baseado na última mensagem ou criação
}

// ListConversations lista todas as conversas do usuário autenticado
func ListConversations(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var conversations []ConversationResponse

	// Buscar todas as conversas do usuário com suas últimas mensagens
	err = config.DB.Raw(`
		WITH LastMessages AS (
			SELECT DISTINCT ON (m.conversation_id)
				m.conversation_id,
				m.id as message_id,
				m.created_at as message_date,
				mr.encrypted_content,
				m.sender_id
			FROM messages m
			JOIN message_recipients mr ON mr.message_id = m.id
			WHERE mr.recipient_id = ?
			ORDER BY m.conversation_id, m.created_at DESC
		),
		UnreadCounts AS (
			SELECT m.conversation_id, COUNT(*) as unread
			FROM messages m
			JOIN message_recipients mr ON mr.message_id = m.id
			WHERE mr.recipient_id = ? AND mr.status = 'SENT'
			GROUP BY m.conversation_id
		)
		SELECT
			c.id,
			c.type,
			CASE
				WHEN c.type = 'GROUP' THEN g.name
				ELSE u.username
			END as name,
			lm.message_id,
			lm.encrypted_content,
			lm.message_date,
			COALESCE(uc.unread, 0) as unread_count,
			COALESCE(lm.message_date, c.created_at) as updated_at
		FROM conversations c
		JOIN conversation_participants cp ON cp.conversation_id = c.id
		LEFT JOIN groups g ON g.conversation_id = c.id
		LEFT JOIN conversation_participants cp2 ON cp2.conversation_id = c.id AND cp2.user_id != ?
		LEFT JOIN users u ON u.id = cp2.user_id
		LEFT JOIN LastMessages lm ON lm.conversation_id = c.id
		LEFT JOIN UnreadCounts uc ON uc.conversation_id = c.id
		WHERE cp.user_id = ?
		ORDER BY COALESCE(lm.message_date, c.created_at) DESC
	`, userID, userID, userID, userID).Scan(&conversations).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar conversas"})
		return
	}

	c.JSON(http.StatusOK, conversations)
}


// GetConversation recupera os detalhes básicos de uma conversa específica para o usuário autenticado
func GetConversation(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	conversationID := c.Param("id")
	includeMessages := c.Query("include_messages") == "true"

	// Buscar conversa com participantes
	var conversation models.Conversation
	query := config.DB.
		Joins("JOIN conversation_participants ON conversation_participants.conversation_id = conversations.id").
		Where("conversations.id = ? AND conversation_participants.user_id = ?", conversationID, userID).
		Preload("Participants.User")

	if includeMessages {
		query = query.Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Messages.Recipients", "recipient_id = ?", userID)
	}

	if err := query.First(&conversation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversa não encontrada"})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

// SendMessage envia uma nova mensagem para uma conversa específica
func SendMessage(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	conversationID := c.Param("id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da conversa é obrigatório"})
		return
	}

	// Verificar se o usuário é participante da conversa
	var participant models.ConversationParticipant
	if err := config.DB.Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&participant).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Usuário não é participante desta conversa"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar se há conteúdo criptografado para todos os participantes
	var participants []models.ConversationParticipant
	if err := config.DB.Where("conversation_id = ?", conversationID).Find(&participants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar participantes"})
		return
	}

	for _, p := range participants {
		if _, exists := req.EncryptedContents[p.UserID]; !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Conteúdo criptografado faltando para algum participante"})
			return
		}
	}

	// Iniciar transação
	tx := config.DB.Begin()

	// Criar mensagem
	message := models.Message{
		ID:             utils.GenerateUUID(),
		ConversationID: conversationID,
		SenderID:       userID,
		CreatedAt:      time.Now(),
	}

	if err := tx.Create(&message).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar mensagem"})
		return
	}

	// Criar MessageRecipient para cada destinatário
	for recipientID, encryptedContent := range req.EncryptedContents {
		messageRecipient := models.MessageRecipient{
			ID:               utils.GenerateUUID(),
			MessageID:        message.ID,
			RecipientID:      recipientID,
			EncryptedContent: encryptedContent,
			Status:           "SENT",
			StatusUpdatedAt:  time.Now(),
		}

		if err := tx.Create(&messageRecipient).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar mensagem para destinatário"})
			return
		}
	}

	// Commit da transação
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao finalizar envio da mensagem"})
		return
	}

	// Notificar os participantes via WebSocket (opcional)
	// Implementar lógica de notificação se necessário

	c.JSON(http.StatusCreated, message)
}

// UpdateMessageStatus atualiza o status de uma mensagem para o usuário autenticado
func UpdateMessageStatus(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	messageID := c.Param("messageId")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da mensagem é obrigatório"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=RECEIVED READ"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := config.DB.Model(&models.MessageRecipient{}).
		Where("message_id = ? AND recipient_id = ?", messageID, userID).
		Updates(map[string]interface{}{
			"status":             req.Status,
			"status_updated_at":  time.Now(),
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar status"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mensagem não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status atualizado com sucesso"})
}