package api

import (
	"io/ioutil"

	"github.com/palantir/go-baseapp/baseapp"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/go-githubapp/githubapp"
)

type Config struct {
	Server baseapp.HTTPConfig `yaml:"server"`
	Github githubapp.Config   `yaml:"github"`

	AppConfig   AppConfig   `yaml:"app_configuration"`
	DebugConfig DebugConfig `yaml:"debug_config"`
}

type AppConfig struct {
	TrackingId                  string `yaml:"tracking_id"`
	AppName                     string `yaml:"app_name"`
	CommentAppName              string `yaml:"comment_app_name"`
	CookieSecret                string `yaml:"cookie_secret"`
	Debug                       bool   `yaml:"debug"`
	DashboardUrl                string `yaml:"dashboard_url"`
	MarketplaceUrl              string `yaml:"marketplace_url"`
	MarketplaceOpenSourcePlanId int64  `yaml:"marketplace_open_source_plan_id"`
	MarketplaceFreePlanId       int64  `yaml:"marketplace_free_plan_id"`
	MarketplaceFreePlanNumber   int    `yaml:"marketplace_free_plan_number"`
	MarketplacePaidPlanNumber   int    `yaml:"marketplace_paid_plan_number"`
	MarketplaceOrgPlanNumber    int    `yaml:"marketplace_org_plan_number"`
}

type DebugConfig struct {
}

func ReadConfig(path string) (*Config, error) {
	var c Config

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading server config file: %s", path)
	}

	if err := yaml.UnmarshalStrict(bytes, &c); err != nil {
		return nil, errors.Wrap(err, "failed parsing configuration file")
	}

	return &c, nil
}
