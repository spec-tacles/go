package types

import "github.com/bwmarrin/snowflake"

// AuditLog represents an audit log on Discord
type AuditLog struct {
	Webhooks        []Webhook       `json:"webhooks"`
	Users           []User          `json:"users"`
	AuditLogEntries []AuditLogEntry `json:"audit_log_entries"`
}

// AuditLogEntry represents an audit log entry on Discord
type AuditLogEntry struct {
	TargetID   string         `json:"target_id"`
	Changes    []AuditChange  `json:"changes,omitempty"`
	UserID     snowflake.ID   `json:"user_id"`
	ID         snowflake.ID   `json:"id"`
	ActionType int            `json:"action_type"`
	Options    AuditEntryInfo `json:"options,omitempty"`
	Reason     string         `json:"reason"`
}

// Audit log events
const (
	AuditLogGuildUpdate            = 1
	AuditLogChannelCreate          = 10
	AuditLogChannelUpdate          = 11
	AuditLogChannelDelete          = 12
	AuditLogChannelOverwriteCreate = 13
	AuditLogChannelOverwriteUpdate = 14
	AuditLogChannelOverwriteDelete = 15
	AuditLogMemberKick             = 20
	AuditLogMemberPrune            = 21
	AuditLogMemberBanAdd           = 22
	AuditLogMemberBanRemove        = 23
	AuditLogMemberUpdate           = 24
	AuditLogMemberRoleUpdate       = 25
	AuditLogRoleCreate             = 30
	AuditLogRoleUpdate             = 31
	AuditLogRoleDelete             = 32
	AuditLogInviteCreate           = 40
	AuditLogInviteUpdate           = 41
	AuditLogInviteDelete           = 42
	AuditLogWebhookCreate          = 50
	AuditLogWebhookUpdate          = 51
	AuditLogWebhookDelete          = 52
	AuditLogEmojiCreate            = 60
	AuditLogEmojiUpdate            = 61
	AuditLogEmojiDelete            = 62
	AuditLogMessageDelete          = 72
)

// AuditEntryInfo represents optional audit log entry info on Discord
type AuditEntryInfo struct {
	DeleteMemberDays string       `json:"delete_member_days"`
	MembersRemoved   string       `json:"members_removed"`
	ChannelID        snowflake.ID `json:"channel_id"`
	Count            string       `json:"count"`
	ID               snowflake.ID `json:"id"`
	Type             string       `json:"type"`
	RoleName         string       `json:"role_name"`
}

// AuditChange represents an audit log change on Discord
type AuditChange struct {
	NewValue interface{} `json:"new_value"`
	OldValue interface{} `json:"old_value"`
	Key      string      `json:"key"`
}

// Audit log change keys
const (
	AuditChangeName                        = "name"
	AuditChangeIconHash                    = "icon_hash"
	AuditChangeSplashHash                  = "splash_hash"
	AuditChangeOwnerID                     = "owner_id"
	AuditChangeRegion                      = "region"
	AuditChangeAFKChannelID                = "afk_channel_id"
	AuditChangeAFKTimeout                  = "afk_timeout"
	AuditChangeMFALevel                    = "mfa_level"
	AuditChangeVerificationLevel           = "verification_level"
	AuditChangeExplicitContentFilter       = "explicit_content_filter"
	AuditChangeDefaultMessageNotifications = "default_message_notifications"
	AuditChangeVanityURLCode               = "vanity_url_code"
	AuditChangeAdd                         = "$add"
	AuditChangeRemove                      = "$remove"
	AuditChangePruneDeleteDays             = "prune_delete_days"
	AuditChangeWidgetEnabled               = "widget_enabled"
	AuditChangeWidgetChannelID             = "widget_channel_id"
	AuditChangePosition                    = "position"
	AuditChangeTopic                       = "topic"
	AuditChangeBitrate                     = "bitrate"
	AuditChangePermissionOverwrites        = "permission_overwrites"
	AuditChangeNSFW                        = "nsfw"
	AuditChangeApplicationID               = "application_id"
	AuditChangePermissions                 = "permissions"
	AuditChangeColor                       = "color"
	AuditChangeHoist                       = "hoist"
	AuditChangeMentionable                 = "mentionable"
	AuditChangeAllow                       = "allow"
	AuditChangeDeny                        = "deny"
	AuditChangeCode                        = "code"
	AuditChangeChannelID                   = "channel_id"
	AuditChangeInviterID                   = "inviter_id"
	AuditChangeMaxUses                     = "max_uses"
	AuditChangeUses                        = "uses"
	AuditChangeMaxAge                      = "max_age"
	AuditChangeTemporary                   = "temporary"
	AuditChangeDeaf                        = "deaf"
	AuditChangeMute                        = "mute"
	AuditChangeNick                        = "nick"
	AuditChangeAvatarHash                  = "avatar_hash"
	AuditChangeID                          = "id"
	AuditChangeType                        = "type"
)
