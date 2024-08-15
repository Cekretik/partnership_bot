package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getClientSecretFile() string {
	return "credentials.json"
}

func GetGoogleSheetData(sheetID, readRange string) ([][]interface{}, error) {
	ctx := context.Background()
	b, err := os.ReadFile(getClientSecretFile())
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	return resp.Values, nil
}

func UpdateGoogleSheet(sheetID, writeRange string, values [][]interface{}) error {
	ctx := context.Background()
	b, err := os.ReadFile(getClientSecretFile())
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	vr := &sheets.ValueRange{
		Values: values,
	}

	_, err = srv.Spreadsheets.Values.Update(sheetID, writeRange, vr).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("unable to update data in sheet: %v", err)
	}

	return nil
}
