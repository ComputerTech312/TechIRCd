package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// OperClass defines an operator class with specific permissions
type OperClass struct {
	Name        string   `json:"name"`
	Rank        int      `json:"rank"`        // Higher number = higher rank
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	Inherits    string   `json:"inherits"`    // Inherit permissions from another class
	Color       string   `json:"color"`       // Color for display purposes
	Symbol      string   `json:"symbol"`      // Symbol to display (*, @, &, etc.)
}

// Oper defines an individual operator
type Oper struct {
	Name        string   `json:"name"`
	Password    string   `json:"password"`
	Host        string   `json:"host"`
	Class       string   `json:"class"`
	Flags       []string `json:"flags"`      // Additional per-user flags
	MaxClients  int      `json:"max_clients"` // Max clients this oper can handle
	Expires     string   `json:"expires"`     // Expiration date (optional)
	Contact     string   `json:"contact"`     // Contact information
	LastSeen    string   `json:"last_seen"`   // Last time this oper was online
}

// RankNames defines custom names for rank levels
type RankNames struct {
	Rank1 string `json:"rank_1"` // Default: Helper
	Rank2 string `json:"rank_2"` // Default: Moderator  
	Rank3 string `json:"rank_3"` // Default: Operator
	Rank4 string `json:"rank_4"` // Default: Administrator
	Rank5 string `json:"rank_5"` // Default: Owner
	// Support for custom ranks beyond 5
	CustomRanks map[string]int `json:"custom_ranks"` // "CustomName": 6
}

// OperConfig holds the complete operator configuration
type OperConfig struct {
	Classes []OperClass `json:"classes"`
	Opers   []Oper      `json:"opers"`
	RankNames RankNames `json:"rank_names"`
	Settings struct {
		RequireSSL          bool     `json:"require_ssl"`
		MaxFailedAttempts   int      `json:"max_failed_attempts"`
		LockoutDuration     int      `json:"lockout_duration_minutes"`
		AllowedCommands     []string `json:"allowed_commands"`
		LogOperActions      bool     `json:"log_oper_actions"`
		NotifyOnOperUp      bool     `json:"notify_on_oper_up"`
		AutoExpireInactive  int      `json:"auto_expire_inactive_days"`
		RequireTwoFactor    bool     `json:"require_two_factor"`
	} `json:"settings"`
}

// LoadOperConfig loads operator configuration from file
func LoadOperConfig(filename string) (*OperConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read oper config file: %v", err)
	}

	var config OperConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse oper config file: %v", err)
	}

	return &config, nil
}

// SaveOperConfig saves operator configuration to file
func SaveOperConfig(config *OperConfig, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal oper config: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write oper config file: %v", err)
	}

	return nil
}

// GetOperClass returns the operator class by name
func (oc *OperConfig) GetOperClass(className string) *OperClass {
	for i := range oc.Classes {
		if oc.Classes[i].Name == className {
			return &oc.Classes[i]
		}
	}
	return nil
}

// GetOper returns the operator by name
func (oc *OperConfig) GetOper(operName string) *Oper {
	for i := range oc.Opers {
		if oc.Opers[i].Name == operName {
			return &oc.Opers[i]
		}
	}
	return nil
}

// GetOperPermissions returns all permissions for an operator (including inherited)
func (oc *OperConfig) GetOperPermissions(operName string) []string {
	oper := oc.GetOper(operName)
	if oper == nil {
		return nil
	}

	class := oc.GetOperClass(oper.Class)
	if class == nil {
		return oper.Flags
	}

	permissions := make(map[string]bool)
	
	// Add class permissions
	for _, perm := range class.Permissions {
		permissions[perm] = true
	}
	
	// Add inherited permissions
	if class.Inherits != "" {
		inherited := oc.GetOperClass(class.Inherits)
		if inherited != nil {
			for _, perm := range inherited.Permissions {
				permissions[perm] = true
			}
		}
	}
	
	// Add individual flags
	for _, flag := range oper.Flags {
		permissions[flag] = true
	}
	
	// Convert back to slice
	result := make([]string, 0, len(permissions))
	for perm := range permissions {
		result = append(result, perm)
	}
	
	return result
}

// GetRankName returns the custom name for a rank level
func (oc *OperConfig) GetRankName(rank int) string {
	switch rank {
	case 1:
		if oc.RankNames.Rank1 != "" {
			return oc.RankNames.Rank1
		}
		return "Helper"
	case 2:
		if oc.RankNames.Rank2 != "" {
			return oc.RankNames.Rank2
		}
		return "Moderator"
	case 3:
		if oc.RankNames.Rank3 != "" {
			return oc.RankNames.Rank3
		}
		return "Operator"
	case 4:
		if oc.RankNames.Rank4 != "" {
			return oc.RankNames.Rank4
		}
		return "Administrator"
	case 5:
		if oc.RankNames.Rank5 != "" {
			return oc.RankNames.Rank5
		}
		return "Owner"
	default:
		// Check custom ranks
		for name, rankNum := range oc.RankNames.CustomRanks {
			if rankNum == rank {
				return name
			}
		}
		return fmt.Sprintf("Rank %d", rank)
	}
}

