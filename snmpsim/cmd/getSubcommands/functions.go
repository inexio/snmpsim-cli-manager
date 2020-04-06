package getsubcommands

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"strconv"
	"strings"
)

/*
Custom function which prints data in a given data format
*/
func printData(rawData interface{}, format string, prettified bool, depth int) error {
	switch format {
	case "xml":
		data, err := toXML(rawData, prettified)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		break
	case "json":
		data, err := toJSON(rawData, prettified)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		break
	case "human-readable":
		data := toHumanReadable(rawData, depth)
		fmt.Println(data)
		break
	default:
		break
	}

	return nil
}

/*
Converts data into the json format
*/
func toJSON(data interface{}, prettified bool) ([]byte, error) {
	var err error
	var jsonData []byte

	if prettified {
		jsonData, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, err
		}
	} else {
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	return jsonData, nil
}

/*
Converts data into the xml format
*/
func toXML(data interface{}, prettified bool) ([]byte, error) {
	var err error
	var xmlData []byte

	//Maps cant be directly marshalled
	if reflect.ValueOf(data).Kind() == reflect.Map {
		data = data.(*[]uint8)
	}

	if prettified {
		xmlData, err = xml.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, err
		}
	} else {
		xmlData, err = xml.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	return xmlData, nil
}

/*
Converts data into an easily human readable format
*/
func toHumanReadable(data interface{}, limit int) string {
	//the standard limit is 3
	output := superReflect(reflect.ValueOf(data), 0, limit)

	return output
}

/*
Custom function to properly generate human readable outputs
*/
func superReflect(reflectedData reflect.Value, depth int, limit int) string {
	var outputString string

	kind := reflectedData.Type().Kind()
	if kind != reflect.String && kind != reflect.Int && kind != reflect.Float64 && kind != reflect.Ptr {
		if depth > 0 {
			outputString += "\n"
		}
		outputString += strings.Repeat("  ", depth)
		outputString += reflectedData.Type().Name()
	}
	depth++

	switch kind {
	case reflect.Struct:
		outputString += "\n"
		for i := 0; i < reflectedData.NumField(); i++ {
			if reflectedData.Type().Field(i).Type.Kind() != reflect.Slice {
				outputString += strings.Repeat("  ", depth)
				outputString += reflectedData.Type().Field(i).Name + ": "
			}
			outputString += superReflect(reflectedData.Field(i), depth, limit)
		}
	case reflect.Slice:
		outputString += "(" + strconv.Itoa(reflectedData.Len()) + ") \n"
		if reflectedData.Len() == 0 {
			outputString += strings.Repeat("  ", depth)
			outputString += "/"
		}
		for j := 0; j < reflectedData.Len(); j++ {
			if depth < limit {
				outputString += superReflect(reflectedData.Index(j), depth, limit)
			}
		}
		outputString += "\n"
	case reflect.Map:
		outputString += "(" + strconv.Itoa(reflectedData.Len()) + ") \n"
		for _, key := range reflectedData.MapKeys() {
			outputString += strings.Repeat("  ", depth)
			outputString += key.String() + ": "
			if depth < limit {
				outputString += superReflect(reflectedData.MapIndex(key), depth, limit)
			}
		}
	case reflect.Int:
		fieldValue := strconv.Itoa(int(reflectedData.Int()))
		outputString += fieldValue + "\n"
	case reflect.String:
		fieldValue := reflectedData.String()
		outputString += fieldValue + "\n"
	case reflect.Float64:
		fieldValue := strconv.FormatFloat(reflectedData.Float(), 'f', -1, 64)
		outputString += fieldValue + "\n"
	case reflect.Ptr:
		indirect := reflect.Indirect(reflectedData)
		if indirect.Kind() == reflect.Invalid {
			outputString += "/\n"
		} else {
			outputString += strconv.Itoa(int(indirect.Int())) + "\n"
		}
	default:
		log.Debug().
			Msg("Could not reflect " + reflectedData.Type().Kind().String())
	}

	depth--
	return outputString
}

/*
Parses the persistent flags of the get sub-commands
*/
func parsePersistentFlags(cmd *cobra.Command) (string, int, bool) {
	format := cmd.Flag("format").Value.String()
	if !validateFormat(format) {
		log.Error().
			Msg("Invalid format")
		os.Exit(1)
	}

	depth, err := cmd.Flags().GetInt("depth")
	if err != nil {
		log.Error().
			Msg("Error while retrieving 'depth' flag")
		os.Exit(1)
	}

	prettified, err := strconv.ParseBool(cmd.Flag("pretty").Value.String())
	if err != nil {
		log.Error().
			Msg("Error during conversion of 'pretty' flag")
		os.Exit(1)
	}

	return format, depth, prettified
}

/*
Performs a check on the given format if it is either valid or invalid
*/
func validateFormat(format string) bool {
	for _, allowedFormat := range []string{"xml", "json", "human-readable"} {
		if format == allowedFormat {
			return true
		}
	}

	return false
}

/*
Parses the given filters of all getComponents functions
*/
func parseFilters(command *cobra.Command) map[string]string {
	//Get the commands local flags
	flags := command.LocalFlags()

	//Create the filters var
	filters := make(map[string]string)

	//Create a reflection of the commands flags
	reflectedFlags := reflect.ValueOf(flags).Elem().FieldByName("formal").MapKeys()

	//fill filters map with the flags and their values
	for _, val := range reflectedFlags {
		if val.String() == "help" {
			continue
		}
		if command.Flag(val.String()).Changed {
			filters[val.String()] = command.Flag(val.String()).Value.String()
		}
	}

	if len(filters) == 0 {
		return nil
	}

	return filters
}
