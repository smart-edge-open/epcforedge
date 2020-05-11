package af

import (
	"context"
	"net/http"
)

// EventSubscResponse struct anyof EventSusbscReqData, EventNotification,
// ProblemDetails
type EventSubscResponse struct {
	eventSubscReq *EventsSubscReqData
	evsNotif      *EventsNotification
	probDetails   *ProblemDetails
	httpResp      *http.Response
	locationURI   string
}

// PcfPAResponse contains Policy auth response from PCF
type PcfPAResponse struct {
	appSessCtx  *AppSessionContext
	probDetails *ProblemDetails
	httpResp    *http.Response
	locationURI string
}

// pcfPolicyAuthAPI defines the interfaces that are exposed for POLICY AUTH
type pcfPolicyAuthAPI interface {
	DeleteAppSession(ctx context.Context, appSessionID string,
		eventSubscReq *EventsSubscReqData) (PcfPAResponse, error)

	GetAppSession(ctx context.Context, appSessionID string) (
		PcfPAResponse, error)

	ModAppSession(ctx context.Context, appSessionID string,
		appSessionContextUpdateData AppSessionContextUpdateData) (
		PcfPAResponse, error)

	PostAppSessions(ctx context.Context,
		appSessionContext AppSessionContext) (PcfPAResponse, error)

	UpdateEventsSubsc(ctx context.Context, appSessionID string,
		eventSubscReq *EventsSubscReqData) (
		EventSubscResponse, error)

	DeleteEventsSubsc(ctx context.Context, appSessionID string) (
		EventSubscResponse, error)
}
