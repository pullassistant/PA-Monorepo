package api

import (
	"encoding/json"
	"github.com/google/go-github/github"
	"net/http"
	"strconv"
)

type InstallationsHandler struct {
	Config Config
}

func (i InstallationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := clientFromCtx(r.Context())

	purchases, _, err := client.Marketplace.ListMarketplacePurchasesForUser(r.Context(), nil)
	if err != nil {
		panic(err)
	}
	var plan github.MarketplacePlan
	if len(purchases) > 0 {
		plan = *purchases[0].Plan
	}

	installations, _, err := client.Apps.ListUserInstallations(r.Context(), nil)
	if err != nil {
		panic(err)
	}

	var installationsResponse []Installation = nil
	for _, v := range installations {
		// todo fix
		//accounts, _, err := appClientFromCtx(r.Context()).Marketplace.ListPlanAccountsForAccount(r.Context(), *v.Account.ID, nil)

		changePlanNumber := i.Config.AppConfig.MarketplaceFreePlanNumber
		if (*plan.ID == i.Config.AppConfig.MarketplaceOpenSourcePlanId || *plan.ID == i.Config.AppConfig.MarketplaceFreePlanId) && *v.Account.Type == "User" {
			changePlanNumber = i.Config.AppConfig.MarketplacePaidPlanNumber
		} else if (*plan.ID == i.Config.AppConfig.MarketplaceOpenSourcePlanId || *plan.ID == i.Config.AppConfig.MarketplaceFreePlanId) && *v.Account.Type == "Organization" {
			changePlanNumber = i.Config.AppConfig.MarketplaceOrgPlanNumber
		}

		installationsResponse = append(installationsResponse, Installation{
			ID:            *v.ID,
			HTMLURL:       *v.HTMLURL,
			Login:         *v.Account.Login,
			Accountid:     *v.Account.ID,
			Accounttype:   *v.Account.Type,
			AvatarURL:     *v.Account.AvatarURL,
			PlanID:        *plan.ID,
			PlanName:      *plan.Name,
			PlanChangeUrl: i.Config.AppConfig.MarketplaceUrl + "/upgrade/" + strconv.Itoa(int(changePlanNumber)) + "/" + strconv.Itoa(int(*v.Account.ID)),
		})
	}

	err = json.NewEncoder(w).Encode(Installations{Installations: installationsResponse})
	if err != nil {
		panic(err)
	}
}

type Installations struct {
	Installations []Installation `json:"installations"`
}
type Installation struct {
	ID            int64  `json:"id"`
	HTMLURL       string `json:"html_url"`
	Login         string `json:"login"`
	Accountid     int64  `json:"accountid"`
	Accounttype   string `json:"accounttype"`
	AvatarURL     string `json:"avatar_url"`
	PlanID        int64  `json:"plan_id"`
	PlanName      string `json:"plan_name"`
	PlanChangeUrl string `json:"plan_change_url"`
}
