package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/senchabot-dev/monorepo/apps/discord-bot/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB

type MySQL struct {
	DB *gorm.DB
}

func NewMySQL() *MySQL {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})

	if err != nil {
		panic("failed to connect database")
	}
	return &MySQL{
		DB: db,
	}
}

func (m *MySQL) SetDiscordBotConfig(ctx context.Context, serverId, key, value string) (bool, error) {
	var discordBotConfig []models.DiscordBotConfigs
	existConfig, err := m.GetDiscordBotConfig(ctx, serverId, key)
	if err != nil {
		return false, err
	}

	if existConfig != nil {
		result := m.DB.Model(&existConfig).Updates(models.DiscordBotConfigs{
			Key:   key,
			Value: value,
		})
		if result.Error != nil {
			return false, errors.New("(SetDiscordBotConfig) db.Updates Error:" + result.Error.Error())
		}

		return true, nil
	}

	discordBotConfig = append(discordBotConfig, models.DiscordBotConfigs{
		Key:      key,
		Value:    value,
		ServerID: serverId,
	})
	result := m.DB.Create(&discordBotConfig)
	if result.Error != nil {
		return false, errors.New("SetDiscordBotConfig db.Create Error:" + result.Error.Error())
	}

	return true, nil
}

func (m *MySQL) GetDiscordBotConfig(ctx context.Context, serverId, key string) (*models.DiscordBotConfigs, error) {
	var discordBotConfig []models.DiscordBotConfigs
	result := m.DB.Where("server_id = ?", serverId).Where("config_key = ?", key).Find(&discordBotConfig)
	if result.Error != nil {
		return nil, errors.New("(GetDiscordBotConfig) db.First Error:" + result.Error.Error())
	}

	if len(discordBotConfig) > 0 {
		return &discordBotConfig[0], nil
	}

	return nil, nil
}

func (m *MySQL) DeleteDiscordBotConfig(ctx context.Context, serverId, key string) (bool, error) {
	existConfig, err := m.GetDiscordBotConfig(ctx, serverId, key)
	if err != nil {
		return false, err
	}

	if existConfig == nil {
		return false, nil
	}

	result := m.DB.Model(&existConfig).Updates(models.DiscordBotConfigs{
		Key:   key,
		Value: "",
	})
	if result.Error != nil {
		return false, errors.New("(SetDiscordBotConfig) db.Updates Error:" + result.Error.Error())
	}

	return true, nil
}

func (m *MySQL) AddAnnouncementChannel(ctx context.Context, channelId, serverId, createdBy string) (bool, error) {
	var announcementChs []*models.DiscordAnnouncementChannels

	foundChannel, err := m.GetAnnouncementChannelByChannelId(ctx, channelId)
	if err != nil {
		return false, errors.New("(AddAnnouncementChannel) GetAnnouncementChannelByChannelId Error: " + err.Error())
	}

	if foundChannel != nil {
		return false, nil
	}

	announcementChs = append(announcementChs, &models.DiscordAnnouncementChannels{
		ChannelID: channelId,
		ServerID:  serverId,
		CreatedBy: createdBy,
	})

	result := m.DB.Create(&announcementChs)
	if result.Error != nil {
		return false, errors.New("(AddAnnouncementChannel) db.Create Error:" + result.Error.Error())
	}

	return true, nil
}

func (m *MySQL) GetAnnouncementChannels(ctx context.Context) ([]*models.DiscordAnnouncementChannels, error) {
	var announcementChs []*models.DiscordAnnouncementChannels

	result := m.DB.Find(&announcementChs)
	if result.Error != nil {
		return nil, errors.New("(GetAnnouncementChannels) db.Find Error:" + result.Error.Error())
	}

	return announcementChs, nil
}

func (m *MySQL) GetAnnouncementChannelByChannelId(ctx context.Context, channelId string) (*models.DiscordAnnouncementChannels, error) {
	var announcementChs []models.DiscordAnnouncementChannels
	result := m.DB.Where("channel_id = ?", channelId).Find(&announcementChs)
	if result.Error != nil {
		return nil, errors.New("(AddAnnouncementChannel) db.Find Error:" + result.Error.Error())
	}

	if len(announcementChs) == 0 {
		return nil, nil
	}

	return &announcementChs[0], nil
}

func (m *MySQL) GetAnnouncementChannelById(ctx context.Context, id int) (*models.DiscordAnnouncementChannels, error) {
	var announcementChs models.DiscordAnnouncementChannels

	result := m.DB.Where("id = ?", id).First(&announcementChs)
	if result.Error != nil {
		return nil, errors.New("(GetAnnouncementChannel) db.Find Error:" + result.Error.Error())
	}
	return &announcementChs, nil
}

