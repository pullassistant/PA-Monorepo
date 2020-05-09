package api

import (
	"encoding/json"
	"net/http"
)

type DashboardConfigHandler struct {
	Config Config
}

func (h DashboardConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(DashboardConfig{
		MarketplaceUrl:              h.Config.AppConfig.MarketplaceUrl,
		MarketplaceOpenSourcePlanId: int(h.Config.AppConfig.MarketplaceOpenSourcePlanId),
		MarketplaceFreePlanId:       int(h.Config.AppConfig.MarketplaceFreePlanId),
		Debug:                       h.Config.AppConfig.Debug,
	})
	if err != nil {
		panic(err)
	}
}

type DashboardConfig struct {
	MarketplaceUrl              string `json:"marketplace_url"`
	MarketplaceOpenSourcePlanId int    `json:"marketplace_open_source_plan_id"`
	MarketplaceFreePlanId       int    `json:"marketplace_free_plan_id"`
	Debug                       bool   `json:"debug"`
}
