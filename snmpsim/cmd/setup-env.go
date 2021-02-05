package cmd

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// setupEnvCmd represents the setupEnv command
var setupEnvCmd = &cobra.Command{
	Use:   "setup-env",
	Args:  cobra.ExactArgs(0),
	Short: "Sets up a lab environment",
	Long: `Sets up a laboratory environment.

When invoked without any flags a prompt will start and you can enter how many components of each kind you want to create. 
Additionally asking for all required data for each component.

When invoked with '--env-config' flag it can read the data contained in the given env-config and set up an environment accordingly.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Load the client data from the config
		baseURL := viper.GetString("mgmt.http.baseURL")
		username := viper.GetString("mgmt.http.authUsername")
		password := viper.GetString("mgmt.http.authPassword")

		//Create a new client
		client, err := snmpsimclient.NewManagementClient(baseURL)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while creating management client")
			os.Exit(1)
		}
		err = client.SetUsernameAndPassword(username, password)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while setting username and password")
			os.Exit(1)
		}

		//Generate a new tag for the environment
		tagName, err := randToken(8)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error during creation of hex token for tag")
			os.Exit(1)
		}
		tag, err := client.CreateTag("Env#"+tagName, "Created via the setup-env command")
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while creating tag")
			os.Exit(1)
		}

		//Check if env-config flag is set
		if cmd.Flag("env-config").Changed {
			//Read in the config file path
			file := cmd.Flag("env-config").Value.String()

			//Read in the data from the config
			fileContents, err := ioutil.ReadFile(file)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Could not read file")
				os.Exit(1)
			}

			//Unmarshal the data
			var environment env
			err = yaml.Unmarshal(fileContents, &environment)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Could not unmarshal")
				os.Exit(1)
			}

			//Create the environment according to the config
			for _, lab := range environment.Labs {
				labID := createObject("lab", tag.ID, lab.Name)
				for _, agent := range lab.Agents {
					agentID := createObject("agent", tag.ID, agent.Name, agent.DataDir, strconv.Itoa(labID))
					for _, engine := range agent.Engines {
						engineID := createObject("engine", tag.ID, engine.Name, engine.engineID, strconv.Itoa(agentID))
						for _, endpoint := range engine.Endpoints {
							createObject("endpoint", tag.ID, endpoint.Name, endpoint.Address, endpoint.Protocol, strconv.Itoa(engineID))
						}
						for _, user := range engine.Users {
							createObject("user", tag.ID, user.Name, user.User, user.AuthKey, user.AuthProto, user.PrivKey, user.PrivProto, strconv.Itoa(engineID))
						}
					}
				}
			}

		} else {
			//Define a structure for the user input
			objects := [5]string{"lab", "agent", "engine", "endpoint", "user"}

			objectFields := map[string][]string{
				"lab": {
					"a name",
				},
				"agent": {
					"a name",
					"a dataDir",
					"a labID",
				},
				"engine": {
					"a name",
					"an engineID",
					"an agentID",
				},
				"endpoint": {
					"a name",
					"an address",
					"a protocol",
					"an engineID",
				},
				"user": {
					"a name",
					"an user",
					"an authKey",
					"an authProto",
					"a privKey",
					"a privProto",
					"an engineID",
				},
			}

			for _, object := range objects {
				createObjects(object, objectFields[object], tag.ID)
			}
		}

		fmt.Print("\n")
		fmt.Println("Environment", tag.Name, "has been created successfully.")
		fmt.Println("Id", tag.ID)
	},
}

func init() {
	rootCmd.AddCommand(setupEnvCmd)
	setupEnvCmd.Flags().String("env-config", "", "Sets the config file for the environment")
}

/*
Function to create multiple objects based on user input
*/
func createObjects(object string, fields []string, tagID int) {
	//Create new reader
	reader := bufio.NewReader(os.Stdin)

	//Read in the amount of objects that will be created
	fmt.Print("How many ", object, "s do you want to create? ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error while retrieving input")
		os.Exit(1)
	}
	input = strings.Replace(input, "\n", "", -1)

	amount, err := strconv.Atoi(input)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error while converting " + input + " from string to int")
		os.Exit(1)
	}

	if amount <= 0 {
		log.Error().
			Err(err).
			Msg("Only values above 0 allowed")
		os.Exit(1)
	}

	//Loop for amount of objects to create
	for i := 1; i <= amount; i++ {
		var userInput []string

		//Loop over the objects fields
		for j := 0; j < len(fields); j++ {
			fmt.Print("Please enter ", fields[j], "(", i, "/", amount, "): ")
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Error().
					Err(err).
					Msg("Error while retrieving input")
				os.Exit(1)
			}
			userInput = append(userInput, strings.Replace(line, "\n", "", -1))
		}

		//Create according object
		switch object {
		case "lab":
			createObject(object, tagID, userInput[0])
		case "agent":
			createObject(object, tagID, userInput[0], userInput[1], userInput[2])
		case "engine":
			createObject(object, tagID, userInput[0], userInput[1], userInput[2])
		case "endpoint":
			createObject(object, tagID, userInput[0], userInput[1], userInput[2], userInput[3])
		case "user":
			createObject(object, tagID, userInput[0], userInput[1], userInput[2], userInput[3], userInput[4], userInput[5], userInput[6])
		default:
			log.Debug().
				Err(err).
				Msg("Invalid object " + object)
			os.Exit(1)
		}
	}

}

/*
Function to create one object based on the inputs
*/
func createObject(objectType string, tagID int, args ...string) int {
	//Load the client data from the config
	baseURL := viper.GetString("mgmt.http.baseURL")
	username := viper.GetString("mgmt.http.authUsername")
	password := viper.GetString("mgmt.http.authPassword")

	//Create a new client
	client, err := snmpsimclient.NewManagementClient(baseURL)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error while creating management client")
		os.Exit(1)
	}
	err = client.SetUsernameAndPassword(username, password)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error while setting username and password")
		os.Exit(1)
	}

	//Create according object
	var id int
	switch objectType {
	case "lab":
		//Create a tagged lab
		lab, err := client.CreateLabWithTag(args[0], tagID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while creating lab")
			os.Exit(1)
		}
		fmt.Println("Lab", args[0], "has been created with the id", lab.ID)
		id = lab.ID
	case "agent":
		//Create a tagged agent
		agent, err := client.CreateAgentWithTag(args[0], args[1], tagID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while creating agent")
			os.Exit(1)
		}
		fmt.Println("Agent", args[0], "has been created with the id", agent.ID)

		//Read in the lab-id
		labID, err := strconv.Atoi(args[2])
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while converting " + args[2] + "from string to int")
			os.Exit(1)
		}

		//Add the agent to the lab
		err = client.AddAgentToLab(labID, agent.ID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while adding agent to lab")
			os.Exit(1)
		}
		fmt.Println("Agent", agent.ID, "has been added to lab", labID)
		id = agent.ID
	case "engine":
		//Create a tagged engine
		engine, err := client.CreateEngineWithTag(args[0], args[1], tagID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while creating engine")
			os.Exit(1)
		}
		fmt.Println("Engine", args[0], "has been created with the id", engine.ID)

		//Read in the agent-id
		agentID, err := strconv.Atoi(args[2])
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while converting " + args[2] + "from string to int")
			os.Exit(1)
		}

		//Add the engine to the agent
		err = client.AddEngineToAgent(agentID, engine.ID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while adding engine to agent")
			os.Exit(1)
		}
		fmt.Println("Engine", engine.ID, "has been added to agent", agentID)
		id = engine.ID
	case "endpoint":
		//Create a tagged endpoint
		endpoint, err := client.CreateEndpointWithTag(args[0], args[1], args[2], tagID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while creating endpoint")
			os.Exit(1)
		}
		fmt.Println("Endpoint", args[0], "has been created with the id", endpoint.ID)

		//Read in the engine-id
		engineID, err := strconv.Atoi(args[3])
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while converting " + args[3] + "from string to int")
			os.Exit(1)
		}

		//Add the endpoint to the engine
		err = client.AddEndpointToEngine(engineID, endpoint.ID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while adding endpoint to engine")
			os.Exit(1)
		}
		fmt.Println("Endpoint", endpoint.ID, "has been added to engine", engineID)
		id = endpoint.ID
	case "user":
		//Create a tagged user
		user, err := client.CreateUserWithTag(args[0], args[1], args[2], args[3], args[4], args[5], tagID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while creating user")
			os.Exit(1)
		}
		fmt.Println("User", args[0], "has been created with the id", user.ID)

		//Read in the engine-id
		engineID, err := strconv.Atoi(args[6])
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while converting " + args[6] + "from string to int")
			os.Exit(1)
		}

		//Add the user to the engine
		err = client.AddUserToEngine(engineID, user.ID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while adding user to engine")
			os.Exit(1)
		}
		fmt.Println("User", user.ID, "has been added to engine", engineID)
		id = user.ID
	default:
		log.Debug().
			Msg("Invalid object-type " + objectType)
	}

	return id
}

/*
Generates a random hex-token
*/
func randToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

//Structs needed to unmarshal the yaml config
type env struct {
	Labs map[int]lab `yaml:"Labs"`
}

type lab struct {
	Name   string `yaml:"Name"`
	Agents agents `yaml:"Agents"`
}

type agents map[int]agent

type agent struct {
	Name    string  `yaml:"Name"`
	DataDir string  `yaml:"DataDir"`
	Engines engines `yaml:"Engines"`
}

type engines map[int]engine

type engine struct {
	Name      string    `yaml:"Name"`
	engineID  string    `yaml:"engineID"`
	Endpoints endpoints `yaml:"Endpoints"`
	Users     users     `yaml:"Users"`
}

type endpoints map[int]endpoint

type endpoint struct {
	Name     string `yaml:"Name"`
	Address  string `yaml:"Address"`
	Protocol string `yaml:"Protocol"`
}

type users map[int]user

type user struct {
	Name      string `yaml:"Name"`
	User      string `yaml:"User"`
	AuthKey   string `yaml:"AuthKey"`
	AuthProto string `yaml:"AuthProto"`
	PrivKey   string `yaml:"PrivKey"`
	PrivProto string `yaml:"PrivProto"`
}