// GetOperRankName returns the rank name for an operator
func (oc *OperConfig) GetOperRankName(operName string) string {
	oper := oc.GetOper(operName)
	if oper == nil {
		return "Unknown"
	}
	
	class := oc.GetOperClass(oper.Class)
	if class == nil {
		return "Unknown"
	}
	
	return oc.GetRankName(class.Rank)
}

// SetCustomRankName allows adding custom rank names at runtime
func (oc *OperConfig) SetCustomRankName(name string, rank int) {
	if oc.RankNames.CustomRanks == nil {
		oc.RankNames.CustomRanks = make(map[string]int)
	}
	oc.RankNames.CustomRanks[name] = rank
}

// HasPermission checks if an operator has a specific permission
func (oc *OperConfig) HasPermission(operName, permission string) bool {
	permissions := oc.GetOperPermissions(operName)
	for _, perm := range permissions {
		if perm == permission || perm == "*" { // * grants all permissions
			return true
		}
	}
	return false
}

// GetOperRank returns the rank of an operator
func (oc *OperConfig) GetOperRank(operName string) int {
	oper := oc.GetOper(operName)
	if oper == nil {
		return 0
	}
	
	class := oc.GetOperClass(oper.Class)
	if class == nil {
		return 0
	}
	
	return class.Rank
}

// CanOperateOn checks if oper1 can perform actions on oper2 (based on rank)
func (oc *OperConfig) CanOperateOn(oper1Name, oper2Name string) bool {
	rank1 := oc.GetOperRank(oper1Name)
	rank2 := oc.GetOperRank(oper2Name)
	
	// Higher rank can operate on lower rank
	// Same rank cannot operate on each other (unless they have override permission)
	return rank1 > rank2 || oc.HasPermission(oper1Name, "override_rank")
}

// DefaultOperConfig returns a default operator configuration
func DefaultOperConfig() *OperConfig {
	return &OperConfig{
		Classes: []OperClass{
			{
				Name:        "helper",
				Rank:        1,
				Description: "Helper - Basic moderation commands",
				Permissions: []string{"kick", "topic", "mode_channel"},
				Color:       "green",
				Symbol:      "%",
			},
			{
				Name:        "moderator", 
				Rank:        2,
				Description: "Moderator - Channel and user management",
				Permissions: []string{"ban", "unban", "kick", "mute", "topic", "mode_channel", "mode_user", "who_override"},
				Inherits:    "helper",
				Color:       "blue",
				Symbol:      "@",
			},
			{
				Name:        "operator",
				Rank:        3,
				Description: "Operator - Server management commands",
				Permissions: []string{"kill", "gline", "rehash", "connect", "squit", "wallops", "operwall"},
				Inherits:    "moderator",
				Color:       "red",
				Symbol:      "*",
			},
			{
				Name:        "admin",
				Rank:        4,
				Description: "Administrator - Full server control",
				Permissions: []string{"*"}, // All permissions
				Color:       "purple",
				Symbol:      "&",
			},
			{
				Name:        "owner",
				Rank:        5,
				Description: "Server Owner - Ultimate authority",
				Permissions: []string{"*", "override_rank", "shutdown", "restart"},
				Color:       "gold",
				Symbol:      "~",
			},
		},
		Opers: []Oper{
			{
				Name:       "admin",
				Password:   "changeme",
				Host:       "*@localhost",
				Class:      "admin",
				MaxClients: 1000,
				Contact:    "admin@example.com",
			},
		},
		RankNames: RankNames{
			Rank1: "Helper",        // Customizable: "Support Staff", "Junior Mod", etc.
			Rank2: "Moderator",     // Customizable: "Channel Mod", "Guard", etc.
			Rank3: "Operator",      // Customizable: "IRC Operator", "Senior Staff", etc. 
			Rank4: "Administrator", // Customizable: "Admin", "Server Admin", etc.
			Rank5: "Owner",         // Customizable: "Root", "Founder", etc.
			CustomRanks: map[string]int{
				// Example: "Super Admin": 6,
				// Example: "Network Founder": 10,
			},
		},
		Settings: struct {
			RequireSSL          bool     `json:"require_ssl"`
			MaxFailedAttempts   int      `json:"max_failed_attempts"`
			LockoutDuration     int      `json:"lockout_duration_minutes"`
			AllowedCommands     []string `json:"allowed_commands"`
			LogOperActions      bool     `json:"log_oper_actions"`
			NotifyOnOperUp      bool     `json:"notify_on_oper_up"`
			AutoExpireInactive  int      `json:"auto_expire_inactive_days"`
			RequireTwoFactor    bool     `json:"require_two_factor"`
		}{
			RequireSSL:          false,
			MaxFailedAttempts:   3,
			LockoutDuration:     30,
			AllowedCommands:     []string{"OPER", "KILL", "GLINE", "REHASH", "WALLOPS"},
			LogOperActions:      true,
			NotifyOnOperUp:      true,
			AutoExpireInactive:  365,
			RequireTwoFactor:    false,
		},
	}
}
