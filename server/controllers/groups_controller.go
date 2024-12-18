package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"server/config"
	"server/models"
	"server/utils"
)

type CreateGroupRequest struct {
	Name           string   `json:"name" binding:"required"`
	ParticipantIDs []string `json:"participant_ids" binding:"required"`
}

func CreateGroup(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar os IDs reais dos contatos
	var contacts []models.Contact
	if err := config.DB.Where("user_id = ? AND id IN ?", userID, req.ParticipantIDs).Find(&contacts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar contatos"})
		return
	}

	// Verificar se todos os contatos foram encontrados
	if len(contacts) != len(req.ParticipantIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Um ou mais contatos não foram encontrados"})
		return
	}

	// Extrair os IDs reais dos contatos
	realParticipantIDs := make([]string, len(contacts))
	for i, contact := range contacts {
		realParticipantIDs[i] = contact.ContactID
	}

	tx := config.DB.Begin()

	conversation := models.Conversation{
		ID:        utils.GenerateUUID(),
		Type:      "GROUP",
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&conversation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar conversa"})
		return
	}

	group := models.Group{
		ConversationID: conversation.ID,
		Name:          req.Name,
		AdminID:       userID,
		CreatedAt:     time.Now(),
	}

	if err := tx.Create(&group).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar grupo"})
		return
	}

	// Adicionar todos os participantes, incluindo o criador do grupo
	allParticipants := append(realParticipantIDs, userID)
	for _, pid := range allParticipants {
		participant := models.ConversationParticipant{
			ID:             utils.GenerateUUID(),
			ConversationID: conversation.ID,
			UserID:         pid,
			JoinedAt:       time.Now(),
		}

		if err := tx.Create(&participant).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar participante"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao finalizar criação do grupo"})
		return
	}

	c.JSON(http.StatusCreated, group)
}