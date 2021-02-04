package lex

import (
	model "cloudrack-lambda-core/lex/model"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

type LexService struct {
	Svc     *lexmodelbuildingservice.LexModelBuildingService
	BotName string
}

func Init(botName string) LexService {
	return LexService{
		Svc:     lexmodelbuildingservice.New(session.New()),
		BotName: botName,
	}
}

func (ls LexService) GetSlotType(name string) (model.SlotConfig, error) {
	input := &lexmodelbuildingservice.GetSlotTypeInput{
		Name:    aws.String(name),
		Version: aws.String("$LATEST"),
	}

	result, err := ls.Svc.GetSlotType(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lexmodelbuildingservice.ErrCodeNotFoundException:
				fmt.Println(lexmodelbuildingservice.ErrCodeNotFoundException, aerr.Error())
			case lexmodelbuildingservice.ErrCodeLimitExceededException:
				fmt.Println(lexmodelbuildingservice.ErrCodeLimitExceededException, aerr.Error())
			case lexmodelbuildingservice.ErrCodeInternalFailureException:
				fmt.Println(lexmodelbuildingservice.ErrCodeInternalFailureException, aerr.Error())
			case lexmodelbuildingservice.ErrCodeBadRequestException:
				fmt.Println(lexmodelbuildingservice.ErrCodeBadRequestException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
			return model.SlotConfig{}, errors.New(aerr.Error())
		}
		return model.SlotConfig{}, err
	}
	slotType := model.SlotConfig{
		Name:     *result.Name,
		Checksum: *result.Checksum,
	}
	return slotType, nil
}

func (ls LexService) UpdateSlotType(slotConfig model.SlotConfig) error {
	//retreiving existing slotype first as overriding required the checksum
	curSlotType, err := ls.GetSlotType(slotConfig.Name)
	input := &lexmodelbuildingservice.PutSlotTypeInput{
		Description:            aws.String(slotConfig.Description),
		EnumerationValues:      []*lexmodelbuildingservice.EnumerationValue{},
		Name:                   aws.String(slotConfig.Name),
		ValueSelectionStrategy: aws.String(slotConfig.ValueSelectionStrategy),
	}
	if err == nil && curSlotType.Checksum != "" {
		input.Checksum = aws.String(curSlotType.Checksum)
	}
	for _, slot := range slotConfig.EnumerationValues {
		enumValue := &lexmodelbuildingservice.EnumerationValue{
			Value:    aws.String(slot.Value),
			Synonyms: []*string{},
		}
		for _, syn := range slot.Synonyms {
			enumValue.Synonyms = append(enumValue.Synonyms, aws.String(syn))
		}
		input.EnumerationValues = append(input.EnumerationValues, enumValue)
	}
	log.Printf("[LEX] Updating slot type: %+v", input)
	_, err = ls.Svc.PutSlotType(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lexmodelbuildingservice.ErrCodeConflictException:
				fmt.Println(lexmodelbuildingservice.ErrCodeConflictException, aerr.Error())
			case lexmodelbuildingservice.ErrCodeLimitExceededException:
				fmt.Println(lexmodelbuildingservice.ErrCodeLimitExceededException, aerr.Error())
			case lexmodelbuildingservice.ErrCodeInternalFailureException:
				fmt.Println(lexmodelbuildingservice.ErrCodeInternalFailureException, aerr.Error())
			case lexmodelbuildingservice.ErrCodeBadRequestException:
				fmt.Println(lexmodelbuildingservice.ErrCodeBadRequestException, aerr.Error())
			case lexmodelbuildingservice.ErrCodePreconditionFailedException:
				fmt.Println(lexmodelbuildingservice.ErrCodePreconditionFailedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
			return errors.New(aerr.Error())
		}
		fmt.Println(err.Error())
		return err
	}
	//updating the bot
	log.Printf("[LEX]Re building the bot %+v", input)
	return ls.RebuildBot()
}

func (ls LexService) RebuildBot() error {
	getIntput := &lexmodelbuildingservice.GetBotInput{
		VersionOrAlias: aws.String("$LATEST"),
		Name:           aws.String(ls.BotName),
	}
	getOutput, err := ls.Svc.GetBot(getIntput)
	if err != nil {
		log.Printf("[LEX] Error while Re building the bot %+v", err)
	}
	putInput := &lexmodelbuildingservice.PutBotInput{
		Name:                getOutput.Name,
		Checksum:            getOutput.Checksum,
		AbortStatement:      getOutput.AbortStatement,
		ChildDirected:       getOutput.ChildDirected,
		ClarificationPrompt: getOutput.ClarificationPrompt,
		Description:         getOutput.Description,
		DetectSentiment:     getOutput.DetectSentiment,
		//EnableModelImprovements:      getOutput.EnableModelImprovements,
		IdleSessionTTLInSeconds: getOutput.IdleSessionTTLInSeconds,
		Intents:                 getOutput.Intents,
		Locale:                  getOutput.Locale,
		//NluIntentConfidenceThreshold: getOutput.NluIntentConfidenceThreshold,
		VoiceId: getOutput.VoiceId,
	}
	_, err = ls.Svc.PutBot(putInput)
	if err != nil {
		log.Printf("[LEX] Error while Re building the bot %+v", err)
	}
	return err
}
