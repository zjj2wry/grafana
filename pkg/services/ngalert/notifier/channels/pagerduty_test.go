package channels

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting"
)

func TestPagerdutyNotifier(t *testing.T) {
	tmpl, err := template.FromGlobs("templates/default.tmpl")
	require.NoError(t, err)

	hostname, err := os.Hostname()
	require.NoError(t, err)

	cases := []struct {
		name         string
		json         string
		alerts       []*types.Alert
		expMsg       *pagerDutyMessage
		expInitError error
		expMsgError  error
	}{
		{
			name: "Default config with one alert",
			json: `{"integrationKey": "abcdefgh0123456789"}`,
			alerts: []*types.Alert{
				{
					Alert: model.Alert{
						Labels:      model.LabelSet{"alertname": "alert1", "lbl1": "val1"},
						Annotations: model.LabelSet{"ann1": "annv1"},
					},
				},
			},
			expMsg: &pagerDutyMessage{
				RoutingKey:  "abcdefgh0123456789",
				DedupKey:    "6e3538104c14b583da237e9693b76debbc17f0f8058ef20492e5853096cf8733",
				Description: "[firing:1]  (val1)",
				EventAction: "trigger",
				Payload: &pagerDutyPayload{
					Summary:   "[FIRING:1]  (val1)",
					Source:    hostname,
					Severity:  "critical",
					Class:     "todo_class",
					Component: "Grafana",
					Group:     "todo_group",
					CustomDetails: map[string]string{
						"firing":       "Labels:\n - alertname = alert1\n - lbl1 = val1\nAnnotations:\n - ann1 = annv1\nSource: \n",
						"num_firing":   "1",
						"num_resolved": "0",
						"resolved":     "",
					},
				},
				Client:    "Grafana",
				ClientURL: "http://localhost",
				Links:     []pagerDutyLink{{HRef: "http://localhost", Text: "External URL"}},
			},
			expInitError: nil,
			expMsgError:  nil,
		}, {
			name: "Custom config with multiple alerts",
			json: `{
				"integrationKey": "abcdefgh0123456789",
				"severity": "warning",
				"class": "{{ .Status }}",
				"component": "My Grafana",
				"group": "my_group"
			}`,
			alerts: []*types.Alert{
				{
					Alert: model.Alert{
						Labels:      model.LabelSet{"alertname": "alert1", "lbl1": "val1"},
						Annotations: model.LabelSet{"ann1": "annv1"},
					},
				}, {
					Alert: model.Alert{
						Labels:      model.LabelSet{"alertname": "alert1", "lbl1": "val2"},
						Annotations: model.LabelSet{"ann1": "annv2"},
					},
				},
			},
			expMsg: &pagerDutyMessage{
				RoutingKey:  "abcdefgh0123456789",
				DedupKey:    "6e3538104c14b583da237e9693b76debbc17f0f8058ef20492e5853096cf8733",
				Description: "[firing:2]  ",
				EventAction: "trigger",
				Payload: &pagerDutyPayload{
					Summary:   "[FIRING:2]  ",
					Source:    hostname,
					Severity:  "warning",
					Class:     "firing",
					Component: "My Grafana",
					Group:     "my_group",
					CustomDetails: map[string]string{
						"firing":       "Labels:\n - alertname = alert1\n - lbl1 = val1\nAnnotations:\n - ann1 = annv1\nSource: \nLabels:\n - alertname = alert1\n - lbl1 = val2\nAnnotations:\n - ann1 = annv2\nSource: \n",
						"num_firing":   "2",
						"num_resolved": "0",
						"resolved":     "",
					},
				},
				Client:    "Grafana",
				ClientURL: "http://localhost",
				Links:     []pagerDutyLink{{HRef: "http://localhost", Text: "External URL"}},
			},
			expInitError: nil,
			expMsgError:  nil,
		}, {
			name:         "Error in initing",
			json:         `{}`,
			expInitError: alerting.ValidationError{Reason: "Could not find integration key property in settings"},
		}, {
			name: "Error in building message",
			json: `{
				"integrationKey": "abcdefgh0123456789",
				"class": "{{ .Status }"
			}`,
			expMsgError: errors.New("failed to template PagerDuty message: template: :1: unexpected \"}\" in operand"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			settingsJSON, err := simplejson.NewJson([]byte(c.json))
			require.NoError(t, err)

			m := &models.AlertNotification{
				Name:     "pageduty_testing",
				Type:     "pagerduty",
				Settings: settingsJSON,
			}

			pn, err := NewPagerdutyNotifier(m, tmpl)
			if c.expInitError != nil {
				require.Error(t, err)
				require.Equal(t, c.expInitError.Error(), err.Error())
				return
			}
			require.NoError(t, err)

			ctx := notify.WithGroupKey(context.Background(), "alertname")
			ctx = notify.WithGroupLabels(ctx, model.LabelSet{"alertname": ""})
			msg, _, err := pn.buildPagerdutyMessage(ctx, types.Alerts(c.alerts...), c.alerts)
			if c.expMsgError != nil {
				require.Error(t, err)
				require.Equal(t, c.expMsgError.Error(), err.Error())
				return
			}
			require.NoError(t, err)

			require.Equal(t, c.expMsg, msg)
		})
	}
}
