package dialogflow

import (
	"context"
	"errors"
	"fmt"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type Logger interface {
	Warnf(format string, args ...interface{})
}

type Usecase interface {
	DetectIntentText(sessionID, text, languageCode string) (string, string, error)
}

type DialogFlow struct {
	ProjectID string
	logger    Logger
}

func New(logger Logger, projectID string) *DialogFlow {
	return &DialogFlow{logger: logger, ProjectID: projectID}
}

func (d DialogFlow) DetectIntentText(sessionID, text, languageCode string) (string, string, error) {
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return "", "", err
	}
	defer sessionClient.Close()

	if d.ProjectID == "" || sessionID == "" {
		return "", "", errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", d.ProjectID, sessionID))
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", d.ProjectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", "", err
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	fmt.Println(queryResult.Action)
	return fulfillmentText, queryResult.Action, nil
}
