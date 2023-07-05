package main

import (
	"log"

	"auth/cmd"

	"github.com/spf13/cobra"
)

type AuthServiceCli struct {
	RootCmd *cobra.Command
}

func ProjectService() *AuthServiceCli {
	cli := &AuthServiceCli{
		RootCmd: &cobra.Command{
			Use:   "project-service",
			Short: "Project Service CLI",
			FParseErrWhitelist: cobra.FParseErrWhitelist{
				UnknownFlags: true,
			},
			// no need to provide the default cobra completion command
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}

	// hide the default help command (allow only `--help` flag)
	cli.RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// register system commands
	cli.RootCmd.AddCommand(cmd.Serve())

	return cli
}

func (cli *AuthServiceCli) Start() error {
	if err := cli.RootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func main() {
	app := ProjectService()
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// package main
//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"os"
// 	"strings"
//
// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// )
//
// // Credentials ...
// type Credentials struct {
// 	Email string `json:"username" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// }
//
// // isValid ...
// func (credentials Credentials) isValid() error {
// 	// check if both are valid strings or not
//
// 	return nil
// }
//
// func getAccessToken() (string, error) {
// 	// get access token from keycloak
// 	KeycloakServer := os.Getenv("KEYCLOAK_SERVER")
// 	KeycloakRealm := os.Getenv("KEYCLOAK_REALM")
// 	KeycloakClientId := os.Getenv("KEYCLOAK_CLIENT_ID")
// 	KeycloakClientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
//
// 	url := KeycloakServer + "/realms/" + KeycloakRealm + "/protocol/openid-connect/token"
// 	// show the url
// 	fmt.Println(url)
// 	method := "POST"
// 	payload := strings.NewReader("grant_type=client_credentials&client_id=" + KeycloakClientId + "&client_secret=" + KeycloakClientSecret)
// 	req, err := http.NewRequest(method, url, payload)
//
// 	if err != nil {
// 		fmt.Println(err)
// 		return "", err
// 	}
//
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 		return "", err
// 	}
//
// 	defer res.Body.Close()
// 	// 	extract the access_token from the response body
//
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return "", err
// 	}
//
// 	var result map[string]interface{}
// 	json.Unmarshal([]byte(body), &result)
//
// 	token := result["access_token"]
// 	if token == nil {
// 		return "", err
// 	}
//
// 	return token.(string), nil
// }
//
// func createKeycloakUser(username string, password string, token string) error {
// 	// create a keycloak user
// 	KeycloakServer := os.Getenv("KEYCLOAK_SERVER")
// 	KeycloakRealm := os.Getenv("KEYCLOAK_REALM")
// 	url := KeycloakServer + "/admin/realms/" + KeycloakRealm + "/users"
// 	method := "POST"
//
// 	// create user with dashboard and project roles
//
// 	payload := strings.NewReader("{\"username\":\"" + username + "\",\"enabled\":true,\"credentials\":[{\"type\":\"password\",\"value\":\"" + password + "\",\"temporary\":false}],\"realmRoles\":[\"dashboard\",\"project\"], \"access\": {\"manageGroupMembership\": true, \"view\": true, \"mapRoles\": true, \"impersonate\": true, \"manage\": true},\"groups\":[\"dashboard\",\"project\"]}")
//
// 	req, err := http.NewRequest(method, url, payload)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
//
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("Authorization", "Bearer "+token)
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
//
// 	defer res.Body.Close()
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
//
// 	fmt.Println(string(body))
//
// 	// attach autorization roles to the user
// 	// get the user id from the response
// 	var result map[string]interface{}
// 	json.Unmarshal([]byte(body), &result)
//
// 	userId := result["id"]
// 	if userId == nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func login(username string, password string) (string, error) {
// 	// 	get access token from keycloak using grant type as password
// 	KeycloakServer := os.Getenv("KEYCLOAK_SERVER")
// 	KeycloakRealm := os.Getenv("KEYCLOAK_REALM")
// 	KeycloakClientId := os.Getenv("KEYCLOAK_CLIENT_ID")
// 	KeycloakClientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
//
// 	url := KeycloakServer + "/realms/" + KeycloakRealm + "/protocol/openid-connect/token"
// 	method := "POST"
// 	payload := strings.NewReader("grant_type=password&client_id=" + KeycloakClientId + "&client_secret=" + KeycloakClientSecret + "&username=" + username + "&password=" + password)
// 	req, err := http.NewRequest(method, url, payload)
//
// 	if err != nil {
// 		fmt.Println(err)
// 		return "", err
// 	}
//
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 		return "", err
// 	}
//
// 	defer res.Body.Close()
// 	// 	extract the access_token from the response body
//
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return "", err
// 	}
//
// 	var result map[string]interface{}
// 	json.Unmarshal([]byte(body), &result)
//
// 	fmt.Println(result)
//
// 	token := result["access_token"]
// 	if token == nil {
// 		return "", err
// 	}
//
// 	return token.(string), nil
// }
//
// func main() {
// 	// load env variables
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		// log error and exit
// 		fmt.Println("Error loading .env file : ", err)
// 	}
//
// 	// Create a new instance of the application.
// 	router := gin.Default()
//
// 	// Define the route.
// 	router.POST("/api/v1/signup", func(c *gin.Context) {
// 		// Get the JSON body and decode into credentials.
// 		var credentials Credentials
// 		if err := c.ShouldBindJSON(&credentials); err != nil {
// 			// If the structure of the body is wrong, return an HTTP error.
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
//
// 		// Check the credentials.
// 		if err := credentials.isValid(); err != nil {
// 			// If the credentials are invalid, return an HTTP error.
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 			return
// 		}
//
// 		// get access token from keycloak
// 		token, err := getAccessToken()
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 			return
// 		}
//
// 		// create a keycloak user
// 		fmt.Println("Creating user with username : ", credentials.Email, token)
// 		err = createKeycloakUser(credentials.Email, credentials.Password, token)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 			return
// 		}
//
// 		// If the credentials are valid, return an HTTP success status code.
// 		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
// 	})
//
// 	// Define the route.
// 	router.POST("/api/v1/login", func(c *gin.Context) {
// 		// Get the JSON body and decode into credentials.
// 		var credentials Credentials
// 		if err := c.ShouldBindJSON(&credentials); err != nil {
// 			// If the structure of the body is wrong, return an HTTP error.
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
//
// 		// Check the credentials.
// 		if err := credentials.isValid(); err != nil {
// 			// If the credentials are invalid, return an HTTP error.
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 			return
// 		}
//
// 		// get access token from keycloak
// 		token, err := login(credentials.Email, credentials.Password)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 			return
// 		}
//
// 		// If the credentials are valid, return an HTTP success status code.
// 		c.JSON(http.StatusOK, gin.H{"status": "you are logged in", "token": token})
// 	})
//
// 	// Run the application.
// 	router.Run(":3000")
// }