func (m *MySQL) DeleteAnnouncementChannel(ctx context.Context, channelId string) (bool, error) {
	existAnnoCH, err := m.GetAnnouncementChannelByChannelId(ctx, channelId)
	if err != nil {
		return false, err
	}

	if existAnnoCH == nil {
		return false, nil
	}

	result := m.DB.Delete(&existAnnoCH)
	if result.Error != nil {
		return false, errors.New("(DeleteAnnouncementChannel) db.Delete Error:" + result.Error.Error())
	}

	return true, nil
}

func (m *MySQL) AddDiscordTwitchLiveAnnos(ctx context.Context, twitchUsername, twitchUserId, annoChannelId, annoServerId, createdBy string) (bool, error) {
	var twitchLiveAnnos []models.DiscordTwitchLiveAnnos

	twitchLiveAnno, err := m.GetDiscordTwitchLiveAnno(ctx, twitchUserId, annoServerId)
	if err != nil {
		return false, errors.New("(AddDiscordTwitchLiveAnnos) CheckDiscordTwitchLiveAnnos Error:" + err.Error())
	}
	if twitchLiveAnno != nil {
		result := m.DB.Model(&twitchLiveAnno).Updates(models.DiscordTwitchLiveAnnos{
			TwitchUsername: twitchUsername,
			TwitchUserID:   twitchUserId,
			AnnoChannelID:  annoChannelId,
			AnnoServerID:   annoServerId,
			CreatedBy:      createdBy,
		})
		if result.Error != nil {
			return false, errors.New("(AddDiscordTwitchLiveAnnos) db.Updates Error:" + result.Error.Error())
		}

		return false, nil
	}

	twitchLiveAnnos = append(twitchLiveAnnos, models.DiscordTwitchLiveAnnos{
		TwitchUsername: twitchUsername,
		TwitchUserID:   twitchUserId,
		AnnoChannelID:  annoChannelId,
		AnnoServerID:   annoServerId,
		Type:           1,
		CreatedBy:      createdBy,
	})

	result := m.DB.Create(&twitchLiveAnnos)
	if result.Error != nil {
		return false, errors.New("(AddDiscordTwitchLiveAnnos) db.Create Error:" + result.Error.Error())
	}

	return true, nil
}

func (m *MySQL) UpdateTwitchStreamerAnnoContent(ctx context.Context, twitchUsername, annoServerId string, annoContent *string) (bool, error) {

	twitchLiveAnno, err := m.GetDiscordTwitchLiveAnnoByUsername(ctx, twitchUsername, annoServerId)
	if err != nil {
		return false, errors.New("(UpdateTwitchStreamerAnnoContent) GetDiscordTwitchLiveAnnoByUsername Error:" + err.Error())
	}
	if twitchLiveAnno != nil {
		result := m.DB.Model(&twitchLiveAnno).Updates(models.DiscordTwitchLiveAnnos{
			AnnoContent: annoContent,
		})
		if result.Error != nil {
			return false, errors.New("(UpdateTwitchStreamerAnnoContent) db.Updates Error:" + result.Error.Error())
		}

		return true, nil
	}

	return false, nil
}

func (m *MySQL) UpdateTwitchStreamerLastAnnoDate(ctx context.Context, twitchUsername, annoServerId string, lastAnnoDate time.Time) (bool, error) {

	twitchLiveAnno, err := m.GetDiscordTwitchLiveAnnoByUsername(ctx, twitchUsername, annoServerId)
	if err != nil {
		return false, errors.New("(UpdateTwitchStreamerLastAnnoDate) GetDiscordTwitchLiveAnnoByUsername Error:" + err.Error())
	}
	if twitchLiveAnno != nil {
		result := m.DB.Model(&twitchLiveAnno).Updates(models.DiscordTwitchLiveAnnos{
			LastAnnoDate: &lastAnnoDate,
		})
		if result.Error != nil {
			return false, errors.New("(UpdateTwitchStreamerLastAnnoDate) db.Updates Error:" + result.Error.Error())
		}

		return true, nil
	}

	return false, nil
}

func (m *MySQL) GetTwitchStreamerLastAnnoDate(ctx context.Context, twitchUsername, annoServerId string) (*time.Time, error) {

	twitchLiveAnno, err := m.GetDiscordTwitchLiveAnnoByUsername(ctx, twitchUsername, annoServerId)
	if err != nil {
		return nil, errors.New("(CheckTwitchStreamerLastAnnoDate) GetDiscordTwitchLiveAnnoByUsername Error:" + err.Error())
	}
	if twitchLiveAnno != nil {
		return twitchLiveAnno.LastAnnoDate, nil
	}

	return nil, nil
}

