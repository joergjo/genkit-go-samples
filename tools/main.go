package main

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

type LightOperation struct {
	State string `json:"state" jsonschema_description:"State of lights, either 'on' or 'off'"`
	Room  string `json:"room" jsonschema_description:"Room where the lights should be operated"`
}

type WindowOperation struct {
	State string `json:"state" jsonschema_description:"State of windows, either 'open' or 'closed'"`
	Room  string `json:"room" jsonschema_description:"Room where the windows should be operated"`
}

type TVOperation struct {
	Movie string `json:"movie" jsonschema_description:"Movie to play on the TV"`
	Room  string `json:"room" jsonschema_description:"Room where the TV should be operated"`
}

type GarageOperation struct {
	Action string `json:"state" jsonschema_description:"Action to perform on the garage door. Must be either 'open' or 'close'"`
}

const systemPrompt = `
    You are an assistant that helps users control their home automation system. You can turn on/off lights, 
    open/close windows, and play movies on the TV. You can also close garage doors.

    Whenever you receive a request, you will fulfill it by calling the appropriate functions. If a request
    involves multiple actions, you will call multiple functions in sequence.
    Deny all requests that are not related to home automation.
`

func main() {
	ctx := context.Background()
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"))

	operateLight := genkit.DefineTool(g, "operateLight", "Turns the lights on or off in the living room, kitchen, bedroom or garage",
		func(ctx *ai.ToolContext, input LightOperation) (string, error) {
			rooms := []string{"living room", "kitchen", "bedroom", "bathroom", "garage"}
			if !slices.Contains(rooms, input.Room) {
				return "", fmt.Errorf("invalid room %s", rooms)
			}
			s := fmt.Sprintf("Changed status of the %s lights to %s", input.Room, input.State)
			fmt.Println(s)
			return s, nil
		})

	operateWindow := genkit.DefineTool(g, "operateWindow", "Opens or closes the windows of the living room or bedroom",
		func(ctx *ai.ToolContext, input WindowOperation) (string, error) {
			rooms := []string{"living room", "bedroom"}
			if !slices.Contains(rooms, input.Room) {
				return "", fmt.Errorf("invalid room %s", rooms)
			}
			s := fmt.Sprintf("Changed status of the %s windows to %s", input.Room, input.State)
			fmt.Println(s)
			return s, nil
		})

	operateTV := genkit.DefineTool(g, "operateTV", "Plays a movie on the TV in the living room or bedroom",
		func(ctx *ai.ToolContext, input TVOperation) (string, error) {
			rooms := []string{"living room", "bedroom"}
			if !slices.Contains(rooms, input.Room) {
				return "", fmt.Errorf("invalid room %s", rooms)
			}
			s := fmt.Sprintf("Playing movie %s on the TV in the %s", input.Movie, input.Room)
			fmt.Println(s)
			return s, nil
		})
	operateGarage := genkit.DefineTool(g, "operateGarage", "Opens or closes the garage door",
		func(ctx *ai.ToolContext, input GarageOperation) (string, error) {
			s := fmt.Sprintf("Garage door is now %s", input.Action)
			fmt.Println(s)
			return s, nil
		})

	prompts := []string{
		"Turn on the lights in the kitchen.",
		"Open the windows of the bedroom, turn the lights off and put on Shawnshank Redemption on the TV.",
		"Close the garage door and turn off the lights in all rooms.",
		"Turn off the lights in all rooms and play a movie in which Tom Cruise plays a lawyer in the living room.",
	}
	for _, p := range prompts {
		fmt.Println("Prompt: ", p)
		resp, err := genkit.Generate(ctx, g, ai.WithSystem(systemPrompt), ai.WithPrompt(p),
			ai.WithTools(operateLight, operateWindow, operateTV, operateGarage))
		if err != nil {
			log.Fatalf("failed to generate response: %v", err)
		}
		fmt.Println("Response:", resp.Text())
	}
}
