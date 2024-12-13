package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/spankie/gymshark/config"
	"github.com/spankie/gymshark/database/models"
)

func getDefaultConfig() config.Configuration {
	return config.Configuration{
		Port:       "8080",
		DbHost:     "localhost",
		DbPort:     5432,
		DbUsername: "spankie",
		DbPassword: "spankie",
		DbName:     "gymshark",
	}
}

func TestCreateOrderHandler(t *testing.T) { //nolint:cyclop
	conf := getDefaultConfig()
	dbService := createDBAndHTTPServer(t, &conf)

	testcases := []struct {
		name           string
		numberOfItems  int
		expectedCode   int
		assertResponse func(res *http.Response)
	}{
		{
			name:          "valid order with one item",
			numberOfItems: 1,
			expectedCode:  http.StatusCreated,
			assertResponse: func(res *http.Response) {
				respMap := map[string]interface{}{}
				err := json.NewDecoder(res.Body).Decode(&respMap)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}

				id, ok := respMap["id"]
				if !ok || id == 0 {
					t.Fatalf("expected id in response, got %v", respMap)
				}

				order, err := dbService.GetOrder(context.Background(), int(id.(float64)))
				if err != nil {
					t.Fatalf("expected nil error finding order in db but got: %v", err)
				}

				if order.NumberOfItems != 1 {
					t.Errorf("expected number of items to be 1, got %d", order.NumberOfItems)
				}

				if len(order.Shipping) != 1 {
					t.Errorf("expected number of shipping to be 1 but got: %v", len(order.Shipping))
				}

				shipping := order.Shipping[0]
				if shipping.PackSize != 250 {
					t.Errorf("expected shipping pack size to be 250 but got: %v", shipping.PackSize)
				}
			},
		},
		{
			name:          "valid order with 501 items",
			numberOfItems: 501,
			expectedCode:  http.StatusCreated,
			assertResponse: func(res *http.Response) {
				respMap := map[string]interface{}{}
				err := json.NewDecoder(res.Body).Decode(&respMap)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}

				id, ok := respMap["id"]
				if !ok || id == 0 {
					t.Fatalf("expected id in response, got %v", respMap)
				}

				order, err := dbService.GetOrder(context.Background(), int(id.(float64)))
				if err != nil {
					t.Fatalf("expected nil error finding order in db but got: %v", err)
				}

				if order.NumberOfItems != 501 {
					t.Errorf("expected number of items to be 1, got %d", order.NumberOfItems)
				}

				if len(order.Shipping) != 2 {
					t.Errorf("expected number of shipping to be 1 but got: %v", len(order.Shipping))
				}

				shipping1 := order.Shipping[0]
				if shipping1.PackSize != 500 {
					t.Errorf("expected shipping pack size to be 500 but got: %v", shipping1.PackSize)
				}

				shipping2 := order.Shipping[1]
				if shipping2.PackSize != 250 {
					t.Errorf("expected shipping pack size to be 250 but got: %v", shipping2.PackSize)
				}
			},
		},
		{
			name:          "empty number of items",
			numberOfItems: 0,
			expectedCode:  http.StatusBadRequest,
			assertResponse: func(res *http.Response) {
				var respMap map[string]interface{}
				err := json.NewDecoder(res.Body).Decode(&respMap)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}

				if _, ok := respMap["error"]; !ok {
					t.Errorf("expected error message in response, got %v", respMap)
				}
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Make a request to the server
			buf := bytes.NewBufferString(fmt.Sprintf(`{ "number_of_items": %v }`, tc.numberOfItems))
			resp, err := http.Post(fmt.Sprintf("http://localhost:%s/orders", conf.Port), "", buf)
			if err != nil {
				t.Errorf("failed to make request to server: %v", err)
			}

			t.Cleanup(func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("failed to close response body: %v", err)
				}
			})

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("expected status code %d, got %d", tc.expectedCode, resp.StatusCode)
			}

			tc.assertResponse(resp)
		})
	}
}

func TestGetOrderHandler(t *testing.T) { //nolint:cyclop
	conf := getDefaultConfig()
	dbService := createDBAndHTTPServer(t, &conf)

	ctx := context.Background()
	order := &models.Order{
		NumberOfItems: 1,
	}
	err := dbService.CreateOrder(ctx, order, nil)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	testcases := []struct {
		name           string
		id             int
		expectedCode   int
		assertResponse func(res *http.Response)
	}{
		{
			name:         "valid id",
			id:           order.ID,
			expectedCode: http.StatusOK,
			assertResponse: func(res *http.Response) {
				respMap := map[string]interface{}{}
				err = json.NewDecoder(res.Body).Decode(&respMap)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}
				if respMap["id"].(float64) != float64(order.ID) {
					t.Errorf("expected response to have id of %d, but got %v", order.ID, respMap["id"])
				}

			},
		},
		{
			name:         "invalid id",
			id:           100,
			expectedCode: http.StatusNotFound,
			assertResponse: func(res *http.Response) {
				respMap := map[string]interface{}{}
				err := json.NewDecoder(res.Body).Decode(&respMap)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}

				if _, ok := respMap["error"]; !ok {
					t.Errorf("expected error message in response, got %v", respMap)
				}
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Make a request to the server
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/orders/%d", conf.Port, tc.id))
			if err != nil {
				t.Errorf("failed to make request to server: %v", err)
			}

			t.Cleanup(func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("failed to close response body: %v", err)
				}
			})

			if resp.StatusCode != tc.expectedCode {
				t.Fatalf("expected status code %d, got %d", tc.expectedCode, resp.StatusCode)
			}

			tc.assertResponse(resp)
		})
	}
}

func TestGetOrderShippingHandler(t *testing.T) {
	conf := getDefaultConfig()
	dbService := createDBAndHTTPServer(t, &conf)

	ctx := context.Background()
	order := &models.Order{
		NumberOfItems: 501,
	}
	err := dbService.CreateOrder(ctx, order, []*models.OrderShipping{
		{
			PackSize:             500,
			ShippingPackQuantity: 1,
		},
		{
			PackSize:             250,
			ShippingPackQuantity: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/shipping", conf.Port))
	if err != nil {
		t.Errorf("failed to make request to server: %v", err)
	}

	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	shipping := []models.OrderShipping{}
	err = json.NewDecoder(resp.Body).Decode(&shipping)
	if err != nil {
		t.Fatalf("unable to decode response: %v", err)
	}

	if len(shipping) != 2 {
		t.Fatalf("expected len of shipping to be %d but got %v", 2, len(shipping))
	}
}