func (m *MySQL) GetTwitchStreamerAnnoContent(ctx context.Context, twitchUsername, annoServerId string) (*string, error) {
	var twitchLiveAnnos []models.DiscordTwitchLiveAnnos

	result := m.DB.Where("twitch_username = ?", twitchUsername).Where("anno_server_id = ?", annoServerId).Find(&twitchLiveAnnos)
	if result.Error != nil {
		return nil, errors.New("(GetTwitchStreamerAnnoContent) db.Find Error:" + result.Error.Error())
	}

	if len(twitchLiveAnnos) == 0 {
		return nil, nil
	}

	return twitchLiveAnnos[0].AnnoContent, nil
}

func (m *MySQL) GetDiscordTwitchLiveAnno(ctx context.Context, twitchUserId, annoServerId string) (*models.DiscordTwitchLiveAnnos, error) {
	var twitchLiveAnnos []models.DiscordTwitchLiveAnnos

	result := m.DB.Where("twitch_user_id = ?", twitchUserId).Where("anno_server_id = ?", annoServerId).Find(&twitchLiveAnnos)
	if result.Error != nil {
		return nil, errors.New("(GetDiscordTwitchLiveAnno) db.Find Error:" + result.Error.Error())
	}

	if len(twitchLiveAnnos) == 0 {
		return nil, nil
	}

	return &twitchLiveAnnos[0], nil
}

func (m *MySQL) GetDiscordTwitchLiveAnnoByUsername(ctx context.Context, twitchUsername, annoServerId string) (*models.DiscordTwitchLiveAnnos, error) {
	var twitchLiveAnnos []models.DiscordTwitchLiveAnnos

	result := m.DB.Where("twitch_username = ?", twitchUsername).Where("anno_server_id = ?", annoServerId).Find(&twitchLiveAnnos)
	if result.Error != nil {
		return nil, errors.New("(GetDiscordTwitchLiveAnnoByUsername) db.Find Error:" + result.Error.Error())
	}

	if len(twitchLiveAnnos) == 0 {
		return nil, nil
	}

	return &twitchLiveAnnos[0], nil
}

func (m *MySQL) GetDiscordTwitchLiveAnnos(ctx context.Context, serverId string) ([]*models.DiscordTwitchLiveAnnos, error) {
	var twitchLiveAnnos []*models.DiscordTwitchLiveAnnos

	result := m.DB.Where("anno_server_id = ?", serverId).Find(&twitchLiveAnnos)
	if result.Error != nil {
		return nil, errors.New("(GetDiscordTwitchLiveAnnos) db.Find Error:" + result.Error.Error())
	}

	return twitchLiveAnnos, nil
}

func (m *MySQL) DeleteDiscordTwitchLiveAnno(ctx context.Context, twitchUserId string, serverId string) (bool, error) {
	existLiveAnno, err := m.GetDiscordTwitchLiveAnno(ctx, twitchUserId, serverId)
	if err != nil {
		return false, err
	}

	if existLiveAnno == nil {
		return false, nil
	}

	result := m.DB.Delete(&existLiveAnno)
	if result.Error != nil {
		return false, errors.New("(DeleteDiscordTwitchLiveAnno) db.Delete Error:" + result.Error.Error())
	}

	return true, nil
}

func (m *MySQL) CheckConfig(ctx context.Context, discordServerId string, configKey string, configValue string) bool {
	configData, err := m.GetDiscordBotConfig(ctx, discordServerId, configKey)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if configData != nil && configData.Value == configValue {
		return true
	}

	return false
}

func (m *MySQL) CreateBotActionActivity(ctx context.Context, botPlatformType, botActivity, discordServerId, activityAuthor string) error {
	botActionActivity := models.BotActionActivity{
		BotPlatformType: botPlatformType,
		BotActivity:     botActivity,
		DiscordServerID: &discordServerId,
		ActivityAuthor:  &activityAuthor,
	}

	result := m.DB.Create(&botActionActivity)

	if result.Error != nil {
		return errors.New("(CreateBotActionActivity) db.Create Error:" + result.Error.Error())
	}

	return nil
}

func (m *MySQL) SaveBotCommandActivity(context context.Context, activity, discordServerId, commandAuthor string) {
	check := m.CheckConfig(context, discordServerId, "bot_activity_enabled", "1")
	if !check {
		return
	}

	if err := m.CreateBotActionActivity(context, "discord", activity, discordServerId, commandAuthor); err != nil {
		fmt.Println(err.Error())
	}
}
