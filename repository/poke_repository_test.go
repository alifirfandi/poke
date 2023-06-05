package repository_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"pokeapi/model"
	"pokeapi/repository"
	"testing"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	res := args.Get(0)
	err := args.Error(1)
	if res != nil {
		return res.(*http.Response), err
	}
	return nil, err
}

func TestGetAllPokemon(t *testing.T) {
	httpClient := new(MockHTTPClient)
	repo := repository.PokeRepository{
		DB:         nil,
		HTTPClient: httpClient,
	}

	expectedURL := "https://pokeapi.co/api/v2/pokemon?limit=10&offset=0"
	expectedRes := model.PokeDataSourceRes{
		Results: []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			{
				Name: "bulbasaur",
				URL:  "https://pokeapi.co/api/v2/pokemon/1/",
			},
			{
				Name: "charmander",
				URL:  "https://pokeapi.co/api/v2/pokemon/4/",
			},
			{
				Name: "squirtle",
				URL:  "https://pokeapi.co/api/v2/pokemon/7/",
			},
		},
	}

	responseBody, _ := json.Marshal(expectedRes)
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(responseBody)),
	}

	httpClient.On("Get", expectedURL).Return(response, nil)

	actualRes, err := repo.GetAllPokemon(0)

	assert.NoError(t, err)
	assert.Equal(t, expectedRes, actualRes)
	httpClient.AssertExpectations(t)
}

func TestGetOnePokemon(t *testing.T) {
	httpClient := new(MockHTTPClient)
	repo := repository.PokeRepository{
		DB:         nil,
		HTTPClient: httpClient,
	}

	expectedName := "bulbasaur"
	expectedURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", expectedName)
	expectedRes := model.PokeDetailDataSourceRes{
		Name: expectedName,
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Stat     struct {
				Name string `json:"name"`
			} `json:"stat"`
		}{
			{
				BaseStat: 45,
				Stat: struct {
					Name string `json:"name"`
				}{
					Name: "hp",
				},
			},
			{
				BaseStat: 49,
				Stat: struct {
					Name string `json:"name"`
				}{
					Name: "attack",
				},
			},
			{
				BaseStat: 49,
				Stat: struct {
					Name string `json:"name"`
				}{
					Name: "defense",
				},
			},
		},
	}

	responseBody, _ := json.Marshal(expectedRes)
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(responseBody)),
	}

	httpClient.On("Get", expectedURL).Return(response, nil)

	actualRes, err := repo.GetOnePokemon(expectedName)

	assert.NoError(t, err)
	assert.Equal(t, expectedRes, actualRes)
	httpClient.AssertExpectations(t)
}
