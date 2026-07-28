package game

import (
	"encoding/json"
	"image/color"
	"log"
	"os"
	"path/filepath"
)

type levelObjectJSON struct {
	Type          string     `json:"type"`
	Name          string     `json:"name"`
	PositionPx    Vector     `json:"positionPx"`
	Radius        float64    `json:"radius"`
	WidthPx       float64    `json:"widthPx"`
	HeightPx      float64    `json:"heightPx"`
	Color         color.RGBA `json:"color"`
	BorderColor   color.RGBA `json:"borderColor"`
	BorderWidthPx float64    `json:"borderWidthPx"`
	Charge        float64    `json:"charge"`
}

// load all objects from JSON file.
// Using JSON so it's easy to edit (at least until I make some kind of level editor app)
// would XML have been better for this?
func loadLevelFromJSON(path string) *Level {
	file, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		log.Printf("failed to read level file %s: %v", path, err)
		return nil
	}

	var objectData []levelObjectJSON
	if err := json.Unmarshal(file, &objectData); err != nil {
		log.Printf("failed to parse level file %s: %v", path, err)
		return nil
	}

	levelName := "NOLEVEL"
	backgroundColor := color.RGBA{10, 10, 20, 255}
	var player *Player
	var playerZone *PlayerZone
	var goal *GoalObject

	objects := make([]ChargedObject, 0, len(objectData))
	obstacles := make([]Obstacle, 0, len(objectData))
	for _, data := range objectData {
		switch data.Type {
		case "obstacle":
			obstacles = append(obstacles, Obstacle{
				positionPx: data.PositionPx,
				widthPx:    data.WidthPx,
				heightPx:   data.HeightPx,
				color:      data.Color,
			})
		case "chargedObject":
			objects = append(objects, ChargedObject{
				positionPx: data.PositionPx,
				radius:     data.Radius,
				color:      data.Color,
				charge:     data.Charge,
			})
		case "goal":
			goal = &GoalObject{
				ChargedObject: ChargedObject{
					positionPx: data.PositionPx,
					radius:     data.Radius,
					color:      data.Color,
					charge:     data.Charge},
				borderColor:   data.BorderColor,
				borderWidthPx: data.BorderWidthPx,
			}
			// add just ChargedObject portion of goal to objects list
			objects = append(objects, goal.ChargedObject)
		case "playerZone":
			playerZone = &PlayerZone{
				positionPx:    data.PositionPx,
				widthPx:       data.WidthPx,
				heightPx:      data.HeightPx,
				color:         data.Color,
				borderColor:   data.BorderColor,
				borderWidthPx: data.BorderWidthPx,
			}
		case "player":
			player = NewPlayer(data.PositionPx)
		case "levelInfo":
			levelName = data.Name
			backgroundColor = data.Color
		default:

		}
	}

	// require a goal and player be defined
	if player == nil {
		log.Printf("Missing Player object definition in %s: %v", path)
		return nil
	}
	if goal == nil {
		log.Printf("Missing Goal object definition in %s: %v", path)
		return nil
	}

	return &Level{
		player:          player,
		goal:            goal,
		name:            levelName,
		objects:         objects,
		obstacles:       obstacles,
		playZone:        playerZone,
		backgroundColor: backgroundColor,
	}
}
